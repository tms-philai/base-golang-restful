package auth

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong     = errors.New("password must not exceed 72 characters")
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial   = errors.New("password must contain at least one special character")
	ErrPasswordMismatch    = errors.New("passwords do not match")
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
	DefaultBcryptCost = 12
)

type PasswordValidator struct {
	MinLength        int
	MaxLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumber    bool
	RequireSpecial   bool
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        MinPasswordLength,
		MaxLength:        MaxPasswordLength,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   true,
	}
}

func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.MinLength {
		return ErrPasswordTooShort
	}

	if len(password) > v.MaxLength {
		return ErrPasswordTooLong
	}

	if v.RequireUppercase && !hasUppercase(password) {
		return ErrPasswordNoUppercase
	}

	if v.RequireLowercase && !hasLowercase(password) {
		return ErrPasswordNoLowercase
	}

	if v.RequireNumber && !hasNumber(password) {
		return ErrPasswordNoNumber
	}

	if v.RequireSpecial && !hasSpecialChar(password) {
		return ErrPasswordNoSpecial
	}

	return nil
}

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`)
	return specialChars.MatchString(s)
}

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: DefaultBcryptCost,
	}
}

func NewPasswordHasherWithCost(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = DefaultBcryptCost
	}
	return &PasswordHasher{
		cost: cost,
	}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	if len(password) > MaxPasswordLength {
		return "", ErrPasswordTooLong
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (h *PasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (h *PasswordHasher) IsHashValid(hashedPassword string) bool {
	cost, err := bcrypt.Cost([]byte(hashedPassword))
	return err == nil && cost >= bcrypt.MinCost && cost <= bcrypt.MaxCost
}

func HashPassword(password string) (string, error) {
	hasher := NewPasswordHasher()
	return hasher.Hash(password)
}

func ComparePassword(hashedPassword, password string) error {
	hasher := NewPasswordHasher()
	return hasher.Compare(hashedPassword, password)
}

func ValidatePassword(password string) error {
	validator := NewPasswordValidator()
	return validator.Validate(password)
}
