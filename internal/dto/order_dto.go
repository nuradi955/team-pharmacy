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
	Status     models.OrderStatus `json:"status"`
	FinalPrice int64              `json:"final_price"`
	CreatedAt  time.Time          `json:"created_at"`
}

type OrderStatusRequest struct {
	Status *models.OrderStatus `json:"status" binding:"required"`
}

type OrderResponse struct {
	UserID          uint                `json:"user_id"`
	Status          string              `json:"status"`
	TotalPrice      int64               `json:"total_price"`
	DiscountTotal   int64               `json:"discount_total"`
	FinalPrice      int64               `json:"final_price"`
	DeliveryAddress string              `json:"delivery_address"`
	Comment         string              `json:"comment"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
	// Payments        []models.Payment   `gorm:"constraint:OnDelete:CASCADE;"`
}

type OrderItemResponse struct {
	MedicineID   uint   `json:"medicine_id"`
	MedicineName string `json:"medicine_name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int64  `json:"price_per_unit"`
	LineTotal    int64  `json:"line_total"`
}
