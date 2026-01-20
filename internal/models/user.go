package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName       string `json:"full_name" gorm:"type:varchar(255);not null"`
	Email          string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Phone          string `json:"phone" gorm:"type:varchar(20);uniqueIndex"`
	DefaultAddress string `json:"default_address" gorm:"type:varchar(255);not null"`
}
