package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name          string        `gorm:"not null" json:"name"`
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID" json:"subcategories, omitempty"`
}
