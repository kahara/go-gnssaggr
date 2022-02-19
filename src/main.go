package main

import (
	"github.com/kahara/go-gnssaggr/src/aggregator"
	"github.com/kahara/go-gnssaggr/src/exporter"
	"github.com/kahara/go-gnssaggr/src/reader"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		reports    = make(chan reader.Report, 2)
		aggregates = make(chan aggregator.Aggregation, 2)
		sigs       = make(chan os.Signal, 1)
		sig        os.Signal
	)

	zerolog.TimeFieldFormat = time.RFC3339Nano

	// FIXME read host and port from command line
	config := reader.Config{
		Host: "green.lan",
		Port: 2947,
	}

	go reader.Read(config, reports)
	go aggregator.Aggregate(reports, aggregates)
	go exporter.Export(aggregates)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig = <-sigs:
		log.Printf("Received %s signal", sig)
	}
}
