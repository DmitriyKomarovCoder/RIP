package ds

import (
	"gorm.io/gorm"
)

type Tenders struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	TenderName  string `json:"tender_name" gorm:"type:varchar(30);not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      string `json:"status" gorm:"type:varchar(20);not null"`
	ImageURL    string `json:"image_url" gorm:"type:varchar(500)"`
}
