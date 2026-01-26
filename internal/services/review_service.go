package services

import (
	"errors"
	"strings"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

type ReviewService interface {
	Create(req dto.ReviewCreate) (*models.Review, error)
	GetAllByMedicine(medicine_id uint) ([]models.Review, error)
	GetAllByUser(user_id uint) ([]models.Review, error)
	GetByID(id uint) (*models.Review, error)
	Update(req dto.ReviewUpdate, id uint) error
	Delete(id uint) error
}

type reviewService struct {
	reviewRepo   repository.ReviewRepository
	medicineRepo repository.MedicineRepository
	userRepo     repository.UserRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository,medicineRepo repository.MedicineRepository,userRepo repository.UserRepository,) ReviewService {
	return &reviewService{reviewRepo:reviewRepo,medicineRepo:medicineRepo,userRepo:userRepo,}
}

func (r *reviewService) Create(req dto.ReviewCreate) (*models.Review, error) {
	if _, err := r.userRepo.GetByID(req.UserID); err != nil {
		return nil, errors.New("Invalid user")
	}

	if _, err := r.medicineRepo.GetByID(req.MedicineID); err != nil {
		return nil, errors.New("Invalid medicine")
	}
    if req.Rating <1 || req.Rating>10 {
		return nil,errors.New("value isnt correct")
	}

	text := strings.TrimSpace(req.Text)

	review := &models.Review{
		UserID:     req.UserID,
		MedicineID: req.MedicineID,
		Rating:     req.Rating,
		Text:       text,
	}

	if err := r.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	//  Нужно тут написать обновление среднего рейтинга 
avg,err := r.reviewRepo.GetAvgRatingByMedicine(review.MedicineID);
	if err != nil {
		return nil,err
	}
	if err:= r.medicineRepo.UpdateAvgRating(review.MedicineID,avg);err != nil{
		return nil,err
	}
	return review,nil 
}

func (r *reviewService) GetAllByUser(user_id uint) ([]models.Review, error) {
	return r.reviewRepo.GetAllByUser(user_id)
}
func (r *reviewService) GetAllByMedicine(medicine_id uint) ([]models.Review, error) {
	return r.reviewRepo.GetAllByMedicine(medicine_id)
}

func (r *reviewService) GetByID(id uint) (*models.Review, error) {
	if id == 0 {
		return nil, errors.New("id cant be zero")
	}
	return r.reviewRepo.GetByID(id)
}

func (r *reviewService) Update(req dto.ReviewUpdate, id uint) error {
	review, err := r.reviewRepo.GetByID(id)
	if err != nil {
		return err
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	

	if req.Text != nil {
		review.Text = strings.TrimSpace(*req.Text)
	}

	if err := r.reviewRepo.Update(review); 
	err != nil {
		return err
	}

	avg,err := r.reviewRepo.GetAvgRatingByMedicine(review.MedicineID);
	if err != nil {
		return err
	}
	if err:= r.medicineRepo.UpdateAvgRating(review.MedicineID,avg);err != nil{
		return err
	}

	return nil
}

func (r *reviewService) Delete(id uint) error {
	review, err := r.reviewRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.reviewRepo.Delete(id); err != nil {
		return err
	}

	avg,err := r.reviewRepo.GetAvgRatingByMedicine(review.MedicineID);
	if err != nil {
		return err
	}
	if err:= r.medicineRepo.UpdateAvgRating(review.MedicineID,avg);err != nil{
		return err
	}

	return nil
}