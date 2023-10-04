package mock

import (
	"errors"
	"rip/internal/pkg/repo"
	"strings"
)

type Repository struct {
	data []repo.DataService
}

func New() *Repository {
	data := []repo.DataService{
		repo.DataService{
			ID:            1,
			Name:          "Encode your secrets",
			Encode_Decode: repo.Encode,
			Blob:          strings.Repeat("secret", 10),
		},
		repo.DataService{
			ID:            2,
			Name:          "Decode your life",
			Encode_Decode: repo.Decode,
			Blob:          strings.Repeat("0100", 20),
		},
	}

	return &Repository{
		data: data,
	}
}

var (
	ErrNotFound = errors.New("Not Found")
)

func (r *Repository) GetDataServiceById(id uint) (repo.DataService, error) {
	for _, d := range r.data {
		if d.ID == id {
			d := d
			return d, nil
		}
	}

	return repo.DataService{}, ErrNotFound
}

func (r *Repository) GetDataServiceAll() ([]repo.DataService, error) {
	return r.data, nil
}

func (r *Repository) GetDataServiceFilteredByName(name string) ([]repo.DataService, error) {

	var filtered []repo.DataService
	for _, d := range r.data {
		if strings.Contains(d.Name, name) {
			filtered = append(filtered, d)
		}
	}

	return filtered, nil
}
