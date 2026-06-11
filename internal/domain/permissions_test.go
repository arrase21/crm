package domain

import (
	"testing"
)

func TestPermissionsList(t *testing.T) {
	expected := map[string]map[string]bool{
		"users":       {"create": true, "read": true, "update": true, "delete": true},
		"employees":   {"create": true, "read": true, "update": true, "delete": true},
		"contracts":   {"create": true, "read": true, "update": true, "delete": true},
		"payroll":     {"calculate": true, "read": true},
		"departments": {"create": true, "read": true, "update": true, "delete": true},
		"positions":   {"create": true, "read": true, "update": true, "delete": true},
		"attendance":  {"create": true, "read": true, "update": true, "delete": true},
		"overtime":    {"create": true, "read": true, "update": true, "delete": true},
		"roles":       {"assign": true},
	}

	if len(Permissions) != 31 {
		t.Errorf("expected 33 permissions, got %d", len(Permissions))
	}

	got := make(map[string]map[string]bool)
	for _, p := range Permissions {
		if got[p.Resource] == nil {
			got[p.Resource] = make(map[string]bool)
		}
		got[p.Resource][p.Action] = true
	}

	for resource, actions := range expected {
		for action := range actions {
			if !got[resource][action] {
				t.Errorf("missing permission: %s %s", resource, action)
			}
		}
	}
}

func TestRole_Fields(t *testing.T) {
	role := Role{
		ID:       1,
		TenantID: 1,
		Name:     "super_admin",
	}
	if role.ID != 1 {
		t.Errorf("expected id 1, got %d", role.ID)
	}
	if role.TenantID != 1 {
		t.Errorf("expected tenant id 1, got %d", role.TenantID)
	}
	if role.Name != "super_admin" {
		t.Errorf("expected name 'super_admin', got '%s'", role.Name)
	}
}

func TestRolePermission_Fields(t *testing.T) {
	rp := RolePermission{
		RoleID:       1,
		PermissionID: 2,
	}
	if rp.RoleID != 1 {
		t.Errorf("expected role id 1, got %d", rp.RoleID)
	}
	if rp.PermissionID != 2 {
		t.Errorf("expected permission id 2, got %d", rp.PermissionID)
	}
}

func TestUserRole_Fields(t *testing.T) {
	ur := UserRole{
		UserID: 1,
		RoleID: 2,
	}
	if ur.UserID != 1 {
		t.Errorf("expected user id 1, got %d", ur.UserID)
	}
	if ur.RoleID != 2 {
		t.Errorf("expected role id 2, got %d", ur.RoleID)
	}
}
