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
	UserID          uint               `json:"user_id"`
	Status          models.OrderStatus `json:"status"`
	TotalPrice      int64              `json:"total_price"`
	DiscountTotal   int64              `json:"discount_total"`
	FinalPrice      int64              `json:"final_price"`
	DeliveryAddress string             `json:"delivery_address"`
	Comment         *string            `json:"comment"`
	Items           []models.OrderItem `json:"items"`
	// Payments        []models.Payment   `gorm:"constraint:OnDelete:CASCADE;"`
}
