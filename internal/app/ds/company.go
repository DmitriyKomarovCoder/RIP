package ds

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	CompanyName string `json:"tender_name" gorm:"type:varchar(30);not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      string `json:"status" gorm:"type:varchar(20);not null"`
	IIN         string `json:"iin" gorm:"type:varchar(50)"`
	ImageURL    string `json:"image_url" gorm:"type:varchar(500)"`
}
