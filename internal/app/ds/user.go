package ds

type User struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	ModeratorName string `json:"moderator_name" gorm:"type:varchar(50)"`
	IsModerator   bool   `json:"is_moderator" gorm:"type:bool"`
}
