package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type PositionService struct {
	positionRepo   domain.PositionRepo
	departmentRepo domain.DepartmentRepo
}

func NewPositionService(p domain.PositionRepo, deptRepo ...domain.DepartmentRepo) *PositionService {
	var dr domain.DepartmentRepo
	if len(deptRepo) > 0 {
		dr = deptRepo[0]
	}
	return &PositionService{
		positionRepo:   p,
		departmentRepo: dr,
	}
}

func (s *PositionService) Create(ctx context.Context, position *domain.Position) error {
	if position == nil {
		return errors.New("position cannot be nil")
	}

	// Validate DepartmentID if provided (only if departmentRepo is provided)
	if position.DepartmentID != 0 && s.departmentRepo != nil {
		_, err := s.departmentRepo.GetByID(ctx, position.DepartmentID)
		if err != nil {
			if errors.Is(err, domain.ErrDepartmentNotFound) {
				return domain.ErrDepartmentNotFound
			}
			return fmt.Errorf("error validating department: %w", err)
		}
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

func (s *PositionService) GetByIDWithDepartment(ctx context.Context, id uint) (*domain.Position, error) {
	if id == 0 {
		return nil, errors.New("invalid position id")
	}
	return s.positionRepo.GetByIDWithDepartment(ctx, id)
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

func (s *PositionService) ListByDepartment(ctx context.Context, departmentID uint) ([]domain.Position, int64, error) {
	if departmentID == 0 {
		return nil, 0, errors.New("department id is required")
	}
	return s.positionRepo.ListByDepartment(ctx, departmentID)
}

func (s *PositionService) CountByDepartment(ctx context.Context, departmentID uint) (int64, error) {
	if departmentID == 0 {
		return 0, errors.New("department id is required")
	}
	return s.positionRepo.CountByDepartment(ctx, departmentID)
}

func (s *PositionService) Update(ctx context.Context, position *domain.Position) error {
	if position == nil || position.ID == 0 {
		return errors.New("position cannot be nil or have zero id")
	}
	if position.DepartmentID != 0 && s.departmentRepo != nil {
		_, err := s.departmentRepo.GetByID(ctx, position.DepartmentID)
		if err != nil {
			if errors.Is(err, domain.ErrDepartmentNotFound) {
				return domain.ErrDepartmentNotFound
			}
			return fmt.Errorf("error validating department: %w", err)
		}
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
