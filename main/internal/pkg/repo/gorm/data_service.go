package repo

import (
	"fmt"
	"rip/internal/pkg/repo"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// nil если такой заявки нет
func (r *Repository) GetDataServiceById(id uint) (*repo.DataService, error) {
	print(id)
	dataService := repo.DataService{DataID: id}
	fmt.Print(dataService)
	res := r.db.Where("active = ?", true).Take(&dataService)
	fmt.Print(dataService)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &dataService, nil
}

func (r *Repository) GetActiveDataServiceFilteredByName(name string) ([]repo.DataService, error) {
	name = strings.ToLower(name)
	var dataService []repo.DataService
	if err := r.db.Where(&repo.DataService{Active: true}).Where("LOWER(data_name) LIKE ?", "%"+name+"%").Find(&dataService).Error; err != nil {
		return nil, err
	}
	fmt.Printf("%v", dataService)

	return dataService, nil
}

func (r *Repository) UpdateDataService(d *repo.DataService) error {
	err := r.db.Model(&d).Where("active = ?", true).Updates(map[string]interface{}{"data_name": d.DataName, "encode": d.Encode, "blob": d.Blob}).Error
	return err
}

func (r *Repository) UpdateImageUUID(imageUUID uuid.UUID, dataID uint) error {
	return r.db.Model(&repo.DataService{DataID: dataID}).Where("active = ?", true).Updates(map[string]interface{}{"image_uuid": imageUUID}).Error
}

func (r *Repository) CreateDataService(d repo.DataService) (uint, error) {
	d.Active = true
	print(d.DataID)
	return d.DataID, r.db.Create(&d).Error
}

func (r *Repository) DeleteDataService(id uint) (uuid.UUID, error) {

	dataS := &repo.DataService{}
	err := r.db.Model(dataS).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "image_uuid"}}}).
		Where("data_id = ?", id).
		Update("active", false).
		Error

	// tx := r.db.Exec("UPDATE data_services SET active = false WHERE data_id = ? RETURNING image_uuid", id)

	if err != nil {
		return uuid.Nil, err
	}
	return dataS.ImageUUID, err
}

/* TODO
func (r *Repository) UpdateImage(dataID uint) error {

} */
