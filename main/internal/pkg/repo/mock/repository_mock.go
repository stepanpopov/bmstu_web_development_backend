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
		{
			DataID:   1,
			DataName: "Encode your secrets",
			Encode:   true,
			Blob:     strings.Repeat("secret ", 10),
		},
		{
			DataID:   2,
			DataName: "Decode your life",
			Encode:   false,
			Blob:     strings.Repeat("0100 ", 20),
		},
		{
			DataID:   3,
			DataName: "Encode your wife",
			Encode:   true,
			Blob:     strings.Repeat("s3cr3t ", 20),
		},
		{
			DataID:   4,
			DataName: "Decode your wife",
			Encode:   false,
			Blob:     strings.Repeat("1030 ", 20),
		},
		{
			DataID:   5,
			DataName: "Encode methan's formauls",
			Encode:   false,
			Blob:     "CH4",
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
		if d.DataID == id {
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
		if strings.Contains(d.DataName, name) {
			filtered = append(filtered, d)
		}
	}

	return filtered, nil
}
