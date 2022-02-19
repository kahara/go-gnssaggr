package aggregator

import "time"

type Aggregation struct {
	Time     time.Time
	Key      string
	Min      float64
	Q1       float64
	Q2       float64
	Q3       float64
	Max      float64
	IQR      float64
	MAD      float64
	StdDev   float64
	Variance float64
}
