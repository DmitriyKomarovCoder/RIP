package ds

import "gorm.io/gorm"

type TenderApplication struct {
	gorm.Model
	ID            uint        `json:"id" gorm:"primary_key"`
	TenderID      uint        `json:"tende_id"`
	Tenders       Tenders     `json:"tender" gorm:"foreignkey:TenderID"`
	ApplicationID uint        `json:"application_id"`
	Application   Application `json:"application" gorm:"foreignkey:ApplicationID"`
}
