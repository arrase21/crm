package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrase21/crm/internal/domain"
)

type EmployeeContractService struct {
	contractRepo     domain.EmployeeContractRepo
	employeeRepo     domain.EmployeeRepo
	contractTypeRepo domain.ContractTypeRepo
}

func NewEmployeeContractService(
	contractRepo domain.EmployeeContractRepo,
	employeeRepo domain.EmployeeRepo,
	contractTypeRepo domain.ContractTypeRepo,
) *EmployeeContractService {
	return &EmployeeContractService{
		contractRepo:     contractRepo,
		employeeRepo:     employeeRepo,
		contractTypeRepo: contractTypeRepo,
	}
}

func (s *EmployeeContractService) Create(ctx context.Context, ec *domain.EmployeeContract) error {
	if ec == nil {
		return errors.New("contract cannot be nil")
	}
	ec.Normalize()
	if err := ec.ValidateAll(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	_, err := s.employeeRepo.GetByID(ctx, ec.EmployeeID)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return domain.ErrEmployeeNotFound
		}
		return fmt.Errorf("error validating employee: %w", err)
	}
	if ec.ContractTypeID != 0 {
		_, err := s.contractTypeRepo.GetByID(ctx, ec.ContractTypeID)
		if err != nil {
			if errors.Is(err, domain.ErrContractTypeNotFound) {
				return domain.ErrContractTypeNotFound
			}
			return fmt.Errorf("error validating contract type: %w", err)
		}
	}
	return s.contractRepo.Create(ctx, ec)
}

func (s *EmployeeContractService) GetByID(ctx context.Context, id uint) (*domain.EmployeeContract, error) {
	if id == 0 {
		return nil, errors.New("invalid contract id")
	}
	return s.contractRepo.GetByID(ctx, id)
}

func (s *EmployeeContractService) GetByEmployeeID(ctx context.Context, employeeID uint) ([]domain.EmployeeContract, error) {
	if employeeID == 0 {
		return nil, errors.New("invalid employee id")
	}
	return s.contractRepo.GetByEmployeeID(ctx, employeeID)
}

func (s *EmployeeContractService) List(ctx context.Context, page, limit int) ([]domain.EmployeeContract, int64, error) {
	return s.contractRepo.List(ctx, page, limit)
}

func (s *EmployeeContractService) Update(ctx context.Context, ec *domain.EmployeeContract) error {
	if ec == nil || ec.ID == 0 {
		return errors.New("contract cannot be nil or have zero id")
	}
	if ec.ContractTypeID != 0 {
		_, err := s.contractTypeRepo.GetByID(ctx, ec.ContractTypeID)
		if err != nil {
			if errors.Is(err, domain.ErrContractTypeNotFound) {
				return domain.ErrContractTypeNotFound
			}
			return fmt.Errorf("error validating contract type: %w", err)
		}
	}
	return s.contractRepo.Update(ctx, ec)
}

func (s *EmployeeContractService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid contract id")
	}
	return s.contractRepo.Delete(ctx, id)
}


