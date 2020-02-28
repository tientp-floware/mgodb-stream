package repo

import (
	"g.ghn.vn/go-common/polyline"
	logger "g.ghn.vn/go-common/zap-logger"
	model "github.com/tientp-floware/mgodb-stream/models"
)

var (
	log = logger.GetLogger("[Trip service]")
)

type (
	// Trip repo
	Trip struct {
		poly polyline.Codec
	}
)

// NewTrip create intance
func NewTrip() *Trip {
	p := polyline.Codec{Dim: 2, Scale: 1e5}
	return &Trip{p}
}

// GPSTrackingByTrip by trip code
func (tr *Trip) GPSTrackingByTrip(trip string) *model.TripLocationTracking {
	var gpsTracking model.TripLocationTracking
	err := NewGPSCRUD().Get(trip, &gpsTracking)
	if err != nil {
		log.Info("err trip:", err)
	}
	return &gpsTracking
}

// GPSCurrentLocation update current location
func (tr *Trip) GPSCurrentLocation(id string, updated interface{}) error {
	return NewGPSCRUD().Update(id, updated)
}

// TripCurrentLocation update current location
func (tr *Trip) TripCurrentLocation(id string, updated interface{}) error {
	return NewTripCRUD().Update(id, updated)
}

// Distance calculate trip len
func (tr *Trip) Distance(trip string) (float64, float64) {
	result := tr.GPSTrackingByTrip(trip)
	bufGps := []byte(result.GPSPolyline)
	bufApp := []byte(result.AppPolyline)
	lenGps, _, _ := polyline.DecodeCoords(bufGps)
	lenApp, _, _ := polyline.DecodeCoords(bufApp)
	return polyline.Distance(lenGps), polyline.Distance(lenApp)
}
