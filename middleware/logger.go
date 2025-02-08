package middleware

import (
	"gobanking/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
