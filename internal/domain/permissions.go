package domain

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	TenantID  uint      `gorm:"not null;index"`
	Name      string    `gorm:"size:50;not null;uniqueIndex:idx_role_name"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Permission struct {
	ID       uint   `gorm:"primaryKey"`
	Resource string `gorm:"size:50;not null;uniqueIndex:idx_perm_resource_action,composite:resource_action"`
	Action   string `gorm:"size:50;not null;uniqueIndex:idx_perm_resource_action,composite:resource_action"`
	Label    string `gorm:"size:100"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
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
