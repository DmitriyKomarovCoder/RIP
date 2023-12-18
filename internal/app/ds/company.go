package ds

// УСЛУГА
type Company struct {
	ID          uint   `json:"company_id" gorm:"primary_key"`
	CompanyName string `json:"name" gorm:"type:varchar(30);not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      string `json:"status" gorm:"type:varchar(20);not null"`
	IIN         string `json:"iin" gorm:"type:varchar(50)"`
	ImageURL    string `json:"image_url" gorm:"type:varchar(500)"`
}

type CompanyList struct {
	DraftID   uint       `json:"draft_id"`
	Companies *[]Company `json:"companies_list"`
}
type AddToCompanyID struct {
	CompanyID uint `json:"company"`
	UserID    int  `json:"user_id"`
}
