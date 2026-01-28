package services

import (
	"errors"
	"log/slog"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/errs"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"

	"gorm.io/gorm"
)

type CartService interface {
	GetOrCreate(userID uint) (*models.Cart, error)

	CreateItem(userID uint, item *dto.AddCartItemRequest) (*dto.CartResponse, error)
	UpdateItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartItemResponse, error)
	GetCartWithItems(userID uint) (*dto.CartResponse, error)
	DeleteItem(userID, itemID uint) error

	ClearCart(userID uint) error
}

type cartService struct {
	carts    repository.CartRepository
	users    repository.UserRepository
	medicine repository.MedicineRepository
	logger   *slog.Logger
}

func NewCartService(cartRepo repository.CartRepository,
	userRepo repository.UserRepository,
	medicineRepo repository.MedicineRepository,
	logger *slog.Logger,
) CartService {
	return &cartService{
		carts:    cartRepo,
		users:    userRepo,
		medicine: medicineRepo,
		logger: logger.With("layer", "service",
			"entity", "cart",
		)}
}

func (s *cartService) GetOrCreate(userID uint) (*models.Cart, error) {
	s.logger.Info("get or create cart started",
		"user_id", userID,
	)
	if userID == 0 {
		s.logger.Warn("invalid user_id",
			"user_id", userID,
		)
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return nil, errs.ErrUserNotFound
		}

		s.logger.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	cart, err := s.carts.GetOrCreate(userID)
	if err != nil {
		s.logger.Error("failed to get or create cart",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	s.logger.Info("cart ready",
		"user_id", userID,
		"cart_id", cart.ID,
	)

	return cart, nil

}

func (s *cartService) CreateItem(userID uint, req *dto.AddCartItemRequest) (*dto.CartResponse, error) {
	s.logger.Info("add item to cart started",
		"user_id", userID,
		"medicine_id", req.MedicineID,
		"quantity", req.Quantity,
	)

	if userID == 0 {
		s.logger.Warn("invalid user_id",
			"user_id", userID,
		)
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return nil, errs.ErrUserNotFound
		}
		s.logger.Info("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	cart, err := s.carts.GetOrCreate(userID)
	if err != nil {
		s.logger.Error("failed to get or create cart",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("failed to get or create cart")
	}

	medicine, err := s.medicine.GetByID(req.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Info("medicine not found",
				"medicine_id", req.MedicineID)
			return nil, errs.ErrMedicineNotFound
		}
		s.logger.Error("failed to get medicine",
			"medicine_id", req.MedicineID,
			"error", err,
		)
		return nil, err
	}

	if int(medicine.StockQuantity) < req.Quantity {
		s.logger.Warn("not enough stock",
			"medicine_id", req.MedicineID,
			"stock", medicine.StockQuantity,
			"requested", req.Quantity,
		)

		return nil, errors.New("не достаточно лекарств на складе")
	}

	existsItem, err := s.carts.GetItem(cart.ID, medicine.ID)
	if err != nil {
		s.logger.Error("failed to check existing item",
			"cart_id", cart.ID,
			"medicine_id", medicine.ID,
			"error", err,
		)
		return nil, err
	}

	if existsItem != nil {
		existsItem.Quantity += req.Quantity
		if existsItem.Quantity > int(medicine.StockQuantity) {
			s.logger.Warn("stock limit exceeded",
				"medicine_id", medicine.ID,
			)
			return nil, errors.New("stock limit exceeded")
		}
		existsItem.PricePerUnit = int64(medicine.Price)
		if err := s.carts.UpdateItem(existsItem); err != nil {
			s.logger.Error("failed to update cart item ",
				"item_id", existsItem.ID,
				"error", err,
			)
			return nil, err
		}
		s.logger.Info("cart item quantity increased",
			"cart_id", cart.ID,
			"medicine_id", medicine.ID,
		)

		return s.GetCartWithItems(userID)
	}

	newItem := models.CartItem{
		CartID:       cart.ID,
		MedicineID:   medicine.ID,
		Quantity:     req.Quantity,
		PricePerUnit: int64(medicine.Price),
	}

	if err := s.carts.CreateItem(&newItem); err != nil {
		s.logger.Error("failed to create cart item ",
			"cart_id", cart.ID,
			"medicine_id", medicine.ID,
			"error", err,
		)

		return nil, err
	}

	s.logger.Info("cart item created",
		"cart_id", cart.ID,
		"medicine_id", medicine.ID,
	)

	return s.GetCartWithItems(userID)

}

func (s *cartService) UpdateItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartItemResponse, error) {
	s.logger.Info("update cart item start",
		"user_id", userID,
		"item_id", itemID,
		"quantity", req.Quantity,
	)

	if userID == 0 || itemID == 0 {
		s.logger.Warn("invalid user_id or item_id",
			"user_id", userID,
			"item_id", itemID,
		)

		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return nil, errs.ErrUserNotFound
		}
		s.logger.Warn("user not found",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	cartWithItems, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		s.logger.Error("failed to get cart with items",
			"user_id", userID,
			"error", err,
		)

		return nil, err
	}

	var item *models.CartItem

	for i := range cartWithItems.Items {
		if cartWithItems.Items[i].ID == itemID {
			item = &cartWithItems.Items[i]
			break
		}
	}

	if item == nil {
		s.logger.Error("the cart is empty",
			"user_id", userID,
			"cart_id", cartWithItems,
		)

		return nil, errs.ErrItemNotFound
	}

	medicine, err := s.medicine.GetByID(item.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("medicine not found",
				"user_id", userID,
				"medicine_id", item.MedicineID,
			)
			return nil, errs.ErrMedicineNotFound
		}

		s.logger.Error("failed to get medicine",
			"user_id", userID,
			"medicine_id", item.MedicineID,
			"error", err,
		)
		return nil, err
	}

	if req.Quantity <= 0 {
		s.logger.Info("deleting cart item (quantity <= 0)",
			"user_id", userID,
			"item_id", item.ID,
		)
		if err := s.carts.DeleteItem(item.ID); err != nil {
			s.logger.Error("failed to delete cart item",
				"user_id", userID,
				"item_id", item.ID,
				"error", err,
			)
			return nil, err
		}
		return nil, nil
	}

	if req.Quantity > int(medicine.StockQuantity) {
		s.logger.Warn("stock limit exceeded",
			"user_id", userID,
			"item_id", item.ID,
			"requested", req.Quantity,
			"available", medicine.StockQuantity,
		)
		return nil, errors.New("stock limit exceeded")
	}

	lineTotal := int64(req.Quantity) * int64(medicine.Price)
	newItem := dto.CartItemResponse{
		ItemID:       item.ID,
		MedicineID:   item.MedicineID,
		Quantity:     req.Quantity,
		PricePerUnit: int64(medicine.Price),
		LineTotal:    lineTotal,
	}

	item.Quantity = req.Quantity
	item.PricePerUnit = int64(medicine.Price)
	if err := s.carts.UpdateItem(item); err != nil {

		s.logger.Error("failed to update cart item",
			"user_id", userID,
			"item_id", item.ID,
			"error", err,
		)
		return nil, err
	}
	s.logger.Info("cart item updated successfully",
		"user_id", userID,
		"item_id", item.ID,
		"quantity", req.Quantity,
		"line_total", lineTotal,
	)
	return newItem, nil
}

func (s *cartService) GetCartWithItems(userID uint) (*dto.CartResponse, error) {
	s.logger.Info("get cart with items started",
		"user_id", userID,
	)
	if userID == 0 {
		s.logger.Warn("invalid user_id",
			"user_id", userID,
		)
		return nil, errors.New("id не может быть равно нулю")
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return nil, errs.ErrUserNotFound
		}
		s.logger.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	cart, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		s.logger.Error("failed to get cart with items",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	var total int64
	items := make([]dto.CartItemResponse, 0, len(cart.Items))

	for _, it := range cart.Items {
		lineTotal := int64(it.Quantity) * it.PricePerUnit
		total += lineTotal

		items = append(items, dto.CartItemResponse{
			ItemID:       it.ID,
			MedicineID:   it.MedicineID,
			Quantity:     it.Quantity,
			PricePerUnit: it.PricePerUnit,
			LineTotal:    lineTotal,
		})
	}

	resp := dto.CartResponse{
		UserID:     cart.UserID,
		Items:      items,
		TotalPrice: total,
	}

	s.logger.Info("get cart with items finished",
		"user_id", userID,
		"items_count", len(items),
		"total_price", total,
	)

	return &resp, nil
}

func (s *cartService) DeleteItem(userID, itemID uint) error {
	s.logger.Info("delete cart item started",
		"user_id", userID,
		"item_id", itemID,
	)
	if itemID == 0 || userID == 0 {
		s.logger.Warn("invalid user_id or item_id",
			"user_id", userID,
			"item_id", itemID,
		)
		return errors.New("id не может быть равно нулю")
	}
	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return errs.ErrUserNotFound
		}
		s.logger.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return err
	}
	if _, err := s.carts.GetOrCreate(userID); err != nil {
		s.logger.Error("failed to get or create cart",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	cartWithItems, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		s.logger.Error("failed to get cart with items",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	if cartWithItems == nil || len(cartWithItems.Items) == 0 {
		s.logger.Warn("cart is empty",
			"user_id", userID,
		)
		return errs.ErrItemNotFound
	}

	for i := range cartWithItems.Items {
		if cartWithItems.Items[i].ID == itemID {
			if err := s.carts.DeleteItem(itemID); err != nil {
				s.logger.Error("failed to delete cart item",
					"user_id", userID,
					"item_id", itemID,
					"error", err,
				)
				return err
			}
			s.logger.Info("cart item deleted successfully",
				"user_id", userID,
				"item_id", itemID,
			)
			return nil
		}

	}
	s.logger.Warn("cart item not found",
		"user_id", userID,
		"item_id", itemID,
	)
	return errs.ErrItemNotFound
}

func (s *cartService) ClearCart(userID uint) error {
	s.logger.Info("clear cart started",
		"user_id", userID,
	)
	if userID == 0 {
		s.logger.Warn("invalid user_id",
			"user_id", userID,
		)
		return errors.New("id не может быть равно нулю")
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("user not found",
				"user_id", userID,
			)
			return errs.ErrUserNotFound
		}

		s.logger.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return err
	}
	if err := s.carts.ClearCart(userID); err != nil {
		s.logger.Error("failed to clear cart",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	s.logger.Info("cart cleared successfully",
		"user_id", userID,
	)

	return nil

}
