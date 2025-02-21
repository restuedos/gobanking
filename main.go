package main

import (
	"fmt"
	"gobanking/config"
	"gobanking/database"
	_ "gobanking/docs"
	"gobanking/model"
	"gobanking/router"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Banking API
// @version 1.0
// @description This is a banking service API with authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer authentication
func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		panic("gagal menghubungkan database")
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&model.Nasabah{}, &model.User{}); err != nil {
		cfg.Logger.Error("gagal migrasi database", "error", err)
		panic("gagal migrasi database")
	}

	// Echo instance
	e := echo.New()

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Setup routes
	router.Setup(e, db, cfg)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	cfg.Logger.Info("memulai server", "alamat", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}
