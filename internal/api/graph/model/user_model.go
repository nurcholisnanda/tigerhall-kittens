package model

import (
	"errors"
	"regexp"
)

type User struct {
	ID       string `json:"id" gorm:"ype:varchar(255);primarykey"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);not null;unique;index"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
	Salt     string `gorm:"type:varchar(24)"`
}

// Validate checks if the input data is valid.
func (u *NewUser) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if len(u.Name) < 2 || len(u.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	// You can add a regular expression here for email validation.
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
