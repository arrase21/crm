package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type DepartmentService struct {
	departmentRepo domain.DepartmentRepo
}

func NewDepartmentService(u domain.DepartmentRepo) *DepartmentService {
	return &DepartmentService{
		departmentRepo: u,
	}
}

func (s *DepartmentService) Create(ctx context.Context, departmet *domain.Department) error {
	if departmet == nil {
		return errors.New("invalid department id")
	}
	departmet.Normalize()
	if err := departmet.ValidateAll(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	existing, err := s.departmentRepo.GetByCode(ctx, departmet.Code)
	if err != nil && !errors.Is(err, domain.ErrDepartmentNotFound) {
		return fmt.Errorf("error checking existing code: %w ", err)
	}
	if existing != nil {
		return domain.ErrDepartmentCodeExists
	}

	existingName, err := s.departmentRepo.GetByName(ctx, departmet.Name)
	if err != nil && !errors.Is(err, domain.ErrDepartmentNotFound) {
		return fmt.Errorf("error checking existing name: %w", err)
	}
	if existingName != nil {
		return domain.ErrDepartmentNameExists
	}

	return s.departmentRepo.Create(ctx, departmet)
}

func (s *DepartmentService) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	if id == 0 {
		return nil, errors.New("invalid department id")
	}
	return s.departmentRepo.GetByID(ctx, id)
}

func (s *DepartmentService) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	if code == "" {
		return nil, errors.New("invalid department code, cannot be nil")
	}
	return s.departmentRepo.GetByCode(ctx, code)
}

func (s *DepartmentService) GetByName(ctx context.Context, name string) (*domain.Department, error) {
	if name == "" {
		return nil, errors.New("invalid department name, cannot be nil")
	}
	return s.departmentRepo.GetByName(ctx, name)
}

func (s *DepartmentService) List(ctx context.Context, page, limit int) ([]domain.Department, int64, error) {
	return s.departmentRepo.List(ctx, page, limit)
}

func (s *DepartmentService) Update(ctx context.Context, dept *domain.Department) error {
	if dept == nil || dept.ID == 0 {
		return errors.New("department cannot be nil or id is required")
	}
	existing, err := s.departmentRepo.GetByCode(ctx, dept.Code)
	if err != nil && !errors.Is(err, domain.ErrDepartmentNotFound) {
		return fmt.Errorf("error checking existing code: %w ", err)
	}
	if existing != nil && existing.ID != dept.ID {
		return domain.ErrDepartmentCodeExists
	}
	return s.departmentRepo.Update(ctx, dept)
}

func (s *DepartmentService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("error department id canno be nil or zero")
	}
	return s.departmentRepo.Delete(ctx, id)
}
