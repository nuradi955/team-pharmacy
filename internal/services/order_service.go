package services

import (
	"errors"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/errs"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"

	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(userID uint, req *dto.OrderCreateRequest) (*dto.OrderResponse, error)
	GetByID(orderID uint) (*dto.OrderResponse, error)
	GetListOrders(userID uint) ([]dto.OrderShortResponse, error)
	UpdateOrder(orderID uint, req *dto.OrderStatusRequest) error
}

type orderService struct {
	orderRepo    repository.OrderRepository
	userRepo     repository.UserRepository
	cartRepo     repository.CartRepository
	medicineRepo repository.MedicineRepository
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository,
	cartRepo repository.CartRepository, medicineRepo repository.MedicineRepository) OrderService {

	return &orderService{orderRepo: orderRepo, userRepo: userRepo, cartRepo: cartRepo, medicineRepo: medicineRepo}
}

func (s *orderService) CreateOrder(userID uint, req *dto.OrderCreateRequest) (*dto.OrderResponse, error) {
	if userID == 0 {
		return nil, errs.ErrInvalidID
	}

	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCartNotFound
		}
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errs.ErrCartIsEmpty
	}

	var (
		orderItems []models.OrderItem
		totalPrice int64
	)

	for _, cartItem := range cart.Items {
		lineTotal := int64(cartItem.Quantity) * cartItem.PricePerUnit

		orderItems = append(orderItems, models.OrderItem{
			MedicineID:   cartItem.MedicineID,
			MedicineName: cartItem.Medicine.Name,
			Quantity:     cartItem.Quantity,
			PricePerUnit: cartItem.PricePerUnit,
			LineTotal:    int64(lineTotal),
		})
		totalPrice += lineTotal
	}

	order := models.Order{
		UserID:          userID,
		Status:          models.OrderStatusPendingPayment,
		TotalPrice:      totalPrice,
		DiscountTotal:   0,
		FinalPrice:      totalPrice,
		DeliveryAddress: req.DeliveryAddress,
		Comment:         req.Comment,
		Items:           orderItems,
	}

	if err := s.orderRepo.CreateOrderWithClearCart(&order, cart.ID); err != nil {
		return nil, err
	}

	return &dto.OrderResponse{UserID: userID,
		Status:          models.OrderStatusPendingPayment,
		TotalPrice:      order.TotalPrice,
		DiscountTotal:   order.DiscountTotal,
		FinalPrice:      order.FinalPrice,
		DeliveryAddress: order.DeliveryAddress,
		Comment:         &order.Comment,
		Items:           orderItems,
	}, nil
}

func (s *orderService) GetByID(orderID uint) (*dto.OrderResponse, error) {
	if orderID == 0 {
		return nil, errs.ErrInvalidID
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{UserID: order.UserID,
		Status:          order.Status,
		TotalPrice:      order.TotalPrice,
		DiscountTotal:   order.DiscountTotal,
		FinalPrice:      order.FinalPrice,
		DeliveryAddress: order.DeliveryAddress,
		Comment:         &order.Comment,
		Items:           order.Items,
		// Payments:        order.Payments,
	}, nil

}

func (s *orderService) GetListOrders(userID uint) ([]dto.OrderShortResponse, error) {

	if userID == 0 {
		return nil, errs.ErrInvalidID
	}

	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	orders, err := s.orderRepo.GetListOrders(userID)
	if err != nil {
		return nil, err
	}
	shortListOrders := make([]dto.OrderShortResponse, 0, len(orders))

	for _, order := range orders {
		shortListOrders = append(shortListOrders, dto.OrderShortResponse{
			ID:         order.ID,
			Status:     order.Status,
			FinalPrice: order.FinalPrice,
			CreatedAt:  order.CreatedAt,
		})
	}
	return shortListOrders, nil

}

func (s *orderService) UpdateOrder(orderID uint, req *dto.OrderStatusRequest) error {
	if orderID == 0 {
		return errs.ErrInvalidID
	}

	if req == nil || req.Status == nil {
		return errs.ErrInvalidStatus
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	newStatus := *req.Status
	if !models.CanChangeOrderStatus(order.Status, newStatus) {
		return errs.ErrInvalidStatus
	}

	return s.orderRepo.UpdateOrder(orderID, &newStatus)
}
