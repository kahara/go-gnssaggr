package reader

import "time"

type Config struct {
	Host string
	Port uint16
}

type Report interface{}

// https://gpsd.gitlab.io/gpsd/gpsd_json.html#_tpv
type TPV struct {
	Class string `json:"class"` // Here, "TPV"

	// Time
	Time        time.Time `json:"time"`
	LeapSeconds float64   `json:"leapseconds"`
	EPT         float64   `json:"ept"` // Estimated time stamp error in seconds

	// Position
	Lon    float64 `json:"lon"`    // Longitude in degrees: +/- signifies East/West
	Lat    float64 `json:"lat"`    // Latitude in degrees: +/- signifies North/South
	AltHAE float64 `json:"altHAE"` // Altitude, height above ellipsoid, in meters. Probably WGS84.
	AltMSL float64 `json:"altMSL"` // MSL Altitude in meters. The geoid used is rarely specified and is often inaccurate.

	// Position error
	EPX float64 `json:"epx"` // Longitude error estimate in meters
	EPY float64 `json:"epy"` // Latitude error estimate in meters
	EPV float64 `json:"epv"` // Estimated vertical error in meters
}

// https://gpsd.gitlab.io/gpsd/gpsd_json.html#_sky
type SKY struct {
	Class string `json:"class"` // Here, "SKY"

	// The DOPs are "dimensionless factors which should be multiplied by a base UERE to get an error estimate"
	GDOP float64 `json:"gdop"` // Geometric (hyperspherical) dilution of precision, a combination of PDOP and TDOP
	HDOP float64 `json:"hdop"` // Horizontal dilution of precision
	PDOP float64 `json:"pdop"` // Position (spherical/3D) dilution of precision
	TDOP float64 `json:"tdop"` // Time dilution of precision
	VDOP float64 `json:"vdop"` // Vertical (altitude) dilution of precision
	XDOP float64 `json:"xdop"` // Longitudinal dilution of precision
	YDOP float64 `json:"ydop"` // Latitudinal dilution of precision

	Satellites []Satellite
}

type Satellite struct {
	GNSSID int  `json:"gnssid"` // 0=GPS, 2=Galileo, 3=Beidou, 5=QZSS, 6=GLONASS
	PRN    int  `json:"PRN"`    // 1-63=GPS, 64-96=GLONASS, 100-164=SBAS
	Used   bool `json:"used"`   // This satellite is used in current solution
}
