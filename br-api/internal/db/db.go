package db

import (
	"github.com/melfish/br-api/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.Manager{},
		&models.Technician{},
		&models.Quote{},
		&models.Job{},
		&models.Notification{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
