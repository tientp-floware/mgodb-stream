package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tientp-floware/mgodb-stream/transport/http/handlers/public"
	logger "go.uber.org/zap"
)

var (
	log = logger.GetLogger("Transport HTTP")
)

type (
	// HTTP name transport
	HTTP struct {
	}
)

// NewHTTP transport
func NewHTTP() *HTTP {
	return &HTTP{}
}

// Server run HTTP transport
func (h *HTTP) Server() *echo.Echo {

	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, "")
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.GET("/health", HealthCheck)
	// Start handler
	// public api
	public.Route(e)
	return e
}

// HealthCheck check status service
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
