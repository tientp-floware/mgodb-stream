package public

import (
	"github.com/labstack/echo/v4"
	repository "github.com/tientp-floware/mgodb-stream/repositories"
)

// Route start route public
func Route(e *echo.Echo) {
	group := e.Group("/api/web")
	srv := New(repository.New())
	group.GET("user/setting/:q", srv.Device.ByPlate)
}
