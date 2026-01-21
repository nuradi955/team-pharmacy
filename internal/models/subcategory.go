package models

import (
	"gorm.io/gorm"
)

type Subcategory struct {
	gorm.Model
	Name       string `gorm:"not null" json:"name"`
	CategoryID uint   `gorm:"not null" json:"category_id"`
}
