package dto

type MedicineCreate struct {
	Name                 string  `json:"name" binding:"required"`
	Description          string  `json:"description"`
	Price                float64 `json:"price" binding:"required"`
	StockQuantity        uint    `json:"stock_quantity" binding:"required"`
	CategoryID           *uint   `json:"category_id" binding:"required"`
	SubcategoryID        *uint   `json:"subcategory_id" binding:"required"`
	Manufacturer         string  `json:"manufacturer" binding:"required"`
	PrescriptionRequired bool    `json:"prescription_required" binding:"required"`
}
type MedicineUpdate struct {
	Name                 *string  `json:"name" binding:"required"`
	Description          *string  `json:"description"`
	Price                *float64 `json:"price" binding:"required"`
	StockQuantity        *uint    `json:"stock_quantity"`
	CategoryID           *uint    `json:"category_id" binding:"required"`
	SubcategoryID        *uint    `json:"subcategory_id" binding:"required"`
	Manufacturer         *string  `json:"manufacturer" binding:"required"`
	PrescriptionRequired *bool    `json:"prescription_required" binding:"required"`
}
