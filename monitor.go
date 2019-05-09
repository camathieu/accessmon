package accessmon

import (
	"time"
)

// Config represent the Monitoring configuration
type Config struct {
	StoreWindow    time.Duration // Time window of parsed lines to keep in-memory
	AlertWindow    time.Duration // Sliding window parameter of the Alerter
	AlertThreshold float64       // Threshold parameter of the Alerter
}

// Monitor holds the different components to analyse a W3C Common Log File line stream
// /!\ NOT THREAD SAFE /!\
type Monitor struct {
	config  *Config  // Monitoring configuration
	parser  Parser   // Parser to parse log entries
	store   *Store   // Store to store parsed lines
	alerter *Alerter // Alerter to check for anomalies

	last time.Time // Time of the last message processed
}

// NewMonitor creates a new monitor from the provided configuration
func NewMonitor(config *Config) (mon *Monitor) {
	mon = &Monitor{
		config: config,
		parser: &W3CParser{},
		store:  &Store{},
	}

	if config.StoreWindow < config.AlertWindow {
		config.StoreWindow = config.AlertWindow
	}

	if config.AlertWindow > 0 && config.AlertThreshold > 0 {
		mon.alerter = NewAlerter(config.AlertWindow, config.AlertThreshold)
	}

	return mon
}

// AddLine parse the line and update the monitor accordingly
// It returns an alert if the line does trigger the start or end of an alert
func (mon *Monitor) AddLine(line string) (alert *Alert, err error) {

	// Parse

	req, err := mon.parser.Parse(line)
	if err != nil {
		return nil, err
	}

	// Store

	err = mon.store.AddRequest(req)
	if err != nil {
		return nil, err
	}

	// Check for alert

	if mon.alerter != nil {
		count := len(mon.store.Since(Deadline(req.Time, mon.config.AlertWindow)))
		value := float64(count) / float64(mon.config.AlertWindow.Seconds())
		alert = mon.alerter.Check(req.Time, value)
	}

	// Clean

	mon.store.Clean(Deadline(req.Time, mon.config.StoreWindow))
	mon.last = req.Time

	return alert, nil
}

// Stats returns summary statistics for the provided time window
func (mon *Monitor) Stats(window time.Duration, top int) (stats *Stats) {
	requests := mon.store.Since(Deadline(mon.last, window))
	return NewStats(requests, top)
}

// Alerts returns any alerts raised by the alerter
func (mon *Monitor) Alerts() (alerts []*Alert) {
	if mon.alerter == nil {
		return nil
	}
	return mon.alerter.Alerts()
}

// Last returns the time of the last message processed
func (mon *Monitor) Last() time.Time {
	return mon.last
}
