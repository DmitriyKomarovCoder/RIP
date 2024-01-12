package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) TenderByUserID(userID string) (*[]ds.TenderResponse, error) {
	var tenders []ds.Tender
	var tenderResponses = []ds.TenderResponse{}
	result := r.db.Preload("User").
		//Preload("TenderCompanies.Tender.User").
		Preload("Moderator").
		Where("user_id = ? AND status != 'удален' AND status != 'черновик'", userID).
		Find(&tenders)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, tender := range tenders {
		tenderResponse := ds.TenderResponse{
			ID:        tender.ID,
			Name:      tender.Name,
			UserName:  tender.User.Name,
			UserLogin: tender.User.Login,
			//UserRole:       tender.User.Role,
			ModeratorName: tender.Moderator.Name,
			Status:        tender.Status,
			StatusCheck:   tender.StatusCheck,
			//ModeratorRole:  tender.Moderator.Role,
			ModeratorLogin:  tender.Moderator.Login,
			CreationDate:    tender.CreationDate,
			FormationDate:   tender.FormationDate,
			CompletionDate:  tender.CompletionDate,
			TenderCompanies: tender.TenderCompanies,
		}
		tenderResponses = append(tenderResponses, tenderResponse)
	}

	return &tenderResponses, result.Error
}

func (r *Repository) TenderByID(id uint) (*ds.TenderResponse, error) {
	tender := ds.Tender{}

	result := r.db.Preload("User").
		//Preload("TenderCompanies.Tenders").
		Preload("TenderCompanies.Company").
		First(&tender, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	tenderResponse := ds.TenderResponse{
		ID:        tender.ID,
		Name:      tender.Name,
		UserName:  tender.User.Name,
		UserLogin: tender.User.Login,
		//UserRole:       tender.User.Role,
		ModeratorName: tender.Moderator.Name,
		Status:        tender.Status,
		StatusCheck:   tender.StatusCheck,
		//ModeratorRole:  tender.Moderator.Role,
		ModeratorLogin:  tender.Moderator.Login,
		CreationDate:    tender.CreationDate,
		FormationDate:   tender.FormationDate,
		CompletionDate:  tender.CompletionDate,
		TenderCompanies: tender.TenderCompanies,
	}
	return &tenderResponse, result.Error
}

func (r *Repository) TenderModel(id uint) (*ds.Tender, error) {
	tender := ds.Tender{}

	result := r.db.Preload("User").
		//Preload("TenderCompanies.Tenders").
		Preload("TenderCompanies.Company").
		First(&tender, id)
	return &tender, result.Error
}

func (r *Repository) TenderDraftId(userId uint) (uint, error) {
	var tender ds.Tender
	result := r.db.
		Where("status = ? AND user_id = ?", "черновик", userId).
		First(&tender)
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return tender.ID, result.Error
}

func (r *Repository) GetTenderDraftID(creatorID uint) (uint, error) {
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
	//userInfo, err := GetUserInfo(r, creatorID)
	//if err != nil {
	//	return 0, err
	//}
	request := ds.Tender{
		UserID:       creatorID,
		Status:       "черновик",
		CreationDate: r.db.NowFunc(),
		ModeratorID:  nil,
		//CreatorLogin: userInfo.Login,
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (r *Repository) GetTenderWithDataByID(requestID uint, userId uint, isAdmin bool) (ds.Tender, []ds.Company, error) {
	var TenderRequest ds.Tender
	var companies []ds.Company

	//ищем такую заявку
	result := r.db.First(&TenderRequest, "id =?", requestID)
	if result.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return ds.Tender{}, nil, result.Error
	}
	if !isAdmin && TenderRequest.UserID == uint(userId) || isAdmin {
		res := r.db.
			Table("tender_companies").
			Select("companies.*").
			Where("status != ?", "удалён").
			Joins("JOIN companies ON tender_companies.\"CompanyID\" = companies.id").
			Where("tender_companies.\"TenderID\" = ?", requestID).
			Find(&companies)
		if res.Error != nil {
			r.logger.Error("error while getting for tender request")
			return ds.Tender{}, nil, res.Error
		}
	} else {
		return ds.Tender{}, nil, errors.New("ошибка доступа к данной заявке")
	}

	return TenderRequest, companies, nil
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

func (r *Repository) TendersList(statusID string, startDate time.Time, endDate time.Time) (*[]ds.TenderResponse, error) {
	var tenders []ds.Tender
	tenderResponses := []ds.TenderResponse{}
	if statusID == "" {
		result := r.db.
			Preload("User").
			Preload("Moderator").
			Where("status != 'удален' AND status != 'черновик' AND creation_date BETWEEN ? AND ?", startDate, endDate).
			Find(&tenders)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil
			} else {
				return nil, result.Error
			}
		}

		for _, tender := range tenders {
			tenderResponse := ds.TenderResponse{
				ID:        tender.ID,
				Name:      tender.Name,
				UserName:  tender.User.Name,
				UserLogin: tender.User.Login,
				//UserRole:       tender.User.Role,
				ModeratorName: tender.Moderator.Name,
				Status:        tender.Status,
				StatusCheck:   tender.StatusCheck,
				//ModeratorRole:  tender.Moderator.Role,
				ModeratorLogin:  tender.Moderator.Login,
				CreationDate:    tender.CreationDate,
				FormationDate:   tender.FormationDate,
				CompletionDate:  tender.CompletionDate,
				TenderCompanies: tender.TenderCompanies,
			}
			tenderResponses = append(tenderResponses, tenderResponse)
		}

		return &tenderResponses, result.Error
	}

	result := r.db.
		Preload("User").
		Where("status = ? AND status != 'черновик' AND creation_date BETWEEN ? AND ?", statusID, startDate, endDate).
		Find(&tenders)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, tender := range tenders {
		tenderResponse := ds.TenderResponse{
			ID:        tender.ID,
			Name:      tender.Name,
			UserName:  tender.User.Name,
			UserLogin: tender.User.Login,
			Status:    tender.Status,
			//UserRole:       tender.User.Role,
			ModeratorName: tender.Moderator.Name,
			//ModeratorRole:  tender.Moderator.Role,
			ModeratorLogin:  tender.Moderator.Login,
			StatusCheck:     tender.StatusCheck,
			CreationDate:    tender.CreationDate,
			FormationDate:   tender.FormationDate,
			CompletionDate:  tender.CompletionDate,
			TenderCompanies: tender.TenderCompanies,
		}
		tenderResponses = append(tenderResponses, tenderResponse)
	}

	return &tenderResponses, result.Error
}

//func (r *Repository) TenderList(status, start, end string, userId string, isAdmin bool) (*[]ds.Tender, error) {
//	var tender []ds.Tender
//	ending := " AND user_id = " + userId
//	if isAdmin {
//		ending = ""
//	}
//
//	query := r.db.Where("status != ?"+ending, "удален")
//
//	if status != "" {
//		query = query.Where("status = ?", status)
//	}
//
//	if start != "" {
//		query = query.Where("creation_date >= ?", start)
//	}
//
//	if end != "" {
//		query = query.Where("creation_date <= ?", end)
//	}
//	query = query.Order("id ASC")
//	result := query.
//		Preload("User").
//		Preload("TenderCompanies.Company").
//		Preload("TenderCompanies.Tenders").
//		Preload("Moderator").
//		Find(&tender)
//	return &tender, result.Error
//
//}

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

func (r *Repository) FormTenderRequestByID(creatorID uint) (error, uint) {
	var req ds.Tender
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, 0
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), 0
	}
	req.StatusCheck = "В обработке"
	req.Status = "сформирован"
	req.FormationDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err, 0
	}

	return nil, req.ID
}

