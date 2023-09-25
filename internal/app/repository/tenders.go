package repository

import "RIP/internal/app/ds"

func (r *Repository) GetOpenTenders() (*[]ds.Tenders, error) {
	var tenders []ds.Tenders
	if err := r.db.Where("status = ?", "действует").Find(&tenders).Error; err != nil {
		return nil, err
	}
	return &tenders, nil
}

func (r *Repository) GetTenderById(id int) (*ds.Tenders, error) {
	var tender ds.Tenders
	if err := r.db.First(&tender, id).Error; err != nil {
		return nil, err
	}
	return &tender, nil
}

func (r *Repository) DeleteTender(id string) {
	query := "UPDATE tenders SET status = 'удален' WHERE id = $1;"
	r.db.Exec(query, id)
}
