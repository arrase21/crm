package service

import (
	"context"

	"github.com/arrase21/crm/internal/domain"
)

type RoleService struct {
	roleRepo domain.RoleRepo
	userRepo domain.UserRepo
}

func NewRoleService(roleRepo domain.RoleRepo, userRepo domain.UserRepo) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (s *RoleService) Assign(ctx context.Context, userID uint, roleName string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return err
	}

	return s.roleRepo.AssignRole(ctx, user.ID, role.ID)
}

func (s *RoleService) Unassign(ctx context.Context, userID uint, roleName string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return err
	}

	return s.roleRepo.UnassignRole(ctx, user.ID, role.ID)
}

func (s *RoleService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.roleRepo.List(ctx)
}

func (s *RoleService) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}
