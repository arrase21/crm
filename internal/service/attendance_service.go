package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/arrase21/crm/internal/domain"
)

type AttendanceService struct {
	attendanceRepo domain.AttendanceRepo
	employeeRepo   domain.EmployeeRepo
}

func NewAttendanceService(attendanceRepo domain.AttendanceRepo, employeeRepo domain.EmployeeRepo) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *AttendanceService) Create(ctx context.Context, a *domain.Attendance) error {
	if a.EmployeeID == 0 {
		return errors.New("employee is required")
	}

	if err := s.checkSupervisorScope(ctx, a.EmployeeID); err != nil {
		return err
	}

	return s.attendanceRepo.Create(ctx, a)
}

func (s *AttendanceService) GetByID(ctx context.Context, id uint) (*domain.Attendance, error) {
	return s.attendanceRepo.GetByID(ctx, id)
}

func (s *AttendanceService) ListByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]domain.Attendance, int64, error) {
	return s.attendanceRepo.ListByEmployee(ctx, employeeID, page, limit)
}

func (s *AttendanceService) ListByDepartment(ctx context.Context, departmentID uint, page, limit int) ([]domain.Attendance, int64, error) {
	return s.attendanceRepo.ListByDepartment(ctx, departmentID, page, limit)
}

func (s *AttendanceService) List(ctx context.Context, page, limit int) ([]domain.Attendance, int64, error) {
	return s.attendanceRepo.List(ctx, page, limit)
}

func (s *AttendanceService) Update(ctx context.Context, a *domain.Attendance) error {
	if a.ID == 0 {
		return errors.New("attendance id is required")
	}
	return s.attendanceRepo.Update(ctx, a)
}

func (s *AttendanceService) Delete(ctx context.Context, id uint) error {
	return s.attendanceRepo.Delete(ctx, id)
}

func (s *AttendanceService) checkSupervisorScope(ctx context.Context, targetEmployeeID uint) error {
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
