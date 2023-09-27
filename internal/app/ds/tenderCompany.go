package ds

import "gorm.io/gorm"

type TenderCompany struct {
	gorm.Model
	CompanyID int     `gorm:"primaryKey;column:CompanyID"`
	TenderID  int     `gorm:"primaryKey;column:TenderID"`
	Tenders   Tender  `json:"tender" gorm:"foreignKey:TenderID"`
	Company   Company `json:"application" gorm:"foreignKey:CompanyID"`
}
