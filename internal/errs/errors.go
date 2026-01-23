package errs

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrItemNotFound            = errors.New("item not found")
	ErrMedicineNotFound        = errors.New("medicine not found")
	ErrCartNotFound            = errors.New("cart not found")
	ErrCartIsEmpty             = errors.New("cart is empty")
	ErrInvalidID               = errors.New("invalid ID")
	ErrOrderNotFound           = errors.New("order not found")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
