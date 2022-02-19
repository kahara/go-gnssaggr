package exporter

import (
	"fmt"
	"github.com/kahara/go-gnssaggr/src/aggregator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Export(config Config, aggregations <-chan aggregator.Aggregation) {
	var (
		err         error
		addr        = fmt.Sprintf(":%d", config.Port)
		aggregation aggregator.Aggregation
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err = http.ListenAndServe(addr, nil); err != nil {
			log.Fatal().Err(err).Msgf("Could not start listening on %s", addr)
		}
	}()

	for {
		select {
		case aggregation = <-aggregations:
			switch aggregation.(type) {
			case aggregator.TPV:
				log.Printf("tpv %+v", aggregation)
			case aggregator.SKY:
				log.Printf("sky %+v", aggregation)
			}
		}
	}
}
