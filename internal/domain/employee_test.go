package domain

import (
	"testing"
)

func TestEmployee_Required(t *testing.T) {
	tests := []struct {
		name     string
		employee Employee
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid employee - should pass",
			employee: Employee{
				UserID: 1,
			},
			wantErr: false,
		},
		{
			name:     "empty user id - should fail",
			employee: Employee{UserID: 0},
			wantErr:  true,
			errMsg:   "user is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.employee.Required()
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
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEmployee_Validate(t *testing.T) {
	tests := []struct {
		name     string
		employee Employee
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid employee - should pass",
			employee: Employee{
				UserID:       1,
				DepartmentID: 1,
				PositionID:   1,
			},
			wantErr: false,
		},
		{
			name: "missing department - should fail",
			employee: Employee{
				UserID:     1,
				DepartmentID: 0,
				PositionID: 1,
			},
			wantErr: true,
			errMsg:  "department is required",
		},
		{
			name: "missing position - should fail",
			employee: Employee{
				UserID:       1,
				DepartmentID: 1,
				PositionID:   0,
			},
			wantErr: true,
			errMsg:  "position is required",
		},
		{
			name: "missing department and position - should fail on first",
			employee: Employee{
				UserID: 1,
			},
			wantErr: true,
			errMsg:  "department is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.employee.Validate()
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
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEmployee_ValidateAll(t *testing.T) {
	tests := []struct {
		name     string
		employee Employee
		wantErr  bool
	}{
		{
			name: "valid employee - should pass",
			employee: Employee{
				UserID:       1,
				DepartmentID: 1,
				PositionID:   1,
			},
			wantErr: false,
		},
		{
			name:     "missing user - should fail",
			employee: Employee{UserID: 0},
			wantErr:  true,
		},
		{
			name: "missing department - should fail",
			employee: Employee{
				UserID: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.employee.ValidateAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmployee_Normalize(t *testing.T) {
	e := Employee{UserID: 1}
	e.Normalize()
}
