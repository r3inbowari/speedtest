package speedtest

import "github.com/r3inbowari/common"

type SpeedTest struct {
	ServerList *[]Server

	Options
}

type Options struct {
}

func InitSpeedTest(opts Options) *SpeedTest {
	return &SpeedTest{ServerList: &[]Server{}}
}

// GetServerList get a list of SpeedTest servers
// host https://www.speedtest.net/
// path /api/js/servers
// proto https://www.speedtest.net/api/js/servers
// origin https://www.speedtest.net/api/js/servers
// param engine
// param limit
// param https_functional
// param lat 39.9 in Beijing
// param lon 116.4 in Beijing
func (st *SpeedTest) GetServerList() error {
	_, err := common.RequestJson(common.RequestOptions{
		Url: "https://www.speedtest.net/api/js/servers",
	}, st.ServerList)
	return err
}
