package repo

import (
	"time"

	"g.ghn.vn/go-common/polyline"
	"github.com/labstack/echo/v4"
	"github.com/tientp-floware/mgodb-stream/config"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	model "github.com/tientp-floware/mgodb-stream/models"
	trip "github.com/tientp-floware/mgodb-stream/repositories/trip"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collectionGPS = mgodb.NewSwitchMogoDB(config.Config.Mongodb.Database, db.Setting)
)

type (
	// TrackingCRUD store data
	TrackingCRUD struct {
		Polyline Polyline
		PointApp [][]float64
		PointGPS [][]float64
		Data     *model.Tracking
		trip     Trip
	}
	// Polyline method
	Polyline interface {
		EncodeCoords([]byte, [][]float64) []byte
		DecodeCoords([]byte) ([][]float64, []byte, error)
	}
	Trip interface {
		GPSTrackingByTrip(string) *model.TripLocationTracking
	}
)

// NewTrackingCRUD collection
func NewTrackingCRUD() *TrackingCRUD {
	p := polyline.Codec{Dim: 2, Scale: 1e5}
	tr := &TrackingCRUD{Polyline: p, trip: trip.NewTrip()}

	return tr
}

// AddPointGPS add new point
func (tr *TrackingCRUD) AddPointGPS(point []float64, t time.Time) *TrackingCRUD {
	tr.PointGPS = append(tr.PointGPS, point)
	tr.Data.TimePointGPS = append(tr.Data.TimePointGPS, t)
	return tr
}

// AddPointApp add new point
func (tr *TrackingCRUD) AddPointApp(point []float64, t time.Time) *TrackingCRUD {
	tr.PointApp = append(tr.PointApp, point)
	tr.Data.TimePointApp = append(tr.Data.TimePointApp, t)
	return tr
}

// EncodePolyline add new point
func (tr *TrackingCRUD) EncodePolyline() *TrackingCRUD {
	tr.Data.PolylineGPS = string(tr.Polyline.EncodeCoords(nil, tr.PointGPS))
	tr.Data.PolylineApp = string(tr.Polyline.EncodeCoords(nil, tr.PointApp))
	return tr
}

// DecodePolylineGPS add new point
func (tr *TrackingCRUD) DecodePolylineGPS() *TrackingCRUD {
	tr.PointGPS, _, _ = tr.Polyline.DecodeCoords([]byte(tr.Data.PolylineGPS))
	return tr
}

// DecodePolylineApp add new point
func (tr *TrackingCRUD) DecodePolylineApp() *TrackingCRUD {
	tr.PointApp, _, _ = tr.Polyline.DecodeCoords([]byte(tr.Data.PolylineApp))
	return tr
}

// Find trip
func (tr *TrackingCRUD) Find(code string) error {
	trip := tr.trip.GPSTrackingByTrip(code)
	tr.Data = &model.Tracking{
		Trip: code,
	}
	if trip != nil {
		tr.Data = &model.Tracking{
			ID:           trip.ID,
			Trip:         trip.Trip,
			PolylineApp:  trip.AppPolyline,
			PolylineGPS:  trip.GPSPolyline,
			TimePointApp: trip.TimePointApp,
			TimePointGPS: trip.TimePointGPS,
			CreatedAt:    trip.CreatedAt,
			UpdatedAt:    trip.UpdatedAt,
		}
	}
	return nil
}

// Create trip
func (tr *TrackingCRUD) Create() error {
	tr.Data.CreatedAt = time.Now()
	tr.Data.UpdatedAt = time.Now()
	tr.Data.ID = primitive.NewObjectID()
	_, err := collectionGPS.InsertOne(tr.Data)
	if err != nil {
		log.Error("Can not insert:", err)
		return err
	}
	return nil
}

// Update  trip
func (tr *TrackingCRUD) Update(code string) error {
	matched, err := collectionGPS.UpdateOne(echo.Map{"trip": code}, echo.Map{"$set": echo.Map{
		"polyline_app":   tr.Data.PolylineApp,
		"polyline_gps":   tr.Data.PolylineGPS,
		"time_point_app": tr.Data.TimePointApp,
		"time_point_gps": tr.Data.TimePointGPS,
		"updated_at":     time.Now(),
	}})
	if err != nil {
		log.Error("Can not updated:", err)
		return err
	}
	log.Info("[Tracking] MatchedCount:", matched.MatchedCount, " ModifiedCount:", matched.ModifiedCount)
	return err
}
