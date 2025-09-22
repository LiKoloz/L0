package models

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

type Delivery struct {
	ID       int    `json:"id" db:"id"`
	OrderUID string `json:"order_uid" db:"order_uid"`
	Name     string `json:"name" db:"name"`
	Phone    string `json:"phone" db:"phone"`
	Zip      string `json:"zip" db:"zip"`
	City     string `json:"city" db:"city"`
	Address  string `json:"address" db:"address"`
	Region   string `json:"region" db:"region"`
	Email    string `json:"email" db:"email"`
}

func (d *Delivery) Validate() error {
	if d.Name == "" {
		return errors.New("delivery name is required")
	}
	if utf8.RuneCountInString(d.Name) > 100 {
		return errors.New("delivery name must be less than 100 characters")
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(d.Phone) {
		return errors.New("invalid phone format")
	}

	if d.Zip == "" {
		return errors.New("zip code is required")
	}

	if d.City == "" {
		return errors.New("city is required")
	}

	if d.Address == "" {
		return errors.New("address is required")
	}

	return nil
}
