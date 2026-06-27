package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type EmployeeService struct {
	employeeRepo   domain.EmployeeRepo
	userRepo       domain.UserRepo
	departmentRepo domain.DepartmentRepo
	positionRepo   domain.PositionRepo
}

func NewEmployeeService(empRepo domain.EmployeeRepo, userRepo domain.UserRepo, deptRepo domain.DepartmentRepo, posRepo domain.PositionRepo) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   empRepo,
		userRepo:       userRepo,
		departmentRepo: deptRepo,
		positionRepo:   posRepo,
	}
}

func (s *EmployeeService) Create(ctx context.Context, emp *domain.Employee) error {
	if emp == nil {
		return errors.New("employee cannot be nil")
	}
	emp.Normalize()
	if err := emp.ValidateAll(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if emp.UserID != 0 {
		_, err := s.userRepo.GetByID(ctx, emp.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return domain.ErrUserNotFound
			}
			return fmt.Errorf("error validating user: %w", err)
		}
	}
	if emp.DepartmentID != 0 {
		_, err := s.departmentRepo.GetByID(ctx, emp.DepartmentID)
		if err != nil {
			if errors.Is(err, domain.ErrDepartmentNotFound) {
				return domain.ErrDepartmentNotFound
			}
			return fmt.Errorf("error validating department: %w", err)
		}
	}
	if emp.PositionID != 0 {
		_, err := s.positionRepo.GetByID(ctx, emp.PositionID)
		if err != nil {
			if errors.Is(err, domain.ErrPositionNotFound) {
				return domain.ErrPositionNotFound
			}
			return fmt.Errorf("error validating position: %w", err)
		}
	}
	if emp.SupervisorID != nil && *emp.SupervisorID != 0 {
		_, err := s.employeeRepo.GetByID(ctx, *emp.SupervisorID)
		if err != nil {
			if errors.Is(err, domain.ErrEmployeeNotFound) {
				return domain.ErrEmployeeNotFound
			}
			return fmt.Errorf("error validating supervisor: %w", err)
		}
	}
	existing, err := s.employeeRepo.GetByUserID(ctx, emp.UserID)
	if err != nil && !errors.Is(err, domain.ErrEmployeeNotFound) {
		return fmt.Errorf("error checking existing employee: %w", err)
	}
	if existing != nil {
		return domain.ErrEmployeeAlreadyExists
	}
	return s.employeeRepo.Create(ctx, emp)
}

func (s *EmployeeService) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return s.employeeRepo.GetByID(ctx, id)
}

func (s *EmployeeService) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.employeeRepo.GetByUserID(ctx, userID)
}

func (s *EmployeeService) List(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	return s.employeeRepo.List(ctx, page, limit)
}

func (s *EmployeeService) ListActive(ctx context.Context, page, limit int) ([]domain.Employee, int64, error) {
	return s.employeeRepo.ListActive(ctx, page, limit)
}

func (s *EmployeeService) ListBySupervisor(ctx context.Context, supervisorEmployeeID uint, page, limit int) ([]domain.Employee, int64, error) {
	return s.employeeRepo.ListBySupervisor(ctx, supervisorEmployeeID, page, limit)
}

func (s *EmployeeService) Update(ctx context.Context, emp *domain.Employee) error {
	if emp == nil || emp.ID <= 0 {
		return errors.New("employee cannot be nil or zero")
	}
	if emp.DepartmentID != 0 {
		_, err := s.departmentRepo.GetByID(ctx, emp.DepartmentID)
		if err != nil {
			if errors.Is(err, domain.ErrDepartmentNotFound) {
				return domain.ErrDepartmentNotFound
			}
			return fmt.Errorf("error validating department: %w", err)
		}
	}
	if emp.PositionID != 0 {
		_, err := s.positionRepo.GetByID(ctx, emp.PositionID)
		if err != nil {
			if errors.Is(err, domain.ErrPositionNotFound) {
				return domain.ErrPositionNotFound
			}
			return fmt.Errorf("error validating position: %w", err)
		}
	}
	if emp.SupervisorID != nil && *emp.SupervisorID != 0 {
		_, err := s.employeeRepo.GetByID(ctx, *emp.SupervisorID)
		if err != nil {
			if errors.Is(err, domain.ErrEmployeeNotFound) {
				return domain.ErrEmployeeNotFound
			}
			return fmt.Errorf("error validating supervisor: %w", err)
		}
	}
	existing, err := s.employeeRepo.GetByUserID(ctx, emp.UserID)
	if err != nil && !errors.Is(err, domain.ErrEmployeeNotFound) {
		return fmt.Errorf("error checking existing employee: %w", err)
	}
	if existing != nil && existing.ID != emp.ID {
		return domain.ErrEmployeeAlreadyExists
	}
	return s.employeeRepo.Update(ctx, emp)
}

func (s *EmployeeService) Delete(ctx context.Context, id uint) error {
	if id <= 0 {
		return errors.New("id employee cannot be zero")
	}
	return s.employeeRepo.Delete(ctx, id)
}
