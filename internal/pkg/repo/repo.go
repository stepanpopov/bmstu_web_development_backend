package repo

import "time"

type Repository interface {
	GetDataServiceById(uint) (*DataService, error)
	GetActiveDataServiceFilteredByName(string) ([]DataService, error)
	CreateDataService(DataService) error
	DeleteDataService(uint) error
	UpdateDataService(*DataService) error

	CreateEncryptDecryptDraft(creatorID uint) (uint, error)
	AddDataServiceToDraft(dataID uint, creatorID uint) (uint, error)
	DeleteDataServiceFromDraft(dataID uint, creatorID uint) error
	GetEncryptDecryptDraftID(creatorID uint) (*uint, error)
	GetEncryptDecryptRequests(status Status, startDate, endDate time.Time) ([]EncryptDecryptRequest, error)
	GetEncryptDecryptRequestWithDataByID(requestID uint) (EncryptDecryptRequest, []DataService, error)
	FormEncryptDecryptRequestByID(requestID, creatorID uint) error
	DeleteEncryptDecryptRequestByID(requestID, creatorID uint) error
	FinishEncryptDecryptRequestByID(requestID, moderatorID uint) error
	RejectEncryptDecryptRequestByID(requestID, moderatorID uint) error
}
