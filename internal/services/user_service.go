package services

import (
	"errors"
	"strings"
	"team-pharmacy/internal/dto"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("пользователь не найден")

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers() ([]models.User, error)
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*models.User, error) {

	user := &models.User{
		FullName:       strings.TrimSpace(req.FullName),
		Email:          strings.TrimSpace(req.Email),
		Phone:          strings.TrimSpace(req.Phone),
		DefaultAddress: strings.TrimSpace(req.DefaultAddress),
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	user, err := s.users.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(id uint, req dto.UpdateUserRequest) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	user, err := s.users.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := s.applyUserUpdate(user, req); err != nil {
		return nil, err
	}

	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	if _, err := s.users.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.users.Delete(id)

}

func (s *userService) ListUsers() ([]models.User, error) {
	users, err := s.users.List()
	if err != nil {
		return nil, err
	}
	return users, nil

}

func (s *userService) applyUserUpdate(user *models.User, req dto.UpdateUserRequest) error {

	if req.FullName != nil {
		fullName := strings.TrimSpace(*req.FullName)
		if fullName == "" {
			return errors.New("имя не может быть пустым")
		}
		user.FullName = fullName
	}

	if req.Phone != nil {

		phone := strings.TrimSpace(*req.Phone)
		if phone == "" {
			return errors.New("поле phone не должно быть пустым")
		}
		user.Phone = phone
	}

	if req.DefaultAddress != nil {

		address := strings.TrimSpace(*req.DefaultAddress)
		if address == "" {
			return errors.New("поле address не должно быть пустым")
		}
		user.DefaultAddress = address

	}
	return nil
}
