package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetTenderDraftID(creatorID int) (uint, error) {
	var draftReq ds.Tender

	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", utils.Draft).Take(&draftReq)
	if errors.Is(gorm.ErrRecordNotFound, res.Error) {
		return 0, nil
	}

	if res.Error != nil {
		return 0, res.Error
	}

	return draftReq.ID, nil
}

func (r *Repository) CreateTenderDraft(creatorID uint) (uint, error) {
	request := ds.Tender{
		// ModeratorID:  creatorID, // просто заглушка, потом придумаю, как сделать норм
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

//func (r *Repository) GetUsersLoginForRequests(tenderRequests []ds.Tender) ([]ds.Tender, error) {
//	for i := range tenderRequests {
//		var user ds.User
//		r.db.Select("login").Where("user_id = ?", tenderRequests[i].UserID).First(&user)
//		tenderRequests[i].UserID = user.Login
//		fmt.Println(monitoringRequests[i].Creator)
//	}
//	return monitoringRequests, nil
//}

func (r *Repository) TenderList(status string, startDate, endDate time.Time, userId int, isAdmin bool) ([]ds.Tender, error) {
	var tender []ds.Tender
	ending := "AND user_id = " + strconv.Itoa(userId)
	if isAdmin {
		ending = ""
	}

	if status != "" {
		if startDate.IsZero() {
			if endDate.IsZero() {
				// фильтрация только по статусу
				res := r.db.Where("status = ? AND status != 'удален'"+ending, status).Find(&tender)
				//tenderRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
				return tender, res.Error
			}

			// фильтрация по статусу и endDate
			res := r.db.Where("status = ? AND status != 'удален'"+ending, status).Where("creation_date < ?", endDate).
				Find(&tender)
			//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
			return tender, res.Error
		}

		// фильтрация по статусу и startDate
		if endDate.IsZero() {
			res := r.db.Where("status = ? AND status != 'удален'"+ending, status).Where("creation_date > ?", startDate).
				Find(&tender)
			//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
			return tender, res.Error
		}

		// фильтрация по статусу, startDate и endDate
		res := r.db.Where("status = ? AND status != 'удален'"+ending, status).Where("creation_date BETWEEN ? AND ?", startDate, endDate).
			Find(&tender)
		//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
		return tender, res.Error
	}

	if startDate.IsZero() {
		if endDate.IsZero() {
			// без фильтрации
			res := r.db.Where("status <> ?"+ending, "удален").Find(&tender)
			//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
			return tender, res.Error
		}

		// фильтрация по endDate
		res := r.db.Where("creation_date < ?"+ending, endDate).
			Find(&tender)
		//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
		return tender, res.Error
	}

	if endDate.IsZero() {
		// фильтрация по startDate
		res := r.db.Where("creation_date > ?"+ending, startDate).
			Find(&tender)
		//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)
		return tender, res.Error
	}

	//фильтрация по startDate и endDate
	res := r.db.Where("creation_date BETWEEN ? AND ?"+ending, startDate, endDate).
		Find(&tender)
	//monitoringRequests, _ = r.GetUsersLoginForRequests(monitoringRequests)

	return tender, res.Error

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

func (r *Repository) FormTenderRequestByID(requestID uint, creatorID uint) error {
	var req ds.Tender
	res := r.db.
		Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "сформирован"
	req.FormationDate = time.Now()

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
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.Status = "удалён"
	delTime := time.Now()
	req.CompletionDate = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	if err := r.db.Model(&req).Update("deleted_at", delTime).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteCompanyFromRequest(deleteFromTender ds.TenderCompany) (ds.Tender, []ds.Company, error) {
	var deletedCompanyTender ds.TenderCompany
	result := r.db.Where("\"CompanyID\" = ? and \"TenderID\" = ?", deleteFromTender.CompanyID,
		deleteFromTender.TenderID).Find(&deletedCompanyTender)
	if result.Error != nil {
		return ds.Tender{}, nil, result.Error
	}

	if result.RowsAffected == 0 {
		return ds.Tender{}, nil, fmt.Errorf("record not found")
	}
	if err := r.db.Delete(&deletedCompanyTender).Error; err != nil {
		return ds.Tender{}, nil, err
	}

	return r.GetTenderWithDataByID(deleteFromTender.TenderID)
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
