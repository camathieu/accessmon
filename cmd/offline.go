package main

import (
	"bufio"
	"os"

	"github.com/camathieu/accessmon"
)

func catLogFile(path string, mon *accessmon.Monitor) (err error) {

	// Open file

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read file line by line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Add to the monitor

		alert, err := mon.AddLine(scanner.Text())
		if err != nil {
			continue
		}

		// Display alert if any

		if alert != nil {
			displayAlertOffline(alert)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
