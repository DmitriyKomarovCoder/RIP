package ds

import "gorm.io/gorm"

type Moderator struct {
	gorm.Model
	ID            uint   `json:"id" gorm:"primary_key"`
	ModeratorName string `json:"moderator_name" gorm:"type:varchar(50)"`
}
