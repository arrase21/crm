package domain

import (
	"testing"
	"time"
)

// ============================================
// DOMAIN TESTS - User Validation & Logic
// ============================================
//
// ¿Qué testea?: La lógica del dominio puramente.
// No tiene dependencias externas (ni DB, ni HTTP).
// Son los tests más rápidos y los que más debes tener.
//
// Estructura recomendada:
//   - Table-Driven Tests: para probar múltiples casos de entrada
//   - Arrange-Act-Assert (AAA): organiza cada test en 3 fases
func TestUser_Required(t *testing.T) {
	// Table-Driven Tests: la mejor forma de probar múltiples escenarios
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user - should pass",
			user: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			wantErr: false,
		},
		{
			name: "empty first name - should fail",
			user: User{
				FirstName: "",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			wantErr: true,
			errMsg:  "first name is required",
		},
		{
			name: "empty last name - should fail",
			user: User{
				FirstName: "Juan",
				LastName:  "",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			wantErr: true,
			errMsg:  "last name is required",
		},
		{
			name: "empty email - should fail",
			user: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "empty dni - should fail",
			user: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "",
				Phone:     "1234567890",
			},
			wantErr: true,
			errMsg:  "dni is required",
		},
		{
			name: "empty phone - should fail",
			user: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "",
			},
			wantErr: true,
			errMsg:  "phone is required",
		},
		{
			name: "whitespace only - should fail (trim spaces)",
			user: User{
				FirstName: "   ",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			wantErr: true,
			errMsg:  "first name is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: preparación (en este caso el setup está en el struct)
			// Act: ejecutamos la acción
			err := tt.user.Required()
			// Assert: verificamos el resultado
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
func TestUser_Validate(t *testing.T) {
	// Helper para crear un usuario válido
	validUser := func() User {
		return User{
			FirstName: "Juan",
			LastName:  "Pérez",
			Email:     "juan@example.com",
			Dni:       "12345678",
			Phone:     "1234567890",
			Gender:    "M",
			BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), // 36 años
		}
	}
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid user - should pass",
			user:    validUser(),
			wantErr: false,
		},
		{
			name: "invalid gender - should fail",
			user: func() User {
				u := validUser()
				u.Gender = "X"
				return u
			}(),
			wantErr: true,
			errMsg:  "Invalid option",
		},
		{
			name: "empty gender - should fail",
			user: func() User {
				u := validUser()
				u.Gender = ""
				return u
			}(),
			wantErr: true,
			errMsg:  "Invalid option",
		},
		{
			name: "future birthday - should fail",
			user: func() User {
				u := validUser()
				u.BirthDay = time.Now().Add(24 * time.Hour)
				return u
			}(),
			wantErr: true,
			errMsg:  "invalid birthday",
		},
		{
			name: "minor user - should fail",
			user: func() User {
				u := validUser()
				u.BirthDay = time.Now().AddDate(-17, 0, 0) // Hace 17 años
				return u
			}(),
			wantErr: true,
			errMsg:  "user must be over 18",
		},
		{
			name: "exactly 18 - should pass",
			user: func() User {
				u := validUser()
				u.BirthDay = time.Now().AddDate(-18, 0, -1) // 18 años menos 1 día
				return u
			}(),
			wantErr: false,
		},
		{
			name: "zero birthday - should fail",
			user: func() User {
				u := validUser()
				u.BirthDay = time.Time{}
				return u
			}(),
			wantErr: true,
			errMsg:  "invalid birthday",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
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
func TestUser_ValidateAll(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid complete user - should pass",
			user: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
				Gender:    "M",
				BirthDay:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "missing required + invalid gender - should fail on first error",
			user: User{
				FirstName: "",
				Gender:    "X",
			},
			wantErr: true,
		},
		{
			name:    "all fields missing - should fail",
			user:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestUser_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    User
		expected User
	}{
		{
			name: "normalize lowercase email and trim spaces",
			input: User{
				FirstName: "  Juan  ",
				LastName:  "  Pérez  ",
				Email:     "  JUAN@EXAMPLE.COM  ",
				Dni:       "  12345678  ",
				Phone:     "  1234567890  ",
			},
			expected: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
		},
		{
			name: "already normalized - no changes",
			input: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
			expected: User{
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Dni:       "12345678",
				Phone:     "1234567890",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.FirstName != tt.expected.FirstName {
				t.Errorf("FirstName: got %s, want %s", tt.input.FirstName, tt.expected.FirstName)
			}
			if tt.input.LastName != tt.expected.LastName {
				t.Errorf("LastName: got %s, want %s", tt.input.LastName, tt.expected.LastName)
			}
			if tt.input.Email != tt.expected.Email {
				t.Errorf("Email: got %s, want %s", tt.input.Email, tt.expected.Email)
			}
			if tt.input.Dni != tt.expected.Dni {
				t.Errorf("Dni: got %s, want %s", tt.input.Dni, tt.expected.Dni)
			}
			if tt.input.Phone != tt.expected.Phone {
				t.Errorf("Phone: got %s, want %s", tt.input.Phone, tt.expected.Phone)
			}
		})
	}
}
func TestUser_IsMinor(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		birthDay  time.Time
		wantMinor bool
	}{
		{
			name:      "30 years old - not minor",
			birthDay:  now.AddDate(-30, 0, 0),
			wantMinor: false,
		},
		{
			name:      "exactly 18 - not minor",
			birthDay:  now.AddDate(-18, 0, -1),
			wantMinor: false,
		},
		{
			name:      "17 years old - is minor",
			birthDay:  now.AddDate(-17, 0, 0),
			wantMinor: true,
		},
		{
			name:      "10 years old - is minor",
			birthDay:  now.AddDate(-10, 0, 0),
			wantMinor: true,
		},
		{
			name:      "birthday tomorrow (17 years) - is minor",
			birthDay:  now.AddDate(-18, 0, 1),
			wantMinor: true,
		},
		{
			name:      "birthday tomorrow (18 years) - not minor",
			birthDay:  now.AddDate(-19, 0, 1),
			wantMinor: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{BirthDay: tt.birthDay}
			if got := user.IsMinor(); got != tt.wantMinor {
				t.Errorf("IsMinor() = %v, want %v", got, tt.wantMinor)
			}
		})
	}
}
