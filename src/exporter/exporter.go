package exporter

import (
	"github.com/kahara/go-gnssaggr/src/aggregator"
	"github.com/rs/zerolog/log"
)

func Export(aggregations <-chan aggregator.Aggregation) {
	var (
		aggregation aggregator.Aggregation
	)

	for {
		select {
		case aggregation = <-aggregations:
			log.Printf("%+v", aggregation)
		}
	}
}
