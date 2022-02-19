package reader

import "time"

type Config struct {
	Host string
	Port uint16
}

type Report interface{}

// https://gpsd.gitlab.io/gpsd/gpsd_json.html#_tpv
type TPV struct {
	Class string `json:"class"`

	// Time
	Time        time.Time `json:"time"`
	LeapSeconds float64   `json:"leapseconds"`
	EPT         float64   `json:"ept"`

	// Position
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	AltHAE   float64 `json:"altHAE"`
	AltMSL   float64 `json:"altMSL"`
	GeoidSep float64 `json:"geoidSep"`

	// Position error
	EPX float64 `json:"epx"`
	EPY float64 `json:"epy"`
	EPV float64 `json:"epv"`
}

// https://gpsd.gitlab.io/gpsd/gpsd_json.html#_sky
type SKY struct {
	Class string `json:"class"`

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
	GNSSID int // 0=GPS, 2=Galileo, 3=Beidou, 5=QZSS, 6=GLONASS
	PRN    int // 1-63=GPS, 64-96=GLONASS, 100-164=SBAS
	Used   bool
}
