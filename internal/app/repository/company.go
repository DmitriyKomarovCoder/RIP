package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"strings"
)

func (r *Repository) GetOpenCompanys() (*[]ds.Company, error) {
	var tenders []ds.Company
	if err := r.db.Where("status = ?", "действует").Find(&tenders).Error; err != nil {
		return nil, err
	}
	return &tenders, nil
}

func (r *Repository) GetCompanyById(id uint) (*ds.Company, error) {
	var company ds.Company
	if err := r.db.Where("status = ?", "действует").First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *Repository) CompaniesList(name string) (*[]ds.Company, error) {
	name = strings.ToLower(name)

	var companies []ds.Company
	if err := r.db.Where("company_name LIKE ? AND status != ?", "%"+name+"%", "удален").Find(&companies).Error; err != nil {
		return nil, err
	}
	return &companies, nil
}

func (r *Repository) AddCompany(company *ds.Company) (uint, error) {
	company.Status = "действует"
	result := r.db.Create(&company)
	return company.ID, result.Error
}

func (r *Repository) DeleteCompany(id uint) error {
	company := ds.Company{}

	if err := r.db.First(&company, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&company).Update("status", "удален").Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateCompany(updatedCompany *ds.Company) (*ds.Company, error) {
	var oldCompany ds.Company

	if result := r.db.First(&oldCompany, updatedCompany.ID); result.Error != nil {
		return updatedCompany, result.Error
	}

	if updatedCompany.CompanyName != "" {
		oldCompany.CompanyName = updatedCompany.CompanyName
	}

	if updatedCompany.Status != "" {
		oldCompany.Status = updatedCompany.Status
	}

	if updatedCompany.ImageURL != "" {
		oldCompany.ImageURL = updatedCompany.ImageURL
	}

	if updatedCompany.Description != "" {
		oldCompany.Description = updatedCompany.Description
	}

	if updatedCompany.IIN != "" {
		oldCompany.IIN = updatedCompany.IIN
	}

	*updatedCompany = oldCompany
	result := r.db.Save(updatedCompany)
	return updatedCompany, result.Error
}

func (r *Repository) DeleteCompanyImage(companyId uint) string {
	company := ds.Company{}

	r.db.First(&company, "threat_id = ?", companyId)
	return company.ImageURL
}

func (r *Repository) AddCompanyToDraft(dataID uint, creatorID uint) (uint, error) {
	// получаем услугу
	data, err := r.GetCompanyById(dataID)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("нет такой услуги")
	}
	if data.Status == "удален" {
		return 0, errors.New("услуга удалена")
	}

	// получаем черновик
	var draftReq ds.Tender
	res := r.db.Where("user_id = ?", creatorID).Where("status != ?", "удалён").Take(&draftReq)

	// создаем черновик, если его нет
	if res.RowsAffected == 0 {
		newDraftRequestID, err := r.CreateTenderDraft(creatorID)
		if err != nil {
			return 0, err
		}

		draftReq.ID = newDraftRequestID
	}

	// добавляем запись в мм
	requestToData := ds.TenderCompany{
		CompanyID: dataID,
		TenderID:  draftReq.ID,
	}

	err = r.db.Create(&requestToData).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return 0, errors.New("услуга уже существует в заявке")
		}

		return 0, err
	}

	return draftReq.ID, nil
}
