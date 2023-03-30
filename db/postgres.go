package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
