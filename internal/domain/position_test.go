package domain

import (
	"testing"
)

func TestPosition_Required(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid position - should pass",
			position: Position{
				Name:        "gerente",
				Description: "Gerente de area",
			},
			wantErr: false,
		},
		{
			name: "empty name - should fail",
			position: Position{
				Name:        "",
				Description: "Descripcion",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "only spaces name - should fail",
			position: Position{
				Name:        "   ",
				Description: "Descripcion",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.position.Required()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("expected error message %q, got %q", tt.errMsg, err.Error())
			}
		})
	}
}

func TestPosition_Validate(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		wantErr  bool
	}{
		{
			name: "valid fields - should pass",
			position: Position{
				Name:        "Gerente",
				Description: "Gerente de area",
			},
			wantErr: false,
		},
		{
			name: "name too long - should fail",
			position: Position{
				Name:        string(make([]byte, 101)),
				Description: "Desc",
			},
			wantErr: true,
		},
		{
			name: "description too long - should fail",
			position: Position{
				Name:        "Name",
				Description: string(make([]byte, 256)),
			},
			wantErr: true,
		},
		{
			name: "max length name is valid",
			position: Position{
				Name:        string(make([]byte, 100)),
				Description: "Desc",
			},
			wantErr: false,
		},
		{
			name: "max length description is valid",
			position: Position{
				Name:        "Name",
				Description: string(make([]byte, 255)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.position.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestPosition_Normalize(t *testing.T) {
	tests := []struct {
		name         string
		inputName    string
		inputDesc    string
		expectedName string
		expectedDesc string
	}{
		{
			name:         "removes whitespace",
			inputName:    "  Gerente  ",
			inputDesc:    "  Descripcion  ",
			expectedName: "Gerente",
			expectedDesc: "Descripcion",
		},
		{
			name:         "empty stays empty",
			inputName:    "",
			inputDesc:    "",
			expectedName: "",
			expectedDesc: "",
		},
		{
			name:         "preserves normal text",
			inputName:    "Gerente",
			inputDesc:    "Descripcion",
			expectedName: "Gerente",
			expectedDesc: "Descripcion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				Name:        tt.inputName,
				Description: tt.inputDesc,
			}
			p.Normalize()
			if p.Name != tt.expectedName {
				t.Errorf("expected Name %q, got %q", tt.expectedName, p.Name)
			}
			if p.Description != tt.expectedDesc {
				t.Errorf("expected Description %q, got %q", tt.expectedDesc, p.Description)
			}
		})
	}
}

func TestPosition_ValidateAll(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		wantErr  bool
	}{
		{
			name: "valid position",
			position: Position{
				Name:        "Gerente",
				Description: "Gerente de area",
			},
			wantErr: false,
		},
		{
			name: "missing name - returns error",
			position: Position{
				Name:        "",
				Description: "Desc",
			},
			wantErr: true,
		},
		{
			name: "name too long - returns error",
			position: Position{
				Name:        string(make([]byte, 101)),
				Description: "Desc",
			},
			wantErr: true,
		},
		{
			name: "description too long - returns error",
			position: Position{
				Name:        "Name",
				Description: string(make([]byte, 256)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.position.ValidateAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
