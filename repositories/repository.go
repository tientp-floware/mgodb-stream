package repository

import (
	model "github.com/tientp-floware/mgodb-stream/models"
	setting "github.com/tientp-floware/mgodb-stream/repositories/setting"
	tracking "github.com/tientp-floware/mgodb-stream/repositories/tracking"
	trip "github.com/tientp-floware/mgodb-stream/repositories/trip"
)

type (
	// Service instance
	Service struct {
		Setting  Setting
		GPS      GPS
		Trip     Trip
		Tracking Tracking
	}
	// Setting list function can use
	Setting interface {
	}
	// GPS hold func's
	GPS interface {
	}
	// Trip hold func's
	Trip interface {
		GPSTrackingByTrip(string) *model.TripLocationTracking
		//GPSCurrentLocation(string, interface{}) error
		TripCurrentLocation(string, interface{}) error
		Distance(string) (float64, float64)
	}
	// Tracking GPS
	Tracking interface {
		Trip(string) *tracking.Tracking
	}
)

// New creates new user application service
func New() *Service {
	return &Service{Setting: setting.NewDevice(), Trip: trip.NewTrip(), Tracking: tracking.NewTracking()}
}
