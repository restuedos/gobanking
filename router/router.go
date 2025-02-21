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
	
	// Auth routes
	authHandler := handler.NewAuthHandler(db, cfg)
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected routes
	nasabahHandler := handler.NewNasabahHandler(db, cfg)
	
	// Create a group for protected routes
	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("/daftar", nasabahHandler.Daftar)
	protected.POST("/tabung", nasabahHandler.Tabung)
	protected.POST("/tarik", nasabahHandler.Tarik)
	protected.GET("/saldo/:no_rekening", nasabahHandler.Saldo)
}
