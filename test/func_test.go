package test

import (
	"fmt"
	"golang.org/x/net/context"
	"speedtest"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	st := speedtest.InitSpeedTest(speedtest.Options{})
	err := st.GetServerList()
	if err != nil {
		t.Fail()
	}
	(st.ServerList)[0].Connect(time.Minute)
	s := (st.ServerList)[0]
	// ping
	err = s.Ping()
	if err != nil {
		t.Fail()
	}
	// print
	fmt.Printf(s.Result.String())
}

func TestDownload(t *testing.T) {
	st := speedtest.InitSpeedTest(speedtest.Options{})
	err := st.GetServerList()
	if err != nil {
		t.Fail()
	}
	s := (st.ServerList)[0]
	// download
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	err = s.DownloadTest(ctx)
	if err != nil {
		t.Fail()
	}
	// print
	fmt.Printf(s.Result.String())
}

func TestUpload(t *testing.T) {
	st := speedtest.InitSpeedTest(speedtest.Options{})
	err := st.GetServerList()
	if err != nil {
		t.Fail()
	}
	s := (st.ServerList)[0]
	// upload
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	err = s.UploadTest(ctx)
	if err != nil {
		t.Fail()
	}
	// print
	fmt.Printf(s.Result.String())
}
