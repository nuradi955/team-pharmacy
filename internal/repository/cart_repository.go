package repository

import (
	"errors"
	"team-pharmacy/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreate(userID uint) (*models.Cart, error)
	GetCartWithItems(userID uint) (*models.Cart, error)

	CreateItem(item *models.CartItem) error
	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID uint) error

	ClearCart(userID uint) error
}

type gormCartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &gormCartRepository{db: db}
}

func (r *gormCartRepository) GetOrCreate(userID uint) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = models.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return &cart, err
	}
	return &cart, nil
}

func (r *gormCartRepository) GetCartWithItems(userID uint) (*models.Cart, error) {

	var cart models.Cart

	if err := r.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *gormCartRepository) CreateItem(item *models.CartItem) error {

	return r.db.Create(item).Error

}

func (r *gormCartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *gormCartRepository) DeleteItem(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

func (r *gormCartRepository) ClearCart(userId uint) error {
	var cart models.Cart

	if err := r.db.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return r.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
}