func (r *Repository) GetTenderByUser(creatorID uint) (error, uint) {
	var req ds.Tender
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, 0
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), 0
	}

	return nil, req.ID
}

//func (r *Repository) RejectTenderRequestByID(requestID, moderatorID uint) error {
//	return r.finishRejectHelper("отклонён", requestID, moderatorID)
//}
//
//func (r *Repository) FinishEncryptDecryptRequestByID(requestID, moderatorID uint) error {
//	return r.finishRejectHelper("завершён", requestID, moderatorID)
//}

func (r *Repository) FinishRejectHelper(status string, requestID, moderatorID uint) error {
	//userInfo, err := GetUserInfo(r, moderatorID)
	//if err != nil {
	//	return err
	//}

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

	req.ModeratorID = &moderatorID
	//req.ModeratorLogin = userInfo.Login
	req.Status = status

	req.CompletionDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

//func (r *Repository) DeleteTenderByID(requestID uint) error { // ?
//	var req ds.Tender
//	res := r.db.
//		Where("user_id = ?", requestID). // ??
//		Take(&req)
//
//	if res.Error != nil {
//		return res.Error
//	}
//	if res.RowsAffected == 0 {
//		return errors.New("нет такой заявки")
//	}
//
//	req.Status = "удалён"
//	delTime := time.Now()
//	req.CompletionDate = time.Now()
//	if err := r.db.Save(&req).Error; err != nil {
//		return err
//	}
//	if err := r.db.Model(&req).Update("deleted_at", delTime).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

func (r *Repository) DeleteTenderByID(requestID uint) error { // ?
	var req ds.Tender
	if result := r.db.First(&req, requestID); result.Error != nil {
		return result.Error
	}

	req.Status = "удален"
	req.CompletionDate = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	result := r.db.Save(&req)

	return result.Error
}

//func (r *Repository) DeleteCompanyFromRequest(userId uint, companyID uint) (ds.Tender, []ds.Company, error) {
//	var request ds.Tender
//	r.db.Where("user_id = ?", userId).First(&request)
//
//	if request.ID == 0 {
//		return ds.Tender{}, nil, errors.New("no such request")
//	}
//	var company ds.Company
//	result := r.db.Where("\"CompanyID\" = ? and \"TenderID\" = ?", companyID,
//		request.ID).Find(&company)
//	if result.Error != nil {
//		return ds.Tender{}, nil, result.Error
//	}
//
//	if result.RowsAffected == 0 {
//		return ds.Tender{}, nil, fmt.Errorf("record not found")
//	}
//	if err := r.db.Delete(&request).Error; err != nil {
//		return ds.Tender{}, nil, err
//	}
//
//	return r.GetTenderWithDataByID(request.ID, userId, false)
//}

func (r *Repository) DeleteCompanyFromRequest(id int) error {
	var dh ds.TenderCompany
	if result := r.db.First(&dh, id); result.Error != nil {
		return result.Error
	}
	return r.db.Delete(&dh).Error
}

func (r *Repository) UpdateTenderCompany(id uint, cash float64) error {
	var updateCompany ds.TenderCompany
	r.db.Where("id = ?", id).First(&updateCompany)

	if updateCompany.TenderID == 0 {
		return errors.New("нет такой заявки")
	}
	updateCompany.Cash = cash

	if err := r.db.Save(&updateCompany).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveRequest(monitoringRequest ds.RequestAsyncService) error {
	var request ds.Tender
	err := r.db.First(&request, "id = ?", monitoringRequest.RequestId)
	if err.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return err.Error
	}
	//request.CompletionDate = time.Now()
	request.StatusCheck = monitoringRequest.Status
	res := r.db.Save(&request)
	return res.Error
}
