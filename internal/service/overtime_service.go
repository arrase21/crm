package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/arrase21/crm/internal/domain"
)

type OvertimeService struct {
	overtimeRepo domain.OvertimeRepo
	employeeRepo domain.EmployeeRepo
}

func NewOvertimeService(overtimeRepo domain.OvertimeRepo, employeeRepo domain.EmployeeRepo) *OvertimeService {
	return &OvertimeService{
		overtimeRepo: overtimeRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *OvertimeService) Create(ctx context.Context, o *domain.Overtime) error {
	if o.EmployeeID == 0 {
		return errors.New("employee is required")
	}

	if err := s.checkSupervisorScope(ctx, o.EmployeeID); err != nil {
		return err
	}

	return s.overtimeRepo.Create(ctx, o)
}

func (s *OvertimeService) GetByID(ctx context.Context, id uint) (*domain.Overtime, error) {
	return s.overtimeRepo.GetByID(ctx, id)
}

func (s *OvertimeService) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Overtime, int64, error) {
	return s.overtimeRepo.ListByEmployee(ctx, employeeID, page, limit)
}

func (s *OvertimeService) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Overtime, int64, error) {
	return s.overtimeRepo.ListByDepartment(ctx, departmentID, page, limit)
}

func (s *OvertimeService) List(ctx context.Context, page, limit int) ([]domain.Overtime, int64, error) {
	return s.overtimeRepo.List(ctx, page, limit)
}

func (s *OvertimeService) Update(ctx context.Context, o *domain.Overtime) error {
	if o.ID == 0 {
		return errors.New("overtime id is required")
	}
	return s.overtimeRepo.Update(ctx, o)
}

func (s *OvertimeService) Delete(ctx context.Context, id uint) error {
	return s.overtimeRepo.Delete(ctx, id)
}

func (s *OvertimeService) checkSupervisorScope(ctx context.Context, targetEmployeeID uint) error {
	claims, ok := ctx.Value(domain.ClaimsKey).(*domain.Claims)
	if !ok {
		return nil
	}

	if slices.Contains(claims.Roles, "super_admin") {
		return nil
	}

	supervisorEmp, err := s.employeeRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("employee record not found: %w", err)
	}

	targetEmp, err := s.employeeRepo.GetByID(ctx, targetEmployeeID)
	if err != nil {
		return fmt.Errorf("target employee not found: %w", err)
	}

	if targetEmp.DepartmentID != supervisorEmp.DepartmentID {
		return errors.New("you can only manage employees in your own department")
	}

	return nil
}
