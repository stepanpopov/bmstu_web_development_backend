package repo

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetDataServiceById(uint) (*DataService, error)
	GetActiveDataServiceFilteredByName(string) ([]DataService, error)
	CreateDataService(DataService) (uint, error)
	DeleteDataService(uint) (uuid.UUID, error)
	UpdateDataService(*DataService) error
	UpdateImageUUID(imageUUID uuid.UUID, dataID uint) error

	CreateEncryptDecryptDraft(creatorID uint) (uint, error)
	AddDataServiceToDraft(dataID uint, creatorID uint) (uint, error)
	DeleteDataServiceFromDraft(dataID uint, creatorID uint) error
	GetEncryptDecryptDraftID(creatorID uint) (*uint, error)
	GetEncryptDecryptRequests(status Status, startDate, endDate time.Time) ([]EncryptDecryptRequestView, error)
	GetEncryptDecryptRequestWithDataByID(requestID uint) (EncryptDecryptRequestView, []DataService, error)
	FormEncryptDecryptRequestByID(requestID uint) error
	DeleteEncryptDecryptRequestByID(requestID uint) error
	FinishEncryptDecryptRequestByID(requestID, moderatorID uint) error
	RejectEncryptDecryptRequestByID(requestID, moderatorID uint) error
	DeleteDataServiceFromEncryptDecryptRequest(dataID uint, reqID uint) error
}

type Avatar interface {
	Put(ctx context.Context, avatar io.Reader, size int64) (uuid.UUID, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}
