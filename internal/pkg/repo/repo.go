package repo

type Repository interface {
	GetDataServiceById(uint) (DataService, error)
	GetDataServiceAll() ([]DataService, error)
	GetActiveDataServiceFilteredByName(string) ([]DataService, error)
	DeleteDataService(uint) error
}
