package repository

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Load migration files
	"github.com/pkg/errors"
	"github.com/uesleicarvalhoo/go-auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{})
}

func DBMigrate(dbInstance *gorm.DB, dbName string) error {
	db, err := dbInstance.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "Failed to instantiate postgres driver")
	}

	migrations, err := migrate.NewWithDatabaseInstance("file://internal/infra/repository/migrations", dbName, driver)
	if err != nil {
		return errors.Wrap(err, "Failed to create migrate instance")
	}

	err = migrations.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "Faile to apply migrations up")
	}

	return nil
}
