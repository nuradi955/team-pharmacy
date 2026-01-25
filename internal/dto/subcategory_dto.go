package dto

type SubcategoryCreate struct {
	Name       string `json:"name" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

type SubcategoryUpdate struct {
	Name       *string `json:"name"`
	CategoryID *uint   `json:"category_id"`
}
