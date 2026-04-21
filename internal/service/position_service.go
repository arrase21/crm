package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type PositionService struct {
	positionRepo domain.PositionRepo
}

func NewPositionService(p domain.PositionRepo) *PositionService {
	return &PositionService{
		positionRepo: p,
	}
}

func (s *PositionService) Create(ctx context.Context, position *domain.Position) error {
	if position == nil {
		return errors.New("position cannot be nil")
	}
	position.Normalize()
	if err := position.ValidateAll(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	existing, err := s.positionRepo.GetByName(ctx, position.Name)
	if err != nil && !errors.Is(err, domain.ErrPositionNotFound) {
		return fmt.Errorf("error checking existing position: %w", err)
	}
	if existing != nil {
		return domain.ErrPositionNameExists
	}

	return s.positionRepo.Create(ctx, position)
}

func (s *PositionService) GetByID(ctx context.Context, id uint) (*domain.Position, error) {
	if id == 0 {
		return nil, errors.New("invalid position id")
	}
	return s.positionRepo.GetByID(ctx, id)
}

func (s *PositionService) GetByName(ctx context.Context, name string) (*domain.Position, error) {
	if name == "" {
		return nil, errors.New("position name cannot be empty")
	}
	return s.positionRepo.GetByName(ctx, name)
}

func (s *PositionService) List(ctx context.Context, page, limit int) ([]domain.Position, int64, error) {
	return s.positionRepo.List(ctx, page, limit)
}

func (s *PositionService) Update(ctx context.Context, position *domain.Position) error {
	if position == nil || position.ID == 0 {
		return errors.New("position cannot be nil or have zero id")
	}
	existing, err := s.positionRepo.GetByName(ctx, position.Name)
	if err != nil && !errors.Is(err, domain.ErrPositionNotFound) {
		return fmt.Errorf("error checking existing position: %w", err)
	}
	if existing != nil && existing.ID != position.ID {
		return domain.ErrPositionNameExists
	}
	return s.positionRepo.Update(ctx, position)
}

func (s *PositionService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid position id")
	}
	return s.positionRepo.Delete(ctx, id)
}
