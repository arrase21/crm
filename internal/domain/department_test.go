package domain

import (
	"strings"
	"testing"
)

func TestDepartment_Required(t *testing.T) {
	tests := []struct {
		name       string
		department Department
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid department - should pass",
			department: Department{
				Name:     "accounting",
				Code:     "acc",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "empty name - should fail",
			department: Department{
				Name:     "",
				Code:     "acc",
				IsActive: true,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty code - should fail",
			department: Department{
				Name:     "accounting",
				Code:     "",
				IsActive: true,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "only spaces name - should fail",
			department: Department{
				Name:     "   ",
				Code:     "acc",
				IsActive: true,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "only spaces code - should fail",
			department: Department{
				Name:     "accounting",
				Code:     "   ",
				IsActive: true,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.department.Required()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error '%s' but got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: '%v'", err)
				}
			}
		})
	}
}

func TestDepartment_Validate(t *testing.T) {
	validDepartment := func() Department {
		return Department{
			Name:     "accounting",
			Code:     "acc",
			IsActive: true,
		}
	}

	tests := []struct {
		name       string
		department Department
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid department - should pass",
			department: validDepartment(),
			wantErr:    false,
		},
		{
			name: "name exactly 100 chars - should pass",
			department: func() Department {
				d := validDepartment()
				d.Name = strings.Repeat("a", 100)
				return d
			}(),
			wantErr: false,
		},
		{
			name: "name > 100 chars - should fail",
			department: func() Department {
				d := validDepartment()
				d.Name = strings.Repeat("a", 101)
				return d
			}(),
			wantErr: true,
			errMsg:  "name must be at most 100 characters",
		},
		{
			name: "code exactly 20 chars - should pass",
			department: func() Department {
				d := validDepartment()
				d.Code = strings.Repeat("c", 20)
				return d
			}(),
			wantErr: false,
		},
		{
			name: "code > 20 chars - should fail",
			department: func() Department {
				d := validDepartment()
				d.Code = strings.Repeat("c", 21)
				return d
			}(),
			wantErr: true,
			errMsg:  "code must be at most 20 characters",
		},
		{
			name: "empty name - should pass (Required checks existence)",
			department: func() Department {
				d := validDepartment()
				d.Name = ""
				return d
			}(),
			wantErr: false,
		},
		{
			name: "empty code - should pass (Required checks existence)",
			department: func() Department {
				d := validDepartment()
				d.Code = ""
				return d
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.department.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error '%s' but got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got '%v'", err)
				}
			}
		})
	}
}

func TestDepartment_Normalize(t *testing.T) {
	tests := []struct {
		name       string
		department Department
		wantName   string
		wantCode   string
	}{
		{
			name:       "no spaces - should remain unchanged",
			department: Department{Name: "accounting", Code: "ACC"},
			wantName:   "accounting",
			wantCode:   "ACC",
		},
		{
			name:       "leading spaces - should trim",
			department: Department{Name: "  accounting", Code: "  acc"},
			wantName:   "accounting",
			wantCode:   "acc",
		},
		{
			name:       "trailing spaces - should trim",
			department: Department{Name: "accounting  ", Code: "acc  "},
			wantName:   "accounting",
			wantCode:   "acc",
		},
		{
			name:       "leading and trailing spaces - should trim both",
			department: Department{Name: "  accounting  ", Code: "  acc  "},
			wantName:   "accounting",
			wantCode:   "acc",
		},
		{
			name:       "multiple spaces in middle - should keep them",
			department: Department{Name: "accounting  dept", Code: "acc"},
			wantName:   "accounting  dept",
			wantCode:   "acc",
		},
		{
			name:       "empty strings - should stay empty",
			department: Department{Name: "", Code: ""},
			wantName:   "",
			wantCode:   "",
		},
		{
			name:       "only spaces - should become empty",
			department: Department{Name: "   ", Code: "   "},
			wantName:   "",
			wantCode:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.department.Normalize()
			if tt.department.Name != tt.wantName {
				t.Errorf("Name: expected '%s' but got '%s'", tt.wantName, tt.department.Name)
			}
			if tt.department.Code != tt.wantCode {
				t.Errorf("Code: expected '%s' but got '%s'", tt.wantCode, tt.department.Code)
			}
		})
	}
}

func TestDepartment_ValidateAll(t *testing.T) {
	validDepartment := func() Department {
		return Department{
			Name:     "accounting",
			Code:     "acc",
			IsActive: true,
		}
	}

	tests := []struct {
		name       string
		department Department
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid department - should pass",
			department: validDepartment(),
			wantErr:    false,
		},
		{
			name: "empty name - should fail with name error",
			department: func() Department {
				d := validDepartment()
				d.Name = ""
				return d
			}(),
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty code - should fail with code error",
			department: func() Department {
				d := validDepartment()
				d.Code = ""
				return d
			}(),
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "name > 100 chars - should fail with length error",
			department: func() Department {
				d := validDepartment()
				d.Name = strings.Repeat("a", 101)
				return d
			}(),
			wantErr: true,
			errMsg:  "name must be at most 100 characters",
		},
		{
			name: "code > 20 chars - should fail with length error",
			department: func() Department {
				d := validDepartment()
				d.Code = strings.Repeat("c", 21)
				return d
			}(),
			wantErr: true,
			errMsg:  "code must be at most 20 characters",
		},
		{
			name: "only spaces - should fail after normalize (spaces trimmed)",
			department: func() Department {
				d := validDepartment()
				d.Name = "   "
				d.Code = "   "
				return d
			}(),
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name with leading/trailing spaces and valid length after trim - should pass",
			department: func() Department {
				d := validDepartment()
				d.Name = "  " + strings.Repeat("a", 96) + "  "
				d.Code = "acc"
				return d
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.department.ValidateAll()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error '%s' but got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got '%v'", err)
				}
			}
		})
	}
}
