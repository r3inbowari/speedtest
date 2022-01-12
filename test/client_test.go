package test

import (
	"speedtest"
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestGetServerList(t *testing.T) {
	st := speedtest.InitSpeedTest(speedtest.Options{})
	err := st.GetServerList()
	if err != nil {
		t.Fail()
	}
	if !assert.NotEqual(t, len(*st.ServerList), 0, "list len not be zero") {
		t.Fail()
	}
}
