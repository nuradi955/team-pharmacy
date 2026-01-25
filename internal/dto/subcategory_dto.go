package dto

type SubcategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}
