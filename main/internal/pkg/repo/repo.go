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

	CreateEncryptDecryptDraft(creatorID uuid.UUID) (uint, error)
	AddDataServiceToDraft(dataID uint, creatorID uuid.UUID) (uint, error)
	DeleteDataServiceFromDraft(dataID uint, creatorID uuid.UUID) error
	GetEncryptDecryptDraftID(creatorID uuid.UUID) (*uint, error)
	GetEncryptDecryptRequests(status Status, startDate, endDate time.Time, creatorID uuid.UUID, isModerator bool) ([]EncryptDecryptRequestView, error)
	GetEncryptDecryptRequestWithDataByID(requestID uint, creatorID uuid.UUID, isModerator bool) (EncryptDecryptRequestView, []DataServiceWithOptResult, error)
	FormEncryptDecryptRequestByID(requestID uint) error
	DeleteEncryptDecryptRequestByID(requestID uint) error
	FinishEncryptDecryptRequestByID(requestID uint, moderatorID uuid.UUID) error
	RejectEncryptDecryptRequestByID(requestID uint, moderatorID uuid.UUID) error
	DeleteDataServiceFromEncryptDecryptRequest(dataID uint, reqID uint) error
	UpdateCalculated(reqID uint, calculated []Calculated) error

	CreateUser(username, passwordHash string, isModerator bool) (uuid.UUID, error)
	CheckUser(username, passwordHash string) (uuid.UUID, bool, error)
}

type Avatar interface {
	Put(ctx context.Context, avatar io.Reader, size int64) (uuid.UUID, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}
