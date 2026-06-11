package domain

import (
	"testing"
)

func TestCountryParam_Required(t *testing.T) {
	valid := func() CountryParam {
		return CountryParam{
			CountryCode: "MX",
			Name:        "Mexico",
			Currency:    "MXN",
		}
	}

	tests := []struct {
		name    string
		param   CountryParam
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid country param - should pass",
			param:   valid(),
			wantErr: false,
		},
		{
			name: "empty country code - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.CountryCode = ""
				return cp
			}(),
			wantErr: true,
			errMsg:  "country code is required",
		},
		{
			name: "empty name - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.Name = ""
				return cp
			}(),
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty currency - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.Currency = ""
				return cp
			}(),
			wantErr: true,
			errMsg:  "currency is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.param.Required()
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

func TestCountryParam_Validate(t *testing.T) {
	valid := func() CountryParam {
		return CountryParam{
			CountryCode: "MX",
			Name:        "Mexico",
			Currency:    "MXN",
			MinWage:     1000,
			HealthRate:  5,
			PensionRate: 6,
		}
	}

	tests := []struct {
		name    string
		param   CountryParam
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - should pass",
			param:   valid(),
			wantErr: false,
		},
		{
			name: "invalid country code length - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.CountryCode = "MEX"
				return cp
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCountryCode.Error(),
		},
		{
			name: "invalid currency length - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.Currency = "MEXI"
				return cp
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCurrency.Error(),
		},
		{
			name: "zero min wage - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.MinWage = 0
				return cp
			}(),
			wantErr: true,
			errMsg:  "minimum wage must be greater than zero",
		},
		{
			name: "negative min wage - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.MinWage = -100
				return cp
			}(),
			wantErr: true,
			errMsg:  "minimum wage must be greater than zero",
		},
		{
			name: "zero health rate - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.HealthRate = 0
				return cp
			}(),
			wantErr: true,
			errMsg:  "health rate must be between 0 and 100",
		},
		{
			name: "health rate >= 100 - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.HealthRate = 100
				return cp
			}(),
			wantErr: true,
			errMsg:  "health rate must be between 0 and 100",
		},
		{
			name: "zero pension rate - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.PensionRate = 0
				return cp
			}(),
			wantErr: true,
			errMsg:  "pension rate must be between 0 and 100",
		},
		{
			name: "pension rate >= 100 - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.PensionRate = 100
				return cp
			}(),
			wantErr: true,
			errMsg:  "pension rate must be between 0 and 100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.param.Validate()
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

func TestCountryParam_Normalize(t *testing.T) {
	tests := []struct {
		name           string
		input          CountryParam
		expectedName   string
		expectedCode   string
		expectedCurr   string
	}{
		{
			name:           "trim spaces and uppercase",
			input:          CountryParam{CountryCode: "  mx  ", Name: "  Mexico  ", Currency: "  mxn  "},
			expectedCode:   "MX",
			expectedName:   "Mexico",
			expectedCurr:   "MXN",
		},
		{
			name:           "already normalized",
			input:          CountryParam{CountryCode: "CO", Name: "Colombia", Currency: "COP"},
			expectedCode:   "CO",
			expectedName:   "Colombia",
			expectedCurr:   "COP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.CountryCode != tt.expectedCode {
				t.Errorf("CountryCode: got '%s', want '%s'", tt.input.CountryCode, tt.expectedCode)
			}
			if tt.input.Name != tt.expectedName {
				t.Errorf("Name: got '%s', want '%s'", tt.input.Name, tt.expectedName)
			}
			if tt.input.Currency != tt.expectedCurr {
				t.Errorf("Currency: got '%s', want '%s'", tt.input.Currency, tt.expectedCurr)
			}
		})
	}
}

func TestCountryParam_ValidateAll(t *testing.T) {
	valid := func() CountryParam {
		return CountryParam{
			CountryCode: "MX",
			Name:        "Mexico",
			Currency:    "MXN",
			MinWage:     1000,
			HealthRate:  5,
			PensionRate: 6,
		}
	}

	tests := []struct {
		name    string
		param   CountryParam
		wantErr bool
	}{
		{
			name:    "valid - should pass",
			param:   valid(),
			wantErr: false,
		},
		{
			name:    "empty - should fail",
			param:   CountryParam{},
			wantErr: true,
		},
		{
			name: "missing currency - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.Currency = ""
				return cp
			}(),
			wantErr: true,
		},
		{
			name: "zero min wage - should fail",
			param: func() CountryParam {
				cp := valid()
				cp.MinWage = 0
				return cp
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.param.ValidateAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
