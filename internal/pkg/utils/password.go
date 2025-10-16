package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) []string {
	var errors []string

	if len(password) < 6 {
		errors = append(errors, "password must be at least 6 characters long")
	}

	if len(password) > 128 {
		errors = append(errors, "password must be less than 128 characters long")
	}

	// Add more password strength validations as needed
	// hasUpper := false
	// hasLower := false
	// hasDigit := false
	// hasSpecial := false

	// for _, char := range password {
	// 	switch {
	// 	case unicode.IsUpper(char):
	// 		hasUpper = true
	// 	case unicode.IsLower(char):
	// 		hasLower = true
	// 	case unicode.IsDigit(char):
	// 		hasDigit = true
	// 	case unicode.IsPunct(char) || unicode.IsSymbol(char):
	// 		hasSpecial = true
	// 	}
	// }

	// if !hasUpper {
	// 	errors = append(errors, "password must contain at least one uppercase letter")
	// }
	// if !hasLower {
	// 	errors = append(errors, "password must contain at least one lowercase letter")
	// }
	// if !hasDigit {
	// 	errors = append(errors, "password must contain at least one digit")
	// }
	// if !hasSpecial {
	// 	errors = append(errors, "password must contain at least one special character")
	// }

	return errors
}
