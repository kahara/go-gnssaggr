package aggregator

import "time"

type Aggregation interface{}

type TPV struct {
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

type SKY struct {
	Time   time.Time
	GDOP   float64
	HDOP   float64
	PDOP   float64
	TDOP   float64
	VDOP   float64
	XDOP   float64
	YDOP   float64
	GNSSID map[int]int
}
