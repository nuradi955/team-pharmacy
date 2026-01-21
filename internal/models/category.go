package models

import "time"

type Category struct {
	ID            uint          `gorm:"primarykey" json:"id"`
	Name          string        `gorm:"not null" json:"name"`
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID" json:"subcategories, omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"created_at"`
}
