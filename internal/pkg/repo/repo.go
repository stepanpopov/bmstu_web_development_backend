package repo

type Repository interface {
	GetDataServiceById(uint) (DataService, error)
	GetDataServiceAll() ([]DataService, error)
	GetDataServiceFilteredByName(string) ([]DataService, error)
}
