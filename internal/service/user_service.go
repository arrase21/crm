package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(u domain.UserRepo) *UserService {
	return &UserService{
		userRepo: u,
	}
}

func (s *UserService) Create(ctx context.Context, usr *domain.User) error {
	if usr == nil {
		return errors.New("user cannot be nil")
	}
	usr.Normalize()
	if err := usr.ValidateAll(); err != nil {
		return fmt.Errorf("Error checking existing user: %w ", err)
	}
	existing, err := s.userRepo.GetByDNI(ctx, usr.Dni)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("Error checking existing user: %w ", err)
	}
	if existing != nil {
		return domain.ErrDniAlreadyExist
	}
	return s.userRepo.Create(ctx, usr)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByDNI(ctx context.Context, dni string) (*domain.User, error) {
	if dni == "" {
		return nil, errors.New("dni cannot be empty")
	}
	return s.userRepo.GetByDNI(ctx, dni)
}

func (s *UserService) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	return s.userRepo.List(ctx, page, limit)
}

func (s *UserService) Update(ctx context.Context, usr *domain.User) error {
	if usr == nil || usr.ID == 0 {
		return errors.New("user cannot be nil or user id is required")
	}
	usr.Normalize()
	if err := usr.ValidateAll(); err != nil {
		return fmt.Errorf("validation error in domain: %w", err)
	}
	existing, err := s.userRepo.GetByDNI(ctx, usr.Dni)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("error checking existing dni: %w", err)
	}
	if existing != nil && existing.ID != usr.ID {
		return domain.ErrDniAlreadyExist
	}
	return s.userRepo.Update(ctx, usr)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}
	return s.userRepo.Delete(ctx, id)
}
