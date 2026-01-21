package repository

import (
	"team-pharmacy/internal/config/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetAllByUser(userID uint) ([]models.Review, error)
	GetAllByMedicine(medicineID uint) ([]models.Review, error)
	GetByID(id uint) (*models.Review, error)
	Delete(id uint) error
}
type ReviewRepo struct {
	db *gorm.DB
}


func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &ReviewRepo{db: db}
}
func (r *ReviewRepo) Create(review *models.Review) error {
	return r.db.Create(review).Error
}
func (r *ReviewRepo) GetAllByUser(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&models.Review{}).Where("user_id = ?", userID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
func (r *ReviewRepo) GetAllByMedicine(medicineID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&models.Review{}).Where("medicine_id = ?", medicineID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
func (r *ReviewRepo) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
func (r *ReviewRepo) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}
