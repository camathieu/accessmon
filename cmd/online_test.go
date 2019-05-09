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

func TestOnlineStreamTime(t *testing.T) {

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

	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	config := &accessmon.Config{AlertWindow: 5 * time.Second, AlertThreshold: 10}
	mon := accessmon.NewMonitor(config)

	shutdown, err := tailLogFile(tmpfile.Name(), 5*time.Second, mon)
	require.NoError(t, err)
	defer shutdown()

	time.Sleep(1 * time.Second)

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

	for _, req := range requests {
		_, err = tmpfile.WriteString(req.String() + "\n")
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)

	shutdown()

	require.Len(t, mon.Alerts(), 2)
}

func TestOnlineRealtime(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	config := &accessmon.Config{AlertWindow: 2 * time.Second, AlertThreshold: 5}
	mon := accessmon.NewMonitor(config)

	shutdown, err := tailLogFile(tmpfile.Name(), 1*time.Second, mon)
	require.NoError(t, err)
	defer shutdown()

	time.Sleep(time.Second)

	stop := time.After(5 * time.Second)
	tick := time.Tick(10 * time.Millisecond)

LOOP:
	for {
		select {
		case <-stop:
			break LOOP
		case <-tick:

			req := &accessmon.Request{
				SourceIP:    net.ParseIP("::1"),
				User:        "user",
				Time:        time.Now(),
				Method:      "GET",
				Path:        "/path",
				Section:     "/",
				HTTPVersion: "HTTP/2.0",
				Code:        500,
				Size:        42,
			}

			_, err = tmpfile.WriteString(req.String() + "\n")
			require.NoError(t, err)
		}
	}

	time.Sleep(time.Second)

	shutdown()

	require.Len(t, mon.Alerts(), 1)
}

func TestOnlineNoData(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	config := &accessmon.Config{AlertWindow: 5 * time.Second, AlertThreshold: 10}
	mon := accessmon.NewMonitor(config)

	shutdown, err := tailLogFile(tmpfile.Name(), 1*time.Second, mon)
	require.NoError(t, err)
	defer shutdown()

	time.Sleep(3 * time.Second)

	shutdown()

	require.Len(t, mon.Alerts(), 0)
}

func TestOnlineFileNotFound(t *testing.T) {
	_, err := tailLogFile("invalid_file_name", 0, nil)
	require.Error(t, err)
}

func TestOnlineInvalidRefreshInterval(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "access.log_")
	require.NoError(t, err)

	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	_, err = tailLogFile(tmpfile.Name(), 0, nil)
	require.Error(t, err)
}
