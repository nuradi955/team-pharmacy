package dto

type ReviewsCreate struct {
	UserID     uint   `json:"user_id" bindng:"required"`
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Rating     uint   `json:"rating" binding:"required"`
	Text       string `json:"text"`
}
