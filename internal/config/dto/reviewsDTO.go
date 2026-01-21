package dto

type ReviewsCreate struct{
	UserID uint `json:"user_id" gorm:"not null,index" binding:"required"`
	MedicineID uint `json:"medicine_id" gorm:"not null,index" binding:"required"`
	Rating uint `json:"rating" gorm:"not null" binding:"required,min=1,max=10"`
	Text string `json:"text" gorm:"size:500"` 
}