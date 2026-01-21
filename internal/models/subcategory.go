package models

import (
	"gorm.io/gorm"
)

type Subcategory struct {
	gorm.Model
	Name       string `gorm:"not null" json:"name"`
	CategoryID uint   `gorm:"constraint:OnDelete:CASCADE" json:"category_id"`
}
