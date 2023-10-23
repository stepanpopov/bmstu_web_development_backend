package repo

type Repository interface {
	GetDataServiceById(uint) (DataService, error)
	GetActiveDataServiceFilteredByName(string) ([]DataService, error)
	CreateDataService(DataService) error
	DeleteDataService(uint) error
	UpdateDataService(*DataService) error
}
