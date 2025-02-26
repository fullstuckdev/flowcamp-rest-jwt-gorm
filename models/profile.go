package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserId uint `json:"user_id"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}