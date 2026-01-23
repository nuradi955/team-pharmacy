package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID uint `gorm:"index;not null"`
	User   *User

	Status OrderStatus `gorm:"type:varchar(32);not null;index"`

	TotalPrice    int64 `gorm:"not null"`
	DiscountTotal int64 `gorm:"not null"`
	FinalPrice    int64 `gorm:"not null"`

	DeliveryAddress string `gorm:"not null"`
	Comment         string `gorm:"type:varchar(255)"`

	Items []OrderItem `gorm:"constraint:OnDelete:CASCADE;"`
	//Payments []Payment   `gorm:"constraint:OnDelete:CASCADE;"`
}

type OrderItem struct {
	gorm.Model
	OrderID      uint `gorm:"index;not null"`
	Order        *Order
	MedicineID   uint `gorm:"index;not null"`
	Medicine     *Medicine
	MedicineName string `gorm:"type:varchar(255)"`
	Quantity     int    `gorm:"not null"`
	PricePerUnit int64  `gorm:"not null"`
	LineTotal    int64  `gorm:"not null"`
}
