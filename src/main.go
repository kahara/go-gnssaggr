package main

import (
	"flag"
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

	// Command line arguments
	gpsdHost := flag.String("host", "localhost", "gpsd host")
	gpsdPort := flag.Uint("port", 2947, "gpsd port")
	prometheusPort := flag.Uint("promport", 9100, "Prometheus exporter port")
	flag.Parse()

	readerConfig := reader.Config{
		Host: *gpsdHost,
		Port: uint16(*gpsdPort),
	}

	exporterConfig := exporter.Config{
		Port: uint16(*prometheusPort),
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
