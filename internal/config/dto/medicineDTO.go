package dto

type MedicineCreate struct {
	Name                 string  `json:"name" gorm:"not null,size:100"`
	Description          string  `json:"description" gorm:"size:500"`
	Price                float64 `json:"price" gorm:"not null"`
	StockQuantity        uint    `json:"stock_quantity" gorm:"not null"`
	CategoryID           uint    `json:"category_id"`
	SubcategoryID        uint    `json:"subcategory_id"`
	Manufacturer         string  `json:"manufacturer" gorm:"size:150,not null"`
	PrescriptionRequired bool    `json:"prescription_required"`
}
type MedicineUpdate struct {
	Name                 *string  `json:"name" gorm:"not null,size:100"`
	Description          *string  `json:"description" gorm:"size:500"`
	Price                *float64 `json:"price" gorm:"not null"`
	StockQuantity        *uint    `json:"stock_quantity" gorm:"not null"`
	CategoryID           *uint    `json:"category_id"`
	SubcategoryID        *uint    `json:"subcategory_id"`
	Manufacturer         *string  `json:"manufacturer" gorm:"size:150,not null"`
	PrescriptionRequired *bool    `json:"prescription_required"`
}
