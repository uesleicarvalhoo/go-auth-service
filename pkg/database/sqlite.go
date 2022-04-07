package database

import (
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteMemoryConnection() (*gorm.DB, error) {
	dsn := "file:memdb1?mode=memory&cached=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
}
