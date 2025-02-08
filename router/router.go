package router

import (
	"gobanking/config"
	"gobanking/handler"
	"gobanking/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Setup(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// Setup middleware
	middleware.SetupMiddleware(e, cfg)

	// Setup handlers
	nasabahHandler := handler.NewNasabahHandler(db, cfg)

	// Routes
	e.POST("/daftar", nasabahHandler.Daftar)
	e.POST("/tabung", nasabahHandler.Tabung)
	e.POST("/tarik", nasabahHandler.Tarik)
	e.GET("/saldo/:no_rekening", nasabahHandler.Saldo)
}
