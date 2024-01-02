package ds

import (
	"time"
)

// ЗАЯВКА ТЕНДЕРА
type Tender struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	Name           string    `json:"application_name" gorm:"type:varchar(255)"`
	Status         string    `json:"status" gorm:"type:varchar(15)"`
	CreationDate   time.Time `json:"creation_date" gorm:"type:date; not null; default:current_date"`
	FormationDate  time.Time `json:"formation_date" gorm:"type:date"`
	CompletionDate time.Time `json:"completion_date" gorm:"type:date"`
	CreatorLogin   string    `json:"creator_login"`
	ModeratorLogin string    `json:"moderator_login"`
	ModeratorID    uint      `json:"moderator_id"`
	//Moderator      User      `json:"moderator" gorm:"foreignkey:ModeratorID"`
	UserID uint `json:"user_id"`
	//User   User `json:"user" gorm:"foreignkey:UserID"`
}

type NewStatus struct {
	Status string `json:"status"`
}

type TenderDetails struct {
	Tender  *Tender    `json:"tender"`
	Company *[]Company `json:"companies"`
}
