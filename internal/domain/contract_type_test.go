package domain

import (
	"strings"
	"testing"
)

func TestContractType_Required(t *testing.T) {
	tests := []struct {
		name         string
		contractType ContractType
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid contract type - should pass",
			contractType: ContractType{
				Name:        "Full Time",
				CountryCode: "MX",
			},
			wantErr: false,
		},
		{
			name: "empty name - should fail",
			contractType: ContractType{
				Name:        "",
				CountryCode: "MX",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty country code - should fail",
			contractType: ContractType{
				Name:        "Full Time",
				CountryCode: "",
			},
			wantErr: true,
			errMsg:  "country code is required",
		},
		{
			name: "whitespace only name - should fail",
			contractType: ContractType{
				Name:        "   ",
				CountryCode: "MX",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contractType.Required()
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

func TestContractType_Validate(t *testing.T) {
	valid := func() ContractType {
		return ContractType{
			Name:        "Full Time",
			CountryCode: "MX",
		}
	}

	tests := []struct {
		name         string
		contractType ContractType
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid - should pass",
			contractType: valid(),
			wantErr:      false,
		},
		{
			name: "name > 100 chars - should fail",
			contractType: func() ContractType {
				ct := valid()
				ct.Name = strings.Repeat("a", 101)
				return ct
			}(),
			wantErr: true,
			errMsg:  "name must be at most 100 characters",
		},
		{
			name: "name exactly 100 chars - should pass",
			contractType: func() ContractType {
				ct := valid()
				ct.Name = strings.Repeat("a", 100)
				return ct
			}(),
			wantErr: false,
		},
		{
			name: "invalid country code length - should fail",
			contractType: func() ContractType {
				ct := valid()
				ct.CountryCode = "MEX"
				return ct
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCountryCode.Error(),
		},
		{
			name: "empty country code after normalize - should pass Validate",
			contractType: func() ContractType {
				ct := valid()
				ct.CountryCode = ""
				return ct
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCountryCode.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contractType.Validate()
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

func TestContractType_Normalize(t *testing.T) {
	tests := []struct {
		name           string
		input          ContractType
		expectedName   string
		expectedCode   string
	}{
		{
			name:           "trim spaces and uppercase country",
			input:          ContractType{Name: "  Full Time  ", CountryCode: "  mx  "},
			expectedName:   "Full Time",
			expectedCode:   "MX",
		},
		{
			name:           "already normalized",
			input:          ContractType{Name: "Full Time", CountryCode: "MX"},
			expectedName:   "Full Time",
			expectedCode:   "MX",
		},
		{
			name:           "lowercase country code",
			input:          ContractType{Name: "Part Time", CountryCode: "co"},
			expectedName:   "Part Time",
			expectedCode:   "CO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.Name != tt.expectedName {
				t.Errorf("Name: got '%s', want '%s'", tt.input.Name, tt.expectedName)
			}
			if tt.input.CountryCode != tt.expectedCode {
				t.Errorf("CountryCode: got '%s', want '%s'", tt.input.CountryCode, tt.expectedCode)
			}
		})
	}
}

func TestContractType_ValidateAll(t *testing.T) {
	valid := func() ContractType {
		return ContractType{
			Name:        "Full Time",
			CountryCode: "MX",
		}
	}

	tests := []struct {
		name         string
		contractType ContractType
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid - should pass",
			contractType: valid(),
			wantErr:      false,
		},
		{
			name: "empty name - should fail",
			contractType: func() ContractType {
				ct := valid()
				ct.Name = ""
				return ct
			}(),
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid country code - should fail",
			contractType: func() ContractType {
				ct := valid()
				ct.CountryCode = "MEX"
				return ct
			}(),
			wantErr: true,
			errMsg:  ErrInvalidCountryCode.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contractType.ValidateAll()
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
