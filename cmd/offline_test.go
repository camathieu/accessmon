package main

import (
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/camathieu/accessmon"
	"github.com/stretchr/testify/require"
)

func TestOffline1(t *testing.T) {
	requests := make([]*accessmon.Request, 550)

	for i := range requests {
		req := &accessmon.Request{
			SourceIP:    net.ParseIP("127.0.0.1"),
			User:        "user",
			Time:        start,
			Method:      "GET",
			Path:        "/path",
			Section:     "/",
			HTTPVersion: "HTTP/1.0",
			Code:        200,
			Size:        42,
		}
		requests[i] = req
	}

	// sequence is :
	//  - 10 s at 5pps  -> 50
	//  - 10 s at 20pps -> 200
	//  - 10 s at 5pps  -> 50
	//  - 10 s at 20pps -> 200
	//  - 10 s at 0pps  -> 0
	//  - 10 s at 5pps  -> 50
	//    60 s          -> 550

	index := 0
	now := start
	seq := []int{5, 20, 5, 20, 0, 5}

	for _, speed := range seq {
		for i := 0; i < 10; i++ {
			for j := 0; j < speed; j++ {
				requests[index].Time = now
				index++
			}
			now = now.Truncate(time.Second).Add(time.Second)
		}
	}

	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	for _, req := range requests {
		_, err = tmpfile.WriteString(req.String() + "\n")
		require.NoError(t, err)
	}

	err = tmpfile.Close()
	require.NoError(t, err)

	config := &accessmon.Config{AlertWindow: 5 * time.Second, AlertThreshold: 10}
	mon := accessmon.NewMonitor(config)

	err = catLogFile(tmpfile.Name(), mon)
	require.NoError(t, err)

	require.Len(t, mon.Alerts(), 2)
}

func TestOffline2(t *testing.T) {

	requests := make([]*accessmon.Request, 830)

	for i := range requests {
		req := &accessmon.Request{
			SourceIP:    net.ParseIP("127.0.0.1"),
			User:        "user",
			Time:        start,
			Method:      "GET",
			Path:        "/path",
			Section:     "/",
			HTTPVersion: "HTTP/1.0",
			Code:        200,
			Size:        42,
		}
		requests[i] = req
	}

	// sequence is :
	//  - 10 s at 1pps  -> 10
	//  - 4 s at 100pps -> 400
	//  - 10 s at 1pps  -> 10
	//  - 4 s at 100pps -> 400
	//  - 10 s at 1pps  -> 10
	//    60 s          -> 830

	index := 0
	now := start

	durations := []int{10, 4, 10, 4, 10}
	values := []int{1, 100, 1, 100, 1}

	for window, speed := range values {
		for i := 0; i < durations[window]; i++ {
			for j := 0; j < speed; j++ {
				requests[index].Time = now
				index++
			}
			now = now.Truncate(time.Second).Add(time.Second)
		}
	}

	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	for _, req := range requests {
		_, err = tmpfile.WriteString(req.String() + "\n")
		require.NoError(t, err)
	}

	err = tmpfile.Close()
	require.NoError(t, err)

	config := &accessmon.Config{AlertWindow: 5 * time.Second, AlertThreshold: 10}
	mon := accessmon.NewMonitor(config)

	err = catLogFile(tmpfile.Name(), mon)
	require.NoError(t, err)

	require.Len(t, mon.Alerts(), 2)
}

func TestOfflineNoFile(t *testing.T) {
	err := catLogFile("invalid_file_name", nil)
	require.Error(t, err)
}
