package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rip/internal/pkg/repo"
)

type Repository struct {
	db *gorm.DB
}

func NewPostgres(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.Logger.LogMode(logger.Info)

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) MigrateAll() error {
	return r.db.AutoMigrate(&repo.DataService{})
}
