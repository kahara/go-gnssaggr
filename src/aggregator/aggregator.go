package aggregator

import (
	"github.com/kahara/go-gnssaggr/src/reader"
	"github.com/rs/zerolog/log"
)

const (
	// For easy times without many allocations, store incoming TPV reports in an array indexed by second of minute.
	period = 61 // Include possible positive https://en.wikipedia.org/wiki/Leap_second
)

func Aggregate(reports <-chan reader.Report, aggregates chan<- Aggregation) {
	var (
		report reader.Report
		tpvs   = make(chan reader.TPV, 2)
		skys   = make(chan reader.SKY, 2)
	)

	go aggregateTPVs(tpvs)
	go aggregateSKYs(skys)

	for {
		select {
		case report = <-reports:
			switch report.(type) {
			case reader.TPV:
				tpvs <- report.(reader.TPV)
			case reader.SKY:
				skys <- report.(reader.SKY)
			}
		}
	}
}

// TPV (time-position-velocity) reports are gathered for the ongoing minute, and are massaged when the incoming
// report's timestamp has zero seconds.
//
// Note that we're relying on the report's timestamps and system clock isn't involved.
func aggregateTPVs(tpvs <-chan reader.TPV) {
	var (
	//tpv [period]reader.TPV
	)

	for {
		select {
		case tpv := <-tpvs:
			log.Print(tpv)
		}
	}
}

// SKY (sky view) reports are massaged and passed to exporter as they come in, without collecting.
func aggregateSKYs(skys <-chan reader.SKY) {
	for {
		select {
		case sky := <-skys:
			log.Print(sky)
		}
	}
}
