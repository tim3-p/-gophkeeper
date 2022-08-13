package common

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Opaque is the type fit for "sub-record" of the record
type Opaque interface {
	Pack() (string, error)
	Check() error
}

// Pack converts Account to string
func (a Account) Pack() (string, error) {
	opaque, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack converts Note to string
func (n Note) Pack() (string, error) {
	opaque, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack converts Account to string
func (c Card) Pack() (string, error) {
	opaque, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack returns Data field of Binary
func (b Binary) Pack() (string, error) {
	return b.Data, nil
}

// ErrDefaultFields is to indicate that some fields are left unset
var ErrDefaultFields = errors.New("mandatory field is set to default value")

// Check checks Account struct for correctness
func (a Account) Check() error {
	// let URL be optional, and username and password be obligatory
	if a.UserName == "" || a.Password == "" {
		return fmt.Errorf("UserName or Password: %w", ErrDefaultFields)
	}
	return nil
}

// Check checks Note struct for correctness
func (n Note) Check() error {
	// do not expose any requirements on the Note content
	return nil
}

// Check checks Card struct for correctness
func (c Card) Check() error {
	// let all fields be obligatory
	if c.Holder == "" || c.Number == "" ||
		c.ExpMonth == 0 || c.ExpYear == 0 || c.CVC == "" {
		return fmt.Errorf("Holder, Number, ExpMonth, ExpYear or CVC: %w",
			ErrDefaultFields)
	}
	return nil
}

// Check checks Binary struct for correctness
func (b Binary) Check() error {
	// do not expose any requirements on the Binary content
	return nil
}
