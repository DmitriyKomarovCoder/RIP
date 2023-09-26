package ds

import "gorm.io/gorm"

type TenderCompany struct {
	gorm.Model
	ID        uint    `json:"id"`
	TenderID  uint    `json:"tende_id"`
	Tenders   Tender  `json:"tender" gorm:"foreignkey:TenderID"`
	CompanyID int     `json:"application_id"`
	Company   Company `json:"application" gorm:"foreignkey:CompanyID"`
}
