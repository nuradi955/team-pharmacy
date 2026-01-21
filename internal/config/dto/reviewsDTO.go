package dto

type ReviewsCreate struct {
	UserID     uint   `json:"user_id" binding:"required"`
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Rating     uint   `json:"rating" binding:"required,min=1,max=10"`
	Text       string `json:"text"`
}
