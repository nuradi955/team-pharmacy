package services

import (
	"errors"
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
}

func NewCartService(cartRepo repository.CartRepository, userRepo repository.UserRepository, medicineRepo repository.MedicineRepository) CartService {
	return &cartService{carts: cartRepo, users: userRepo, medicine: medicineRepo}
}

func (s *cartService) GetOrCreate(userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return s.carts.GetOrCreate(userID)

}

func (s *cartService) CreateItem(userID uint, req *dto.AddCartItemRequest) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	cart, err := s.carts.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("failed to get or create cart")
	}

	medicine, err := s.medicine.GetByID(req.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrMedicineNotFound
		}

		return nil, err
	}

	if int(medicine.StockQuantity) < req.Quantity {
		return nil, errors.New("не достаточно лекарств на складе")
	}

	existsItem, err := s.carts.GetItem(cart.ID, medicine.ID)
	if err != nil {
		return nil, err
	}

	if existsItem != nil {
		existsItem.Quantity += req.Quantity
		if existsItem.Quantity > int(medicine.StockQuantity) {
			return nil, errors.New("stock limit exceeded")
		}
		existsItem.PricePerUnit = int64(medicine.Price)
		if err := s.carts.UpdateItem(existsItem); err != nil {
			return nil, err
		}
		return s.GetCartWithItems(userID)
	}

	newItem := models.CartItem{
		CartID:       cart.ID,
		MedicineID:   medicine.ID,
		Quantity:     req.Quantity,
		PricePerUnit: int64(medicine.Price),
	}

	if err := s.carts.CreateItem(&newItem); err != nil {
		return nil, err
	}

	return s.GetCartWithItems(userID)

}

func (s *cartService) UpdateItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartItemResponse, error) {
	if userID == 0 || itemID == 0 {
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	cartWithItems, err := s.carts.GetCartWithItems(userID)
	if err != nil {
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
		return nil, errs.ErrItemNotFound
	}

	medicine, err := s.medicine.GetByID(item.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrMedicineNotFound
		}
		return nil, err
	}
	if req.Quantity <= 0 {
		if err := s.carts.DeleteItem(item.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if req.Quantity > int(medicine.StockQuantity) {
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
		return nil, err
	}
	return &newItem, nil
}

func (s *cartService) GetCartWithItems(userID uint) (*dto.CartResponse, error) {
	if userID == 0 {
		return nil, errors.New("id не может быть равно нулю")
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	cart, err := s.carts.GetCartWithItems(userID)
	if err != nil {
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

	return &resp, nil
}

func (s *cartService) DeleteItem(userID, itemID uint) error {
	if itemID == 0 || userID == 0 {
		return errors.New("id не может быть равно нулю")
	}
	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrUserNotFound
		}
		return err
	}
	if _, err := s.carts.GetOrCreate(userID); err != nil {
		return err
	}

	cartWithItems, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		return err
	}

	if cartWithItems == nil {
		return errs.ErrItemNotFound
	}

	for i := range cartWithItems.Items {
		if cartWithItems.Items[i].ID == itemID {
			if err := s.carts.DeleteItem(itemID); err != nil {
				return err
			}
			return nil
		}

	}
	return errs.ErrItemNotFound
}

func (s *cartService) ClearCart(userID uint) error {
	if userID == 0 {
		return errors.New("id не может быть равно нулю")
	}

	if _, err := s.users.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrUserNotFound
		}
		return err
	}

	return s.carts.ClearCart(userID)

}
