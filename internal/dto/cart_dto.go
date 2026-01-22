package dto

type AddCartItemRequest struct {
	MedicineID uint `json:"medicine_id" binding:"required,gt=0"`
	Quantity   int  `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartResponse struct {
	UserID     uint               `json:"user_id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice int64              `json:"total_price"`
}

type CartItemResponse struct {
	ItemID       uint  `json:"item_id"`
	MedicineID   uint  `json:"medicine_id"`
	Quantity     int   `json:"quantity"`
	PricePerUnit int64 `json:"price_per_unit"`
	LineTotal    int64 `json:"line_total"`
}
