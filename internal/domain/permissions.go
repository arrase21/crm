package domain

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"not null;index" json:"tenant_id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex:idx_role_name" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Permission struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Resource string `gorm:"size:50;not null;uniqueIndex:idx_perm_resource_action" json:"resource"`
	Action   string `gorm:"size:50;not null;uniqueIndex:idx_perm_resource_action" json:"action"`
	Label    string `gorm:"size:100" json:"label"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
}

var Permissions = []Permission{
	{Resource: "users", Action: "create"}, {Resource: "users", Action: "read"},
	{Resource: "users", Action: "update"}, {Resource: "users", Action: "delete"},
	{Resource: "employees", Action: "create"}, {Resource: "employees", Action: "read"},
	{Resource: "employees", Action: "update"}, {Resource: "employees", Action: "delete"},
	{Resource: "contracts", Action: "create"}, {Resource: "contracts", Action: "read"},
	{Resource: "contracts", Action: "update"}, {Resource: "contracts", Action: "delete"},
	{Resource: "payroll", Action: "calculate"}, {Resource: "payroll", Action: "read"},
	{Resource: "departments", Action: "create"}, {Resource: "departments", Action: "read"},
	{Resource: "departments", Action: "update"}, {Resource: "departments", Action: "delete"},
	{Resource: "positions", Action: "create"}, {Resource: "positions", Action: "read"},
	{Resource: "positions", Action: "update"}, {Resource: "positions", Action: "delete"},
	{Resource: "attendance", Action: "create"}, {Resource: "attendance", Action: "read"},
	{Resource: "attendance", Action: "update"}, {Resource: "attendance", Action: "delete"},
	{Resource: "overtime", Action: "create"}, {Resource: "overtime", Action: "read"},
	{Resource: "overtime", Action: "update"}, {Resource: "overtime", Action: "delete"},
	{Resource: "roles", Action: "assign"},
}
