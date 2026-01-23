package dto

type CategoryCreate struct {
	Name string `json:"name" binding:"required"`
}
