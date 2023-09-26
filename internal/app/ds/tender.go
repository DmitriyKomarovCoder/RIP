package ds

import (
	"time"

	"gorm.io/gorm"
)

type Tender struct {
	gorm.Model
	ID             uint      `json:"id" grom:"primary_key"`
	Name           string    `json:"application_name" gorm:"type:varchar(255)"`
	Status         string    `json:"status" gorm:"type:varchar(15)"`
	CreationDate   time.Time `json:"creation_date" gorm:"type:date; not null; default:current_date"`
	FormationDate  time.Time `json:"formation_date" gorm:"type:date"`
	CompletionDate time.Time `json:"completion_date" gorm:"type:date"`
	ModeratorID    uint      `json:"moderator_id"`
	Moderator      User      `json:"moderator" gorm:"foreignkey:ModeratorID"`
	UserID         int       `json:"user_id"`
	User           User      `json:"user" gorm:"foreignkey:UserID"`
}
