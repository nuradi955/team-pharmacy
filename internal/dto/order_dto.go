package dto

import (
	"team-pharmacy/internal/models"
	"time"
)

type OrderCreateRequest struct {
	DeliveryAddress string `json:"delivery_address" binding:"required"`
	Comment         string `json:"comment"`
	Promocode       string `json:"promocode"`
}

type OrderShortResponse struct {
	ID         uint               `json:"id"`
	Status     models.OrderStatus `json:"status"`
	FinalPrice int64              `json:"final_price"`
	CreatedAt  time.Time          `json:"created_at"`
}

type OrderStatusRequest struct {
	Status *models.OrderStatus `json:"status"`
}

type OrderResponse struct {
	UserID          uint               `gorm:"index;not null"`
	Status          models.OrderStatus `gorm:"type:varchar(32);not null;index"`
	TotalPrice      int64              `gorm:"not null"`
	DiscountTotal   int64              `gorm:"not null"`
	FinalPrice      int64              `gorm:"not null"`
	DeliveryAddress string             `gorm:"not null"`
	Comment         *string            `gorm:"type:varchar(255)"`
	Items           []models.OrderItem `gorm:"constraint:OnDelete:CASCADE;"`
	// Payments        []models.Payment   `gorm:"constraint:OnDelete:CASCADE;"`
}
