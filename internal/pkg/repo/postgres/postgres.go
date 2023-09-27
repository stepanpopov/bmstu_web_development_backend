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

func (r *Repository) GetDataServiceById(id uint) (repo.DataService, error) {
	return repo.DataService{}, nil
}

func (r *Repository) GetDataServiceAll() ([]repo.DataService, error) {
	return nil, nil
}

func (r *Repository) GetDataServiceFilteredByName(name string) ([]repo.DataService, error) {
	return nil, nil
}

/* func (r *Repository) GetProductByID(id int) (*ds.Product, error) {
	product := &ds.Product{}

	err := r.db.First(product, "id = ?", "1").Error // find product with id = 1
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) CreateProduct(product ds.Product) error {
	return r.db.Create(product).Error
} */
