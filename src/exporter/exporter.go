package exporter

import (
	"fmt"
	"github.com/cespare/xxhash"
	"github.com/kahara/go-gnssaggr/src/aggregator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	Namespace    = "gnssaggr"
	TPVSubsystem = "tpv"
	SKYSubsystem = "sky"
)

var (
	metrics = make(map[uint64]*prometheus.GaugeVec) // Keep metrics in a map keyed with xxHash(subsystem + name)
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
				exportTPV(aggregation.(aggregator.TPV))
			case aggregator.SKY:
				exportSKY(aggregation.(aggregator.SKY))
			}
		}
	}
}

func exportTPV(tpv aggregator.TPV) {
	var (
		key string = strings.ToLower(tpv.Key)
	)

	updateMetric(TPVSubsystem, fmt.Sprintf("%s_min", key), fmt.Sprintf("%s minimum", tpv.Key), []string{}, []string{}, tpv.Min)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_q1", key), fmt.Sprintf("Q1 of %s", tpv.Key), []string{}, []string{}, tpv.Q1)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_q2", key), fmt.Sprintf("Q2 of %s", tpv.Key), []string{}, []string{}, tpv.Q2)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_q3", key), fmt.Sprintf("Q3 of %s", tpv.Key), []string{}, []string{}, tpv.Q3)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_max", key), fmt.Sprintf("%s maximum", tpv.Key), []string{}, []string{}, tpv.Max)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_iqr", key), fmt.Sprintf("Interquartile range of %s", tpv.Key), []string{}, []string{}, tpv.IQR)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_mad", key), fmt.Sprintf("Median absolute deviation of %s", tpv.Key), []string{}, []string{}, tpv.MAD)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_stddev", key), fmt.Sprintf("Standard deviation of %s", tpv.Key), []string{}, []string{}, tpv.StdDev)
	updateMetric(TPVSubsystem, fmt.Sprintf("%s_variance", key), fmt.Sprintf("Variance of %s", tpv.Key), []string{}, []string{}, tpv.Variance)
}

func exportSKY(sky aggregator.SKY) {
	updateMetric(SKYSubsystem, "gdop", "Geometric (hyperspherical) dilution of precision, a combination of PDOP and TDOP", []string{}, []string{}, sky.GDOP)
	updateMetric(SKYSubsystem, "hdop", "Horizontal dilution of precision", []string{}, []string{}, sky.HDOP)
	updateMetric(SKYSubsystem, "pdop", "Position (spherical/3D) dilution of precision", []string{}, []string{}, sky.PDOP)
	updateMetric(SKYSubsystem, "tdop", "Time dilution of precision", []string{}, []string{}, sky.TDOP)
	updateMetric(SKYSubsystem, "vdop", "Vertical (altitude) dilution of precision", []string{}, []string{}, sky.VDOP)
	updateMetric(SKYSubsystem, "xdop", "Longitudinal dilution of precision", []string{}, []string{}, sky.XDOP)
	updateMetric(SKYSubsystem, "ydop", "Latitudinal dilution of precision", []string{}, []string{}, sky.YDOP)

	for gnssid, count := range sky.GNSSID {
		updateMetric(SKYSubsystem, "satellites", "Satellite count", []string{"gnssid"}, []string{fmt.Sprintf("%d", gnssid)}, float64(count))
	}
}

// Create metrics on the fly. Everything is a Gauge. Timestamps aren't fiddled with at this time, maybe later.
func updateMetric(subsystem string, name string, help string, labelNames []string, labelValues []string, value float64) {
	var (
		metricHash uint64 = xxhash.Sum64String(fmt.Sprintf("%s%s", subsystem, name))
		metric     *prometheus.GaugeVec
		ok         bool
	)

	if metric, ok = metrics[metricHash]; !ok {
		log.Printf("Creating metric %s", name)
		metric = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		}, labelNames)
		metrics[metricHash] = metric
	}

	metrics[metricHash].WithLabelValues(labelValues...).Set(value)
}
