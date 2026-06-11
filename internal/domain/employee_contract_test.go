package domain

import (
	"testing"
	"time"
)

func TestEmployeeContract_Required(t *testing.T) {
	valid := func() EmployeeContract {
		return EmployeeContract{
			EmployeeID:  1,
			CountryCode: "MX",
			Currency:    "MXN",
		}
	}

	tests := []struct {
		name     string
		contract EmployeeContract
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid contract - should pass",
			contract: valid(),
			wantErr:  false,
		},
		{
			name: "empty employee id - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.EmployeeID = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "employee is required",
		},
		{
			name: "empty country code - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.CountryCode = ""
				return ec
			}(),
			wantErr: true,
			errMsg:  "country code is required",
		},
		{
			name: "empty currency - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.Currency = ""
				return ec
			}(),
			wantErr: true,
			errMsg:  "currency is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contract.Required()
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

func TestEmployeeContract_Validate(t *testing.T) {
	now := time.Now()
	valid := func() EmployeeContract {
		return EmployeeContract{
			EmployeeID:      1,
			CountryCode:     "MX",
			Currency:        "MXN",
			BaseSalary:      10000,
			StartDate:       now.AddDate(0, -1, 0),
			WorkHoursPerDay: 8,
			WorkDaysPerWeek: 5,
		}
	}

	tests := []struct {
		name     string
		contract EmployeeContract
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid contract - should pass",
			contract: valid(),
			wantErr:  false,
		},
		{
			name: "zero base salary - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.BaseSalary = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "base salary must be greater than zero",
		},
		{
			name: "negative base salary - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.BaseSalary = -100
				return ec
			}(),
			wantErr: true,
			errMsg:  "base salary must be greater than zero",
		},
		{
			name: "invalid country code length - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.CountryCode = "MEX"
				return ec
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCountryCode.Error(),
		},
		{
			name: "invalid currency length - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.Currency = "MEXICAN"
				return ec
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCurrency.Error(),
		},
		{
			name: "zero start date - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.StartDate = time.Time{}
				return ec
			}(),
			wantErr: true,
			errMsg:  "start date is required",
		},
		{
			name: "end date before start date - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				end := ec.StartDate.AddDate(0, -1, 0)
				ec.EndDate = &end
				return ec
			}(),
			wantErr: true,
			errMsg:  "end date must be after start date",
		},
		{
			name: "end date after start date - should pass",
			contract: func() EmployeeContract {
				ec := valid()
				end := ec.StartDate.AddDate(0, 1, 0)
				ec.EndDate = &end
				return ec
			}(),
			wantErr: false,
		},
		{
			name: "work hours per day zero - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.WorkHoursPerDay = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "work hours per day must be between 1 and 24",
		},
		{
			name: "work hours per day > 24 - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.WorkHoursPerDay = 25
				return ec
			}(),
			wantErr: true,
			errMsg:  "work hours per day must be between 1 and 24",
		},
		{
			name: "work days per week zero - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.WorkDaysPerWeek = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "work days per week must be between 1 and 7",
		},
		{
			name: "work days per week > 7 - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.WorkDaysPerWeek = 8
				return ec
			}(),
			wantErr: true,
			errMsg:  "work days per week must be between 1 and 7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contract.Validate()
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

func TestEmployeeContract_Normalize(t *testing.T) {
	tests := []struct {
		name             string
		input            EmployeeContract
		expectedCode     string
		expectedCurrency string
	}{
		{
			name:             "trim spaces and uppercase",
			input:            EmployeeContract{CountryCode: "  mx  ", Currency: "  mxn  "},
			expectedCode:     "MX",
			expectedCurrency: "MXN",
		},
		{
			name:             "already normalized",
			input:            EmployeeContract{CountryCode: "CO", Currency: "COP"},
			expectedCode:     "CO",
			expectedCurrency: "COP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.CountryCode != tt.expectedCode {
				t.Errorf("CountryCode: got '%s', want '%s'", tt.input.CountryCode, tt.expectedCode)
			}
			if tt.input.Currency != tt.expectedCurrency {
				t.Errorf("Currency: got '%s', want '%s'", tt.input.Currency, tt.expectedCurrency)
			}
		})
	}
}

func TestEmployeeContract_ValidateAll(t *testing.T) {
	now := time.Now()
	valid := func() EmployeeContract {
		return EmployeeContract{
			EmployeeID:      1,
			CountryCode:     "MX",
			Currency:        "MXN",
			BaseSalary:      10000,
			StartDate:       now.AddDate(0, -1, 0),
			WorkHoursPerDay: 8,
			WorkDaysPerWeek: 5,
		}
	}

	tests := []struct {
		name     string
		contract EmployeeContract
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid - should pass",
			contract: valid(),
			wantErr:  false,
		},
		{
			name: "missing employee - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.EmployeeID = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "employee is required",
		},
		{
			name: "missing country code - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.CountryCode = ""
				return ec
			}(),
			wantErr: true,
			errMsg:  "country code is required",
		},
		{
			name: "zero base salary - should fail",
			contract: func() EmployeeContract {
				ec := valid()
				ec.BaseSalary = 0
				return ec
			}(),
			wantErr: true,
			errMsg:  "base salary must be greater than zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contract.ValidateAll()
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
