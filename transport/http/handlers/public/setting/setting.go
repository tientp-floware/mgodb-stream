package setting

import (
	"github.com/labstack/echo/v4"
	repository "github.com/tientp-floware/mgodb-stream/repositories"
)

type (
	// HTTP in package device
	HTTP struct {
		srv *repository.Service
	}
)

// New create new device handler
func New(srv *repository.Service) *HTTP {
	return &HTTP{srv}
}

//ByPlate handle get device by vehicle plate
func (h *HTTP) ByPlate(c echo.Context) error {
	qr := c.Param("q")
	result := h.srv.Device.ByVehiclePlate(qr)
	return c.JSON(200, result)
}
