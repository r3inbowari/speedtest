package main

import (
	"github.com/r3inbowari/zlog"
	"gopkg.in/alecthomas/kingpin.v2"
	"speedtest"
	"speedtest/internal/location"
)

var (
	list     = kingpin.Flag("list", "show a list of nearby servers depends on ip").Short('l').Bool()
	download = kingpin.Flag("download", "select a server and download test only").Short('d').Bool()
	upload   = kingpin.Flag("upload", "select a server and upload test only").Short('u').Bool()
)

func main() {
	kingpin.Parse()
	l := zlog.NewLogger()
	l.SetScreen(true)
	st := speedtest.InitSpeedTest(speedtest.Options{Log: &l.Logger, Location: location.HongKong})
	st.Log.Info("Hi, SpeedTest!")

	if *list {
		st.CmdShowServerList()
		return
	}

	if *download {
		st.CmdDownloadTest()
		return
	}

	if *upload {
		st.CmdUploadTest()
		return
	}

	st.CmdSingleTest()
}
