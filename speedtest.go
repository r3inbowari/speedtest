package speedtest

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/r3inbowari/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"speedtest/internal/api"
	"speedtest/internal/location"
	"sync"
	"time"
)

const PrimaryServerApi = "https://www.speedtest.net/api/js/servers"

type SpeedTest struct {
	ServerList []*api.Server

	Options
}

type Options struct {
	Log      *logrus.Logger
	Location *location.Location
}

func InitSpeedTest(opts Options) *SpeedTest {
	if opts.Log == nil {
		opts.Log = logrus.New()
	}
	return &SpeedTest{Options: opts}
}

// GetServerList get a list of SpeedTest servers
// host https://www.speedtest.net/
// path /api/js/servers
// proto https://www.speedtest.net/api/js/servers
// origin https://www.speedtest.net/api/js/servers
// param engine
// param limit
// param https_functional
// param lat 39.56 in Beijing
// param lon 116.2 in Beijing
func (st *SpeedTest) GetServerList() error {
	url := PrimaryServerApi
	if st.Location != nil {
		url += fmt.Sprintf("?lat=%.4f&lon=%.4f", st.Location.Lat, st.Location.Lon)
	}
	_, err := common.RequestJson(common.RequestOptions{
		Url: url,
	}, &st.ServerList)
	return err
}

func (st *SpeedTest) CmdShowServerList() {
	st.Log.Info("fetching server data...")
	err := st.GetServerList()
	if err != nil {
		st.Log.WithField("err", err.Error()).Error("[SpeedTest] could not fetch the server list")
		time.Sleep(time.Second * 5)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(st.ServerList))
	for i := range st.ServerList {
		st.pingGo(i, &wg)
	}
	wg.Wait()
	for _, server := range st.ServerList {
		_, _ = color.New(color.FgGreen).Printf("[ID:%6s] [%2s] [Region:%12s] [Distance:%5dkm, Lat:%5s, Lot: %5s] [Latency: %7sms] [Jitter: %7sms] -> %s\n", server.Id, server.CC, server.Name, server.Distance, server.Lat, server.Lon, api.MinString(server.Result.LocalRTT), api.JitterString(server.Result.LocalRTT), server.Sponsor)
	}

}

func (st *SpeedTest) pingGo(i int, wg *sync.WaitGroup) {
	s := (st.ServerList)[i]
	st.Log.Info("[SpeedTest] connecting to " + s.Host)
	go func() {
		defer wg.Done()
		if errPing := s.Connect(time.Second * 5); errPing != nil {
			st.Log.WithField("id", s.Id).WithField("name", s.Sponsor).Error("[SpeedTest] could not establish connection to server or it is not supported websocket protocol")
			return
		}
		if errPing := s.Ping(); errPing != nil {
			st.Log.Errorf("[SpeedTest] ping timeout to %s", s.Host)
			return
		}
		st.Log.WithField("id", s.Id).WithField("name", s.Sponsor).Info("[SpeedTest] done")
	}()
}

func (st *SpeedTest) CmdDownloadTest() {

	st.CmdShowServerList()
	_, _ = color.New(color.FgGreen).Print("Enter a ServerID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	id := scanner.Text()
	for i := range st.ServerList {
		if st.ServerList[i].Id == id {

			err := st.ServerList[i].DownloadTest(context.Background())
			if err != nil {
				st.Log.WithField("id", id).Error("[SpeedTest] download failed")
				time.Sleep(time.Second * 5)
				return
			}
			_, _ = color.New(color.FgGreen).Printf("[ID:%6s] [DownloadRate: %.2fMbps]\n", st.ServerList[i].Id, st.ServerList[i].Result.DownloadRate)
			return
		}
	}
	if id != "" {
		st.Log.WithField("id", id).Error("[SpeedTest] not found server")
	}
}

func (st *SpeedTest) CmdUploadTest() {

	st.CmdShowServerList()
	_, _ = color.New(color.FgGreen).Print("Enter a ServerID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	id := scanner.Text()
	for i := range st.ServerList {
		if st.ServerList[i].Id == id {

			err := st.ServerList[i].UploadTest(context.Background())
			if err != nil {
				st.Log.WithField("id", id).Error("[SpeedTest] upload failed")
				time.Sleep(time.Second * 5)
				return
			}
			_, _ = color.New(color.FgGreen).Printf("[ID:%6s] [UploadRate: %.2fMbps]\n", st.ServerList[i].Id, st.ServerList[i].Result.UploadRate)
			return
		}
	}
	if id != "" {
		st.Log.WithField("id", id).Error("[SpeedTest] not found server")
	}
}

// CmdSingleTest single test with ping/down/up
func (st *SpeedTest) CmdSingleTest() {
	err := st.GetServerList()
	if err != nil {
		st.Log.Error(err.Error())
	}

	if len(st.ServerList) == 0 {
		st.Log.WithField("err", err.Error()).Warn("could not find any servers")
		return
	}

	s := (st.ServerList)[0]

	_, _ = color.New(color.FgGreen).Printf("[SpeedTest] testing -> %s %dkm\n", s.Sponsor, s.Distance)

	err = s.Connect(time.Second * 10)
	if err != nil {
		st.Log.WithField("err", err.Error()).Error("unable to connect to the current server, please check your network")
		return
	}

	// ping
	err = s.Ping()
	if err != nil {
		st.Log.WithField("err", err.Error()).Error("ping error")
	}

	// download
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	err = s.DownloadTest(ctx)
	if err != nil {
		st.Log.WithField("err", err.Error()).Error("download error")
	}
	// upload
	ctx, _ = context.WithTimeout(context.Background(), time.Minute)
	err = s.UploadTest(ctx)
	if err != nil {
		st.Log.WithField("err", err.Error()).Error("upload error")
	}

	// print
	s.Result.Print((st.ServerList)[0].Sponsor)
}
