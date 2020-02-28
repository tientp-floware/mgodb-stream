package public

import (
	"github.com/labstack/echo/v4"
	repository "github.com/tientp-floware/mgodb-stream/repositories"
	"github.com/tientp-floware/mgodb-stream/transport/http/handlers/public/setting"
)

type (
	// Service Public package
	Service struct {
		Setting Setting
	}
	// Setting handler method
	Setting interface {
		ByPlate(echo.Context) error
	}
)

// New add public services
func New(srv *repository.Service) *Service {
	return &Service{
		Setting: setting.New(srv),
	}
}
