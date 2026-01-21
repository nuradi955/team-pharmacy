package dto

type MedicineCreate struct {
	Name                 string  `json:"name" binding:"required"`
	Description          string  `json:"description"`
	Price                float64 `json:"price" binding:"required,min=0.01,max=999999999"`
	StockQuantity        uint    `json:"stock_quantity" binding:"required"`
	CategoryID           *uint   `json:"category_id" binding:"required"`
	SubcategoryID        *uint   `json:"subcategory_id" binding:"required"`
	Manufacturer         string  `json:"manufacturer" binding:"required"`
	PrescriptionRequired bool    `json:"prescription_required" binding:"required"`
}

type MedicineUpdate struct {
	Name                 *string  `json:"name" binding:"omitempty"`
	Description          *string  `json:"description"`
	Price                *float64 `json:"price" binding:"omitempty,min=0.01,max=999999999"`
	StockQuantity        *uint    `json:"stock_quantity"`
	CategoryID           *uint    `json:"category_id" binding:"omitempty"`
	SubcategoryID        *uint    `json:"subcategory_id" binding:"omitempty"`
	Manufacturer         *string  `json:"manufacturer" binding:"omitempty"`
	PrescriptionRequired *bool    `json:"prescription_required" binding:"omitempty"`
}
