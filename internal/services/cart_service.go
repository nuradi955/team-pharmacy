package services

import (
	"errors"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
)

type CartService interface {
	GetOrCreate(userID uint) (*models.Cart, error)
	GetCartWithItems(userID uint) (*models.Cart, error)

	CreateItem(userID uint, item *dto.AddCartItemRequest) error
	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID uint) error

	ClearCart(userID uint) error
}

type cartService struct {
	carts repository.CartRepository
	users repository.UserRepository
}

func NewCartService(cartRepo repository.CartRepository, userRepo repository.UserRepository) CartService {
	return &cartService{carts: cartRepo, users: userRepo}
}

func (s *cartService) GetOrCreate(userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("invalid ID")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return s.carts.GetOrCreate(userID)

}

func (s *cartService) GetCartWithItems(userID uint) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("id must be greater than 0")
	}
	_, err := s.users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	cart, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, errors.New("cart not found")
	}
	cart.TotalPrice = 0
	for _, it := range cart.Items {
		cart.TotalPrice += float64(it.Quantity) * float64(it.PricePerUnit)
	}

	return cart, nil

}

func (s *cartService) CreateItem(userID uint, req *dto.AddCartItemRequest) (*models.Cart, error) {
	if userID == 0 {
		return nil, errors.New("id must be greater than 0")
	}
	if req == nil || req.MedicineID == 0 || req.Quantity == 0 {
		return nil, errors.New("invalid request")
	}

	if _, err := s.users.GetByID(userID); err != nil {
		return nil, err
	}
	
	cartWithItems, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		return nil, err
	}

	for i := range cartWithItems.Items {
		it := &cartWithItems.Items[i]
		if it.MedicineID == req.MedicineID {
			it.Quantity += req.Quantity
			if err := s.carts.UpdateItem(it); err != nil {
				return nil, err
			}
			updated, err := s.carts.GetCartWithItems(userID)
			if err != nil {
				return nil, err
			}
			updated.TotalPrice = 0
			for _, v := range updated.Items {
				updated.TotalPrice += float64(v.Quantity) * float64(v.PricePerUnit)
			}
			return updated, nil
		}
	}

	newitem := models.CartItem{
		CartID:     cartWithItems.ID,
		MedicineID: req.MedicineID,
		Quantity:   req.Quantity,
		//PricePerUnit:
	}
	if err := s.carts.CreateItem(&newitem); err != nil {
		return nil, err
	}

	updated, err := s.carts.GetCartWithItems(userID)
	if err != nil {
		return nil, err
	}
	updated.TotalPrice = 0
	for _, it := range updated.Items {
		updated.TotalPrice += float64(it.Quantity) * float64(it.PricePerUnit)
	}
	return updated, nil
}
