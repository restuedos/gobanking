package main

import (
	"fmt"
	"gobanking/config"
	"gobanking/database"
	"gobanking/model"
	"gobanking/router"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		panic("gagal menghubungkan database")
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&model.Nasabah{}); err != nil {
		cfg.Logger.Error("gagal migrasi database", "error", err)
		panic("gagal migrasi database")
	}

	// Echo instance
	e := echo.New()

	// Setup routes
	router.Setup(e, db, cfg)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	cfg.Logger.Info("memulai server", "alamat", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}
