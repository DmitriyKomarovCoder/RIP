package ds

import (
	"RIP/internal/app/role"
	"time"
)

type User struct {
	UserId           int       `json:"id" gorm:"primaryKey"`
	Login            string    `json:"login" binding:"required,max=64"`
	Role             role.Role `json:"role" sql:"type:string"`
	Name             string    `json:"name"`
	Password         string    `json:"password" binding:"required,min=8,max=64"`
	RegistrationDate time.Time
}

//type UserDetail struct {
//	//UserId           int       `json:"id" gorm:"primaryKey"`
//	Login string `json:"login" binding:"required,max=64"`
//	//Role             role.Role `json:"role" sql:"type:string"`
//	Name string `json:"name"`
//	//Password         string    `json:"password" binding:"required,min=8,max=64"`
//	//RegistrationDate time.Time
//}

//type UserLogin struct {
//	Login    string `json:"login" binding:"required,max=64"`
//	Password string `json:"password" binding:"required,min=8,max=64"`
//}
//
//type UserSignUp struct {
//	Login    string `json:"login" binding:"required,max=64"`
//	Name     string `json:"name"`
//	Password string `json:"password" binding:"required,min=8,max=64"`
//}
//
//type RegisterResp struct {
//	Ok bool `json:"ok"`
//}
