package gps

import (
	logger "g.ghn.vn/go-common/zap-logger"
)

// Service represents the client
type Service struct{}

var (
	log = logger.GetLogger("[GPS service]")
)

func init() {
	log.Info("GPS init service")
}

// New returns a new Fuel database instance
func New() *Service {
	return &Service{}
}

// ToServiceGPS get last info
func (s *Service) ToServiceGPS() {

}
