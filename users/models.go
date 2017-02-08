package users

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	AuthID  string `json:"AuthId" gorm:"not null;unique"`
	IsAdmin bool   `json:"IsAdmin"`
}
