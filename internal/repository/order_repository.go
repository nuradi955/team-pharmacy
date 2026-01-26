package repository

import (
	"errors"
	"team-pharmacy/internal/errs"
	"team-pharmacy/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetByID(orderID uint) (*models.Order, error)
	GetListOrders(userID uint) ([]models.Order, error)
	UpdateOrder(orderID uint, status *models.OrderStatus) error
	CreateOrderWithClearCart(order *models.Order, cartID uint) error
}

type gormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &gormOrderRepository{db: db}
}

func (r *gormOrderRepository) CreateOrder(order *models.Order) error {

	return r.db.Create(&order).Error

}

func (r *gormOrderRepository) GetByID(orderID uint) (*models.Order, error) {
	var order *models.Order

	if err := r.db.Preload("Items").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}
	return order, nil
}

func (r *gormOrderRepository) GetListOrders(userID uint) ([]models.Order, error) {
	var list []models.Order

	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormOrderRepository) UpdateOrder(orderID uint, status *models.OrderStatus) error {

	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error

}

func (r *gormOrderRepository) CreateOrderWithClearCart(order *models.Order, cartID uint) error {

	if err := r.db.Create(order).Error; err != nil {
		return err
	}

	if err := r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error; err != nil {
		return err
	}

	return nil
}
