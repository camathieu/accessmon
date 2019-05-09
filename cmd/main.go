package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camathieu/accessmon"
)

func cleanDisplay() {
	fmt.Print("\033[H\033[2J")
}

func perSecond(value int, window time.Duration) float64 {
	return float64(value) / window.Seconds()
}

func displayStats(stats *accessmon.Stats, now time.Time, window time.Duration) {
	fmt.Println("")

	// This time is the date of the last log received

	fmt.Printf("Date : %s\n", now)

	if stats == nil {
		fmt.Printf("Nothing to process in the last %s\n", window)
		fmt.Println("")
		return
	}

	fmt.Printf("Total %.3f req/s\n", perSecond(stats.Count, window))
	fmt.Printf(" ServerError : %.3f%%\n", stats.ServerError)
	fmt.Printf(" HTTP2 : %.1f%%\n", stats.HTTP2)
	fmt.Printf(" IPv6 : %.1f%%\n", stats.Ipv6)
	fmt.Println("")
	fmt.Printf(" top source : %s (%.1f req/s)\n", stats.TopSources[0].Key, perSecond(stats.TopSources[0].Count, window))
	fmt.Printf(" top section : %s (%.1f req/s)\n", stats.TopSection[0].Key, perSecond(stats.TopSection[0].Count, window))
	fmt.Printf(" top user : %s (%.1f req/s)\n", stats.TopUsers[0].Key, perSecond(stats.TopUsers[0].Count, window))
	fmt.Println("")
}

func displayAlerts(alerts []*accessmon.Alert) {
	for _, alert := range alerts {
		displayAlert(alert)
	}
}

func displayAlert(alert *accessmon.Alert) {
	fmt.Printf("AL - High traffic above threshold at %s ( %.3f requests per second )\n", alert.Start, alert.Value)
	if !alert.IsOngoing() {
		fmt.Printf("OK - High traffic under threshold at %s. Alert duration %s\n", alert.End, alert.End.Sub(alert.Start))
	}
}

func displayAlertOffline(alert *accessmon.Alert) {
	if alert.IsOngoing() {
		fmt.Printf("AL - High traffic above threshold at %s ( %.3f requests per second )\n", alert.Start, alert.Value)
	} else {
		fmt.Printf("OK - High traffic under threshold at %s. Alert duration %s\n", alert.End, alert.End.Sub(alert.Start))
	}
}

func main() {
	path := flag.String("logfile", "/tmp/access.log", "log file path")
	refresh := flag.Duration("refresh", 10*time.Second, "screen refresh interval ( online mode only )")
	offline := flag.Bool("offline", false, "offline mode ( cat )")
	generate := flag.Bool("generate", false, "generator mode")

	config := &accessmon.Config{}
	flag.DurationVar(&config.AlertWindow, "window", 2*time.Minute, "total request per second moving average alerting window")
	flag.Float64Var(&config.AlertThreshold, "threshold", 10, "total request per second moving average alerting threshold")

	flag.Parse()

	if *generate {

		// Handy generator mode

		err := generator(*path)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	// we need to store at least enough requests in memory to generate statistics
	// for the last refresh interval
	config.StoreWindow = *refresh

	mon := accessmon.NewMonitor(config)

	if *offline {
		err := catLogFile(*path, mon)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		shutdown, err := tailLogFile(*path, *refresh, mon)
		if err != nil {
			log.Fatal(err)
		}

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			shutdown()
			os.Exit(0)
		}()

		select {}
	}
}
