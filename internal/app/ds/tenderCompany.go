package ds

type TenderCompany struct {
	CompanyID uint    `json:"company_id" gorm:"primaryKey;column:CompanyID"`
	TenderID  uint    `json:"tender_id" gorm:"primaryKey;column:TenderID"`
	Tenders   Tender  `json:"tender" gorm:"foreignKey:TenderID"`
	Company   Company `json:"application" gorm:"foreignKey:CompanyID"`
	Cash      float64 `json:"cash" gorm:"numeric(10, 2)"`
}
