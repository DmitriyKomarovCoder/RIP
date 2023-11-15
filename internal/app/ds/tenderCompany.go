package ds

import "gorm.io/gorm"

type TenderCompany struct {
	gorm.Model
	CompanyID uint    `gorm:"primaryKey;column:CompanyID"`
	TenderID  uint    `gorm:"primaryKey;column:TenderID"`
	Tenders   Tender  `json:"tender" gorm:"foreignKey:TenderID"`
	Company   Company `json:"application" gorm:"foreignKey:CompanyID"`
	Cash      float64 `json:"cash" gorm:"numeric(10, 2)"`
}
