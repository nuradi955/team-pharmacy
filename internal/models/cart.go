package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     uint       `gorm:"uniqueIndex"`
	Items      []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
	TotalPrice float64    `json:"total_price"`
}

type CartItem struct {
	gorm.Model
	CartID       uint      `gorm:"index;not null"`
	Cart         *Cart     `gorm:"constraint:OnDelete:CASCADE;"`
	MedicineID   uint      `gorm:"index;not null"`
	Medicine     *Medicine `gorm:"constraint:OnDelete:CASCADE;"`
	Quantity     int       `gorm:"not null"`
	PricePerUnit int64     `gorm:"not null"`
	LineTotal    int64     `gorm:"not null"`
}
