package ds

type TenderCompany struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	CompanyID uint    `json:"company_id"`
	TenderID  uint    `json:"tender_id"`
	Tenders   Tender  `json:"tender" gorm:"foreignKey:TenderID"`
	Company   Company `json:"company" gorm:"foreignKey:CompanyID"`
	Cash      float64 `json:"cash" gorm:"numeric(10, 2)"`
}

type TenderCompanyUpdate struct {
	ID   uint    `json:"id" gorm:"primary_key"`
	Cash float64 `json:"cash" gorm:"numeric(10, 2)"`
}
