package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB initialize postgres db connection
// and return gorm.DB instance without any modification
func NewPostgresDB(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
