package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetTenderDraftID(moderator_id uint) (*uint, error) {
	var tenderReq ds.Tender

	err := r.db.First(&tenderReq, "moderator_id = ? and status = 'черновик'", moderator_id)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		r.logger.Error("error while getting monitoring request draft", err)
		return nil, err.Error
	}
	return &tenderReq.ID, nil
}

func (r *Repository) CreateTenderDraft(creatorID uint) (uint, error) {
	request := ds.Tender{
		UserID:       creatorID,
		Status:       "черновик",
		CreationDate: r.db.NowFunc(),
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (r *Repository) GetTenderWithDataByID(requestID uint) (ds.Tender, []ds.Company, error) {
	if requestID == 0 {
		return ds.Tender{}, nil, errors.New("record not found")
	}

	request := ds.Tender{ID: requestID}
	res := r.db.Take(&request)
	if err := res.Error; err != nil {
		return ds.Tender{}, nil, err
	}

	var dataService []ds.Company

	res = r.db.
		Table("tender_companies").
		Select("companies.*").
		Where("status != ?", "удалён").
		Joins("JOIN companies ON tender_companies.\"CompanyID\" = companies.id").
		Where("tender_companies.\"TenderID\" = ?", requestID).
		Find(&dataService)

	if err := res.Error; err != nil {
		return ds.Tender{}, nil, err
	}

	return request, dataService, nil
}

func (r *Repository) TenderList(status, start, end string) (*[]ds.Tender, error) {
	var tender []ds.Tender
	query := r.db.Where("status != ? AND status != ?", "удалён", "черновик")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if start != "" {
		query = query.Where("creation_date >= ?", start)
	}

	if end != "" {
		query = query.Where("creation_date <= ?", end)
	}
	result := query.Find(&tender)
	return &tender, result.Error
}

func (r *Repository) UpdateTender(updatedTender *ds.Tender) error {
	oldTender := ds.Tender{}
	if result := r.db.First(&oldTender, updatedTender.ID); result.Error != nil {
		return result.Error
	}
	if updatedTender.Name != "" {
		oldTender.Name = updatedTender.Name
	}
	if updatedTender.CreationDate.String() != utils.EmptyDate {
		oldTender.CreationDate = updatedTender.CreationDate
	}
	if updatedTender.CompletionDate.String() != utils.EmptyDate {
		oldTender.CompletionDate = updatedTender.CompletionDate
	}
	if updatedTender.FormationDate.String() != utils.EmptyDate {
		oldTender.FormationDate = updatedTender.FormationDate
	}
	if updatedTender.Status != "" {
		oldTender.Status = updatedTender.Status
	}

	*updatedTender = oldTender
	result := r.db.Save(updatedTender)
	return result.Error
}

func (r *Repository) FormTenderRequestByID(requestID uint) error {
	var req ds.Tender
	res := r.db.
		Where("id = ?", requestID).
		Where("status = ?", "черновик").
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "сформирован"
	req.CreationDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) RejectTenderRequestByID(requestID, moderatorID uint) error {
	return r.finishRejectHelper("отклонён", requestID, moderatorID)
}

func (r *Repository) FinishEncryptDecryptRequestByID(requestID, moderatorID uint) error {
	return r.finishRejectHelper("завершён", requestID, moderatorID)
}

func (r *Repository) finishRejectHelper(status string, requestID, moderatorID uint) error {
	var req ds.Tender
	res := r.db.
		Where("id = ?", requestID).
		Where("status = ?", "сформирован").
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.ModeratorID = moderatorID
	req.Status = status

	req.CompletionDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteTenderByID(requestID uint) error { // ?
	var req ds.Tender
	res := r.db.
		Where("id = ?", requestID). // ??
		Where("status in (?)", "черновик", "сформирован").
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "удалён"
	req.CompletionDate = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteCompanyFromRequest(userId, companyId uint) (ds.Tender, []ds.Company, error) {
	var request ds.Tender
	r.db.Where("user_id = ? and status = 'сформирован'", userId).First(&request)

	if request.ID == 0 {
		return ds.Tender{}, nil, errors.New("no such request")
	}

	var companyRequesTender ds.TenderCompany
	err := r.db.Where("CompanyID = ? AND TenderID = ?", companyId, request.ID).First(&companyRequesTender).Error
	if err != nil {
		return ds.Tender{}, nil, errors.New("такой компании нет в заявке")
	}

	err = r.db.Where("CompanyID = ? AND TenderID = ?", companyId, request.ID).
		Delete(ds.Tender{}).Error

	if err != nil {
		return ds.Tender{}, nil, err
	}

	return r.GetTenderWithDataByID(request.ID)
}

func (r *Repository) UpdateTenderCompany(tenderID uint, companyID uint, cash float64) error {
	var updateCompany ds.TenderCompany
	r.db.Where(" \"TenderID\" = ? and \"CompanyID\" = ?", tenderID, companyID).First(&updateCompany)

	if updateCompany.TenderID == 0 {
		return errors.New("нет такой заявки")
	}
	updateCompany.Cash = cash

	if err := r.db.Save(&updateCompany).Error; err != nil {
		return err
	}

	return nil
}
