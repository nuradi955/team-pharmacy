package dto

type CreateUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required,min=11"`
	DefaultAddress string `json:"default_address" binding:"required,max=255"`
}


type UpdateUserRequest struct {
	FullName       *string `json:"full_name" binding:"omitempty"`
	Email          *string `json:"email" binding:"omitempty,email"`
	Phone          *string `json:"phone" binding:"omitempty,min=11"`
	DefaultAddress *string `json:"default_address" binding:"omitempty,max=255"`
}
