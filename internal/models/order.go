package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          uint   `json:"user_id" gorm:"not null"`
	Status          string `json:"status" gorm:"not null"`
	TotalPrice      float64
	DiscountTotal   float64
	FinalPrice      float64
	DeliveryAddress string
	Comment         string
}
