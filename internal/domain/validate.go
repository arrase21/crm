package domain

import (
	"errors"
	"strings"
	"time"
)

func (u *User) Validate() error {
	if u.Gender != "M" && u.Gender != "F" {
		return errors.New("Invalid option")
	}
	if u.BirthDay.IsZero() || u.BirthDay.After(time.Now()) {
		return errors.New("invalid birthday")
	}
	if u.IsMinor() {
		return errors.New("user must be over 18")
	}
	return nil
}

func (u *User) Normalize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Dni = strings.TrimSpace(u.Dni)
	u.Phone = strings.TrimSpace(u.Phone)
}

func (u *User) Required() error {
	u.Normalize()

	switch {
	case u.FirstName == "":
		return errors.New("first name is required")
	case u.LastName == "":
		return errors.New("last name is required")
	case u.Email == "":
		return errors.New("email is required")
	case u.Dni == "":
		return errors.New("dni is required")
	case u.Phone == "":
		return errors.New("phone is required")
	}

	return nil
}

func (u *User) IsMinor() bool {
	now := time.Now()
	age := now.Year() - u.BirthDay.Year()
	if now.Month() < u.BirthDay.Month() || (now.Month() == u.BirthDay.Month() && now.Day() < u.BirthDay.Day()) {
		age--
	}
	return age < 18
}

func (u *User) ValidateAll() error {
	if err := u.Required(); err != nil {
		return err
	}
	if err := u.Validate(); err != nil {
		return err
	}
	return nil
}

func (d *Department) Normalize() {
	d.Name = strings.TrimSpace(d.Name)
	d.Code = strings.TrimSpace(d.Code)
}
func (d *Department) Required() error {
	d.Normalize()

	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
func (d *Department) Validate() error {
	if len(d.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	if len(d.Code) > 20 {
		return errors.New("code must be at most 20 characters")
	}
	return nil
}
func (d *Department) ValidateAll() error {
	d.Normalize() // Normalizar UNA sola vez al inicio
	if err := d.Required(); err != nil {
		return err
	}
	if err := d.Validate(); err != nil {
		return err
	}
	return nil
}
