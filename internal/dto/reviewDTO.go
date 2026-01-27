package dto

type ReviewCreate struct {
	UserID     uint   `json:"user_id" binding:"required"`
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Rating     uint   `json:"rating" binding:"required,min=1,max=10"`
	Text       string `json:"text"`
}

type ReviewUpdate struct {
	// UserID     *uint   `json:"user_id"  binding:"omitempty"`
	// MedicineID *uint   `json:"medicine_id" binding:"omitempty"`
	Rating     *uint   `json:"rating" binding:"omitempty,min=1,max=10"`
	Text       *string `json:"text"`
}
