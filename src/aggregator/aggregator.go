package aggregator

import (
	"github.com/kahara/go-gnssaggr/src/reader"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
	"reflect"
)

const (
	// For easy times without many allocations, store incoming TPV reports in an array indexed by second of minute.
	period = 60 // Not taking possible positive https://en.wikipedia.org/wiki/Leap_second into account, skip it instead
)

func Aggregate(reports <-chan reader.Report, aggregates chan<- Aggregation) {
	var (
		report reader.Report
		tpvs   = make(chan reader.TPV, 2)
		skys   = make(chan reader.SKY, 2)
	)

	go aggregateTPVs(tpvs, aggregates)
	go aggregateSKYs(skys, aggregates)

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
func aggregateTPVs(tpvs <-chan reader.TPV, aggregates chan<- Aggregation) {
	var (
		tpv        [period]reader.TPV
		second     int
		prevSecond int = -1
	)

	for {
		select {
		case _tpv := <-tpvs:
			second = _tpv.Time.Second()
			switch {
			case second == prevSecond:
				// Skip duplicate reports for the same second
				continue
			case second > 59:
				// Skip positive leap second
				continue
			case second == 59:
				// Aggregate reports for current minute
				tpv[second] = _tpv

				// Check that there's a report for each second of current minute
				for _, x := range tpv {
					if reflect.DeepEqual(x, reader.TPV{}) {
						log.Printf("Incomplete set of TPV reports for this minute, skipping")
						goto cleanup
					}
				}

				for _, key := range []string{"EPT", "Lon", "Lat", "AltHAE", "AltMSL", "EPX", "EPY", "EPV"} {
					aggregates <- func(key string) TPV {
						var (
							aggregation TPV = TPV{
								Time: tpv[0].Time,
								Key:  key,
							}
							values    []float64
							value     reflect.Value
							quartiles stats.Quartiles
						)

						for _, x := range tpv {
							value = reflect.ValueOf(x)
							values = append(values, reflect.Indirect(value).FieldByName(key).Float())
						}

						aggregation.Min, _ = stats.Min(values)
						quartiles, _ = stats.Quartile(values)
						aggregation.Q1 = quartiles.Q1
						aggregation.Q2 = quartiles.Q2
						aggregation.Q3 = quartiles.Q3
						aggregation.Max, _ = stats.Max(values)
						aggregation.IQR, _ = stats.InterQuartileRange(values)
						aggregation.MAD, _ = stats.MedianAbsoluteDeviation(values)
						aggregation.StdDev, _ = stats.StdDevP(values)
						aggregation.Variance, _ = stats.Variance(values)

						return aggregation
					}(key)
				}

				// Clean the slate for new one-minute period
			cleanup:
				tpv = [period]reader.TPV{}
			default:
				tpv[second] = _tpv
			}
			prevSecond = second
		}
	}
}

// SKY (sky view) reports are massaged and passed to exporter as they come in, without collecting.
func aggregateSKYs(skys <-chan reader.SKY, aggregates chan<- Aggregation) {
	var (
		aggregation SKY
	)

	for {
		select {
		case sky := <-skys:
			aggregation = SKY{
				Time:   sky.Time,
				GDOP:   sky.GDOP,
				HDOP:   sky.HDOP,
				PDOP:   sky.PDOP,
				TDOP:   sky.TDOP,
				VDOP:   sky.VDOP,
				XDOP:   sky.XDOP,
				YDOP:   sky.YDOP,
				GNSSID: map[int]int{},
			}

			for _, satellite := range sky.Satellites {
				aggregation.GNSSID[satellite.GNSSID] += 1
			}

			aggregates <- aggregation
		}
	}
}
