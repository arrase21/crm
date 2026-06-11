package main

import (
	"github.com/arrase21/crm/internal/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20260610-initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&domain.User{},
					&domain.Department{},
					&domain.Position{},
					&domain.Employee{},
					&domain.ContractType{},
					&domain.EmployeeContract{},
					&domain.CountryParam{},
					&domain.PayrollRecord{},
					&domain.Attendance{},
					&domain.Overtime{},
					&domain.Role{},
					&domain.Permission{},
					&domain.UserRole{},
					&domain.RolePermission{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"user_roles",
					"role_permissions",
					"permissions",
					"roles",
					"overtime",
					"attendances",
					"payroll_records",
					"country_params",
					"employee_contracts",
					"contract_types",
					"employees",
					"positions",
					"departments",
					"users",
				)
			},
		},
	}
}
