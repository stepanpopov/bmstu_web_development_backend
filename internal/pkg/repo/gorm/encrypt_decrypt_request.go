package repo

import (
	"errors"
	"rip/internal/pkg/repo"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: можно ли добавить услугу после формирования???
// PS: ща нельзя

func (r *Repository) CreateEncryptDecryptDraft(creatorID uuid.UUID) (uint, error) {
	request := repo.EncryptDecryptRequest{
		CreatorID:    &creatorID,
		Status:       repo.Draft,
		CreationDate: r.db.NowFunc(),
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.RequestID, nil
}

func (r *Repository) AddDataServiceToDraft(dataID uint, creatorID uuid.UUID) (uint, error) {
	// получаем услугу
	data, err := r.GetDataServiceById(dataID)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("нет такой услуги")
	}
	if !data.Active {
		return 0, errors.New("услуга удалена")
	}

	// получаем черновик
	var draftReq repo.EncryptDecryptRequest
	res := r.db.Where("creator_id = ?", creatorID).Where("status = ?", repo.Draft).Take(&draftReq)

	// создаем черновик, если его нет
	if res.RowsAffected == 0 {
		newDraftRequestID, err := r.CreateEncryptDecryptDraft(creatorID)
		if err != nil {
			return 0, err
		}

		draftReq.RequestID = newDraftRequestID
	}

	// добавляем запись в мм
	requestToData := repo.EncryptDecryptToData{
		DataID:    dataID,
		RequestID: draftReq.RequestID,
	}

	err = r.db.Create(&requestToData).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return 0, errors.New("услуга уже существует в заявке")
		}

		return 0, err
	}

	return draftReq.RequestID, nil
}

func (r *Repository) DeleteDataServiceFromDraft(dataID uint, creatorID uuid.UUID) error {
	// получаем услугу
	data, err := r.GetDataServiceById(dataID)
	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("нет такой услуги")
	}
	if !data.Active {
		return errors.New("услуга удалена")
	}

	// получаем черновик
	draftRequestID, err := r.GetEncryptDecryptDraftID(creatorID)
	if err != nil {
		return err
	}
	if draftRequestID == nil {
		return errors.New("у пользователя нет черновика-заявки")
	}

	// удаляем услугу из черновика
	requestToData := repo.EncryptDecryptToData{
		DataID:    dataID,
		RequestID: *draftRequestID,
	}

	// TODO: если не нашли??
	if err := r.db.Delete(&requestToData).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteDataServiceFromEncryptDecryptRequest(dataID uint, reqID uint) error {
	// получаем услугу
	data, err := r.GetDataServiceById(dataID)
	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("нет такой услуги")
	}
	if data.Active {
		return errors.New("услуга удалена")
	}

	// получаем заявку
	// TODO: проверить заявку

	// удаляем услугу из черновика
	requestToData := repo.EncryptDecryptToData{
		DataID:    data.DataID,
		RequestID: reqID,
	}

	if err := r.db.Delete(&requestToData).Error; err != nil {
		return err
	}

	return nil
}

// returns nil if there is no draft
func (r *Repository) GetEncryptDecryptDraftID(creatorID uuid.UUID) (*uint, error) {
	var draftReq repo.EncryptDecryptRequest
	res := r.db.Where("creator_id = ?", creatorID).Where("status = ?", repo.Draft).Take(&draftReq)

	if errors.Is(gorm.ErrRecordNotFound, res.Error) {
		return nil, nil
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &draftReq.RequestID, nil
}

func (r *Repository) GetEncryptDecryptRequests(status repo.Status, startDate, endDate time.Time, creatorID uuid.UUID, isModerator bool) ([]repo.EncryptDecryptRequestView, error) {
	var requests []repo.EncryptDecryptRequestView

	filterCond := r.db.
		Table("encrypt_decrypt_requests AS e")

	if !isModerator {
		filterCond = filterCond.Where("e.creator_id = ?", creatorID)
	}

	if status != repo.UnknownStatus {
		filterCond = filterCond.Where("e.status = ?", status)
	}

	if !startDate.IsZero() {
		filterCond = filterCond.Where("e.form_date > ?", startDate)
	}

	if !endDate.IsZero() {
		filterCond = filterCond.Where("e.form_date < ?", endDate)
	}

	if err := filterCond.
		Joins("LEFT JOIN users u1 on e.moderator_id = u1.user_id").
		Joins("LEFT JOIN users u2 on e.creator_id = u2.user_id").
		Select([]string{"e.request_id", "e.status", "e.creation_date", "e.finish_date", "e.form_date",
			"u1.username", "u2.username",
		}).
		Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetEncryptDecryptRequestWithDataByID(requestID uint, creatorID uuid.UUID, isModerator bool) (repo.EncryptDecryptRequestView, []repo.DataService, error) {
	if requestID == 0 {
		return repo.EncryptDecryptRequestView{}, nil, errors.New("record not found")
	}

	reqView := repo.EncryptDecryptRequestView{RequestID: requestID}

	filter := r.db.Table("encrypt_decrypt_requests AS e")

	if !isModerator {
		filter = filter.Where("creator_id = ?", creatorID)
	}

	res := filter.
		Joins("LEFT JOIN users u1 on e.moderator_id = u1.user_id").
		Joins("LEFT JOIN users u2 on e.creator_id = u2.user_id").
		Select([]string{"e.request_id", "e.status", "e.creation_date", "e.finish_date", "e.form_date",
			"u1.username", "u2.username",
		}).
		Take(&reqView)

	if err := res.Error; err != nil {
		return repo.EncryptDecryptRequestView{}, nil, err
	}

	var dataService []repo.DataService
	// TODO: test
	res = r.db.
		Table("data_services").
		Where("active = ?", true).
		Joins("JOIN encrypt_decrypt_to_data e on data_services.data_id = e.data_id and e.request_id = ?", requestID).
		Find(&dataService)

	if err := res.Error; err != nil {
		return repo.EncryptDecryptRequestView{}, nil, err
	}

	return reqView, dataService, nil
}

// creator
func (r *Repository) FormEncryptDecryptRequestByID(requestID uint) error {
	var req repo.EncryptDecryptRequest
	res := r.db.
		Where("request_id = ?", requestID).
		Where("status = ?", repo.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = repo.Formed
	now := r.db.NowFunc()
	req.FormDate = &now // наверно не прокнет тк это алиас к time.Now().Local()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteEncryptDecryptRequestByID(requestID uint) error {
	var req repo.EncryptDecryptRequest
	res := r.db.
		Where("request_id = ?", requestID).
		Where("status in (?)", []repo.Status{repo.Draft, repo.Formed}).
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = repo.Deleted
	// надо ли ставить finish date

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

// moderator
func (r *Repository) FinishEncryptDecryptRequestByID(requestID uint, moderatorID uuid.UUID) error {
	return r.finishRejectHelper(repo.Finished, requestID, moderatorID)
}

func (r *Repository) RejectEncryptDecryptRequestByID(requestID uint, moderatorID uuid.UUID) error {
	return r.finishRejectHelper(repo.Rejected, requestID, moderatorID)
}

func (r *Repository) finishRejectHelper(status repo.Status, requestID uint, moderatorID uuid.UUID) error {
	var req repo.EncryptDecryptRequest
	res := r.db.
		Where("request_id = ?", requestID).
		Where("status = ?", repo.Formed).
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.ModeratorID = &moderatorID
	req.Status = status
	/*if req.Status == repo.Finished {
		TODO: добавить результат в мм
	}*/
	now := r.db.NowFunc()
	req.FinishDate = &now // наверно не прокнет тк это алиас к time.Now().Local()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}
