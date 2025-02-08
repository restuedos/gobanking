package database

import (
	"gobanking/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	cfg.Logger.Info("menghubungkan database",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		cfg.Logger.Error("gagal menghubungkan database", "error", err)
		return nil, err
	}

	return db, nil
}
