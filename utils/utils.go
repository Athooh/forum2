package utils

import (
	"errors"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain text password.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateSessionToken creates a unique session token using UUID.
func GenerateSessionToken() (string, error) {
	token := uuid.New().String()
	return token, nil
}

// ValidatePasswordStrength checks if a password is strong enough.
func ValidatePasswordStrength(password string) error {
	var hasMinLength, hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) >= 8 {
		hasMinLength = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLength {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper {
		return errors.New("password must include at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must include at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must include at least one numeric digit")
	}
	if !hasSpecial {
		return errors.New("password must include at least one special character")
	}

	return nil
}