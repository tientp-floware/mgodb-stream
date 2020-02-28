package repo

import (
	"time"

	logger "go.uber.org/zap"
)

var (
	log = logger.GetLogger("[Tracking service]")
)

type (
	// Tracking store data
	Tracking struct {
		Code string
		Poly *TrackingCRUD
	}
)

// NewTracking collection
func NewTracking() *Tracking {
	return &Tracking{Poly: NewTrackingCRUD()}
}

// Trip trip
func (tr *Tracking) Trip(code string) *Tracking {
	tr.Code = code
	return tr.getrip()
}

func (tr *Tracking) getrip() *Tracking {
	err := tr.Poly.Find(tr.Code)
	if tr.Poly.Data.ID.IsZero() {
		log.Info("not found:", tr.Code, " - err:", err)
		tr.Poly.Data.Trip = tr.Code
		tr.Poly.Create()
		return tr
	}
	return tr
}

// DecodePolylineGPS in trip
func (tr *Tracking) DecodePolylineGPS() *Tracking {
	tr.Poly.DecodePolylineGPS()
	return tr
}

// DecodePolylineApp in trip
func (tr *Tracking) DecodePolylineApp() *Tracking {
	tr.Poly.DecodePolylineApp()
	return tr
}

// EncodePolyline in trip
func (tr *Tracking) EncodePolyline() *Tracking {
	tr.Poly.EncodePolyline()
	return tr.update()
}

// AddPointAPP in trip
func (tr *Tracking) AddPointAPP(point []float64, t time.Time) *Tracking {
	tr.Poly.AddPointApp(point, t)
	return tr
}

// AddPointGPS in trip
func (tr *Tracking) AddPointGPS(point []float64, t time.Time) *Tracking {
	tr.Poly.AddPointGPS(point, t)
	return tr
}

func (tr *Tracking) update() *Tracking {
	tr.Poly.Update(tr.Code)
	return tr
}
