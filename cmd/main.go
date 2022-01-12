package main

import (
	"fmt"
	"golang.org/x/net/context"
	"speedtest"
	"time"
)

func main() {
	st := speedtest.InitSpeedTest(speedtest.Options{})
	err := st.GetServerList()
	if err != nil {
		println(err.Error())
	}
	(*st.ServerList)[0].Connect()
	s := (*st.ServerList)[0]
	// ping
	err = s.Ping()
	if err != nil {
		println(err.Error())
	}

	// download
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	err = s.DownloadTest(ctx)
	if err != nil {
		println(err.Error())
	}
	// upload
	ctx, _ = context.WithTimeout(context.Background(), time.Minute)
	err = s.UploadTest(ctx)
	if err != nil {
		println(err.Error())
	}
	// print
	fmt.Printf(s.Result.String())

	defer func() { time.Sleep(time.Minute) }()
}
