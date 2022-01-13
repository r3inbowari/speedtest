package test

import (
	"fmt"
	"speedtest/internal/api"
	"testing"
	"time"
)

func TestRemoteAvg(t *testing.T) {
	ss := []string{
		"1642004506713",
		"1642004506740",
		"1642004506776",
		"1642004506791",
		"1642004506814",
	}
	f := api.RemoteAverage(ss)
	fmt.Printf("%.2fms\n", f)

	ss = []string{
		"1642004506713",
		"1642004506740",
		"1642004506sd6",
		"1642004506791",
		"1642004506814",
	}
	f = api.RemoteAverage(ss)
	fmt.Printf("%.2fms\n", f)
}

func TestLocalAvg(t *testing.T) {
	dd := []time.Duration{
		time.Millisecond * 25,
		time.Millisecond * 33,
		time.Millisecond * 27,
		time.Millisecond * 41,
		time.Millisecond * 38,
	}
	f := api.LocalAverage(dd)
	fmt.Printf("%.2fms\n", f)
}
