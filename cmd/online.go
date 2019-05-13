package main

import (
	"errors"
	"io"
	"math"
	"time"

	"github.com/camathieu/accessmon"
	"github.com/hpcloud/tail"
)

func tailLogFile(path string, refreshInterval time.Duration, mon *accessmon.Monitor) (shutdown func(), err error) {

	if refreshInterval <= 0 {
		return func() {}, errors.New("missing refresh interval")
	}

	// Open and tail file

	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true, MustExist: true, Location: &tail.SeekInfo{Whence: io.SeekEnd}})
	if err != nil {
		return func() {}, err
	}

	// Cleanup to call on exit

	shutdown = func() {
		_ = t.Stop()
		t.Cleanup()
	}

	ticker := time.Tick(refreshInterval)

	window := int(math.Max(10, float64(refreshInterval)))

	go func() {
	LOOP:
		for {

			// This select statement allow everything to run in a single goroutine
			// avoiding the need for synchronization in the accessmon package

			select {
			case <-ticker:
				// Update display
				cleanDisplay()

				deadline := time.Now().Add(-time.Duration(window))
				if mon.Last().After(deadline) {
					displayStats(mon.Stats(refreshInterval, 1), mon.Last(), refreshInterval)
				} else {

					// It's important to note that this program does event time stream processing
					// which is very well explained in this doc : https://ci.apache.org/projects/flink/flink-docs-stable/dev/event_time.html
					// Not having received any data in the last interval does not mean there were actually no requests on the server,
					// logs storage can just be slow (NFS mount points) for example.
					// As we can make no assumptions on the time of the next log line we will receive
					// It's better to display a proper warning than just updating the display with 0 request per seconds
					// It's also more effective to change the format of the output to catch the operator eyes in case of such event

					displayStats(nil, time.Now(), refreshInterval)
				}

				displayAlerts(mon.Alerts())
			case line := <-t.Lines:
				// Update monitor
				if line == nil {
					break LOOP
				}

				// just discard invalid lines
				_, _ = mon.AddLine(line.Text)
			}
		}
	}()

	return shutdown, nil
}
