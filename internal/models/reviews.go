package models

import (
	"time"
)

type Review struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null,index"`
	MedicineID uint      `json:"medicine_id" gorm:"not null,index"`
	Rating     uint      `json:"rating" gorm:"not null check:rating>=1 AND rating<=10"`
	Text       string    `json:"text" gorm:"size:500"`
	CreatedAt  time.Time `json:"created_at"`

	// User     Models.User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Medicine Models.Medicine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
