package ds

import (
	"RIP/internal/app/role"
	"github.com/golang-jwt/jwt"
)

//type JwtClaims struct {
//	jwt.StandardClaims
//	UserId  int  `json:"userId"`
//	IsAdmin bool `json:"isAdmin"`
//}
//
//type Role int
//
//const (
//	Client Role = iota // 0
//	Admin              // 1
//)

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
	Role   role.Role
}
