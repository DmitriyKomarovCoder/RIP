package repository

import "RIP/internal/app/ds"

func (r *Repository) GetOpenCompanys() (*[]ds.Company, error) {
	var tenders []ds.Company
	if err := r.db.Where("status = ?", "действует").Find(&tenders).Error; err != nil {
		return nil, err
	}
	return &tenders, nil
}

func (r *Repository) GetCompanyById(id int) (*ds.Company, error) {
	var tender ds.Company
	if err := r.db.First(&tender, id).Error; err != nil {
		return nil, err
	}
	return &tender, nil
}

func (r *Repository) DeleteCompany(id string) {
	query := "UPDATE companies SET status = 'удален' WHERE id = $1;"
	r.db.Exec(query, id)
}
