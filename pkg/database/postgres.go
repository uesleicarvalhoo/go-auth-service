package database

import (
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQLConnection(cfg Config) (*gorm.DB, error) {
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/Sao_Paulo",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
}
