package repo

import (
	"github.com/labstack/echo/v4"

	"github.com/tientp-floware/mgodb-stream/config"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	model "github.com/tientp-floware/mgodb-stream/models"
)

var (
	dbgps  = mgodb.NewSwitchMogoDB(config.Config.Mongodb.DatabaseIOT, db.GPSTrackingCollection)
	dbtrip = mgodb.NewMogoDB(db.TripCollection)
)

// TripCRUD represents the client
type TripCRUD struct {
}

// GPSTrackingCRUD represents the client
type GPSTrackingCRUD struct {
}

// NewGPSCRUD crud
func NewGPSCRUD() *GPSTrackingCRUD {
	return &GPSTrackingCRUD{}
}

// NewTripCRUD crud
func NewTripCRUD() *TripCRUD {
	return &TripCRUD{}
}

// Get location tracking
func (gps *GPSTrackingCRUD) Get(trip string, result *model.TripLocationTracking) error {
	return dbgps.FindOne(echo.Map{"trip": trip}).Decode(&result)
}

// Update location tracking
func (gps *GPSTrackingCRUD) Update(trip string, data interface{}) error {
	updated, err := dbgps.UpdateOne(echo.Map{"trip": trip}, echo.Map{"$set": data})
	if err != nil {
		log.Error("Can not updated:", err)
		return err
	}
	log.Info("[Tracking] MatchedCount:", updated.MatchedCount, " ModifiedCount:", updated.ModifiedCount)
	return err
}

// Distance location tracking
func (gps *GPSTrackingCRUD) Distance(code string, result *model.TripLocationTracking) *GPSTrackingCRUD {
	return gps
}

// Update location tracking
func (gps *TripCRUD) Update(trip string, data interface{}) error {
	updated, err := dbtrip.UpdateOne(echo.Map{"code": trip}, echo.Map{"$set": data})
	if err != nil {
		log.Error("Can not updated:", err)
		return err
	}
	log.Info("[Trip] MatchedCount:", updated.MatchedCount, " ModifiedCount:", updated.ModifiedCount)
	return err
}
