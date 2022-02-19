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
	readerConfig := reader.Config{
		Host: "green.lan",
		Port: 2947,
	}

	// FIXME read port from command line
	exporterConfig := exporter.Config{
		Port: 9101,
	}

	go reader.Read(readerConfig, reports)
	go aggregator.Aggregate(reports, aggregates)
	go exporter.Export(exporterConfig, aggregates)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig = <-sigs:
		log.Printf("Received %s signal", sig)
	}
}
