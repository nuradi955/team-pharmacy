package dto

type CreateUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required,min=11"`
	DefaultAddress string `json:"default_address" binding:"required,max=255"`
}

type UpdateUserRequest struct {
	FullName       *string `json:"full_name" omitempty:"required"`
	Email          *string `json:"email" omitempty:"required,email"`
	Phone          *string `json:"phone" omitempty:"required,min=11"`
	DefaultAddress *string `json:"default_address" omitempty:"required,max=255"`
}
