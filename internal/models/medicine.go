package models

import "gorm.io/gorm"

type Medicine struct {
	gorm.Model
	Name                 string  `json:"name" gorm:"not null,size:100"`
	Description          string  `json:"description" gorm:"size:500"`
	Price                float64 `json:"price" gorm:"not null"`
	InStock              bool    `json:"in_stock" gorm:"not null"`
	StockQuantity        uint    `json:"stock_quantity" gorm:"not null"`
	CategoryID           *uint   `json:"category_id"`
	SubcategoryID        *uint   `json:"subcategory_id"`
	Manufacturer         string  `json:"manufacturer" gorm:"size:150,not null"`
	PrescriptionRequired bool    `json:"prescription_required"`
	AvgRating            float64 `json:"avg_rating" gorm:"index,not null check:rating>=1 AND rating<=10"`

	
	Category    Category    `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SETNULL;"`
	Subcategory Subcategory `json:"subcategory" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SETNULL;"`
}

func (m *Medicine) BeforeSave(tx *gorm.DB) error {
	m.InStock = m.StockQuantity > 0
	return nil
}
