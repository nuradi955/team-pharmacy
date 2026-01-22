package errs

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrItemNotFound     = errors.New("item not found")
	ErrMedicineNotFound = errors.New("medicine not found")
)
