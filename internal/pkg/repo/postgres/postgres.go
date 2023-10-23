package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rip/internal/pkg/repo"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) MigrateAll() error {
	return r.db.AutoMigrate(&repo.DataService{})
}
