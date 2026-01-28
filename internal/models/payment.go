package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	OrderID uint `json:"order_id" gorm:"not null,index"`
	Amount uint `json:"amount" gorm:"not null"`
	Status string `json:"status" gorm:"not null,size:31,index"`
	Method string `json:"method" gorm:"not null,size:31,index"`
	PaidAT time.Time `json:"paid_at"`

	Order Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

}