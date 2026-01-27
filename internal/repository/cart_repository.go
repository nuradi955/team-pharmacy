package repository

import (
	"errors"
	"log/slog"
	"team-pharmacy/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreate(userID uint) (*models.Cart, error)
	GetCartWithItems(userID uint) (*models.Cart, error)

	GetItem(cartID uint, medicineID uint) (*models.CartItem, error)
	CreateItem(item *models.CartItem) error
	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID uint) error

	ClearCart(userID uint) error
}

type gormCartRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewCartRepository(db *gorm.DB, logger *slog.Logger) CartRepository {
	return &gormCartRepository{db: db, logger: logger.With("layer", "repository", "entity", "cart")}
}

func (r *gormCartRepository) GetOrCreate(userID uint) (*models.Cart, error) {
	const op = "repo.cart.get_or_create"
	r.logger.Debug(op,
		"user_id", userID,
	)
	var cart models.Cart

	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = models.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				r.logger.Error(op,
					"user_id", userID,
					"error", err,
				)
				return nil, err
			}
			return &cart, nil
		}
		r.logger.Error(op,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	return &cart, nil
}

func (r *gormCartRepository) GetCartWithItems(userID uint) (*models.Cart, error) {
	const op = "repo.cart.get_with_items"
	r.logger.Debug(op,
		"user_id", userID,
	)
	var cart models.Cart

	if err := r.db.Preload("Items.Medicine").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		r.logger.Error(op,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	return &cart, nil
}

func (r *gormCartRepository) CreateItem(item *models.CartItem) error {
	const op = "repo.cart_item.update"

	r.logger.Debug(op,
		"item_id", item.ID,
	)
	if err := r.db.Save(item).Error; err != nil {
		r.logger.Error(op,
			"item_id", item.ID,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *gormCartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *gormCartRepository) DeleteItem(itemID uint) error {

	const op = "repo.cart_item.delete"

	r.logger.Debug(op,
		"item_id", itemID,
	)

	if err := r.db.Delete(&models.CartItem{}, itemID).Error; err != nil {
		r.logger.Error(op,
			"item_id", itemID,
			"error", err,
		)
		return err
	}

	return nil
}

func (r *gormCartRepository) ClearCart(userID uint) error {
	const op = "repo.cart.clear"

	r.logger.Debug(op,
		"user_id", userID,
	)

	if err := r.db.Where("cart_id IN (?)",
		r.db.Model(&models.Cart{}).Select("id").Where("user_id = ?", userID),
	).
		Delete(&models.CartItem{}).
		Error; err != nil {
		r.logger.Error(op,
			"user_id", userID,
			"error", err,
		)
		return err
	}
	return nil
}

func (r *gormCartRepository) GetItem(cartID uint, medicineID uint) (*models.CartItem, error) {
	const op = "repo.cart_item.get"

	r.logger.Debug(op,
		"cart_id", cartID,
		"medicine_id", medicineID,
	)
	var item models.CartItem

	err := r.db.
		Where("cart_id = ? AND medicine_id = ?", cartID, medicineID).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		r.logger.Error(op,
			"cart_id", cartID,
			"medicine_id", medicineID,
			"error", err,
		)
		return nil, err
	}
	return &item, nil
}
