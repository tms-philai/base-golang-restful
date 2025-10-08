package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordValidator_Validate(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name          string
		password      string
		expectError   bool
		expectedError error
	}{
		{
			name:        "valid password",
			password:    "ValidPass123!",
			expectError: false,
		},
		{
			name:          "too short",
			password:      "Short1!",
			expectError:   true,
			expectedError: ErrPasswordTooShort,
		},
		{
			name:          "too long",
			password:      "VeryLongPassword123!" + string(make([]byte, 100)),
			expectError:   true,
			expectedError: ErrPasswordTooLong,
		},
		{
			name:          "no uppercase",
			password:      "lowercase123!",
			expectError:   true,
			expectedError: ErrPasswordNoUppercase,
		},
		{
			name:          "no lowercase",
			password:      "UPPERCASE123!",
			expectError:   true,
			expectedError: ErrPasswordNoLowercase,
		},
		{
			name:          "no number",
			password:      "NoNumber!",
			expectError:   true,
			expectedError: ErrPasswordNoNumber,
		},
		{
			name:          "no special character",
			password:      "NoSpecial123",
			expectError:   true,
			expectedError: ErrPasswordNoSpecial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordValidator_CustomRules(t *testing.T) {
	validator := &PasswordValidator{
		MinLength:        6,
		MaxLength:        20,
		RequireUppercase: false,
		RequireLowercase: true,
		RequireNumber:    false,
		RequireSpecial:   false,
	}

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "valid with custom rules",
			password:    "simple",
			expectError: false,
		},
		{
			name:        "too short for custom rules",
			password:    "short",
			expectError: true,
		},
		{
			name:        "no lowercase required",
			password:    "NOLOWERCASE",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "TestPassword123!"
	hash, err := hasher.Hash(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestPasswordHasher_Hash_TooLong(t *testing.T) {
	hasher := NewPasswordHasher()

	password := string(make([]byte, MaxPasswordLength+1))
	hash, err := hasher.Hash(password)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPasswordTooLong)
	assert.Empty(t, hash)
}

func TestPasswordHasher_Compare(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "TestPassword123!"
	hash, err := hasher.Hash(password)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		testPassword  string
		expectError   bool
	}{
		{
			name:         "correct password",
			testPassword: password,
			expectError:  false,
		},
		{
			name:         "incorrect password",
			testPassword: "WrongPassword123!",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Compare(hash, tt.testPassword)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordHasher_IsHashValid(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "TestPassword123!"
	hash, err := hasher.Hash(password)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		expected bool
	}{
		{
			name:     "valid hash",
			hash:     hash,
			expected: true,
		},
		{
			name:     "invalid hash",
			hash:     "not-a-valid-hash",
			expected: false,
		},
		{
			name:     "empty hash",
			hash:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasher.IsHashValid(tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPasswordHasher_WithCustomCost(t *testing.T) {
	tests := []struct {
		name         string
		cost         int
		expectedCost int
	}{
		{
			name:         "valid cost",
			cost:         10,
			expectedCost: 10,
		},
		{
			name:         "cost too low",
			cost:         2,
			expectedCost: DefaultBcryptCost,
		},
		{
			name:         "cost too high",
			cost:         50,
			expectedCost: DefaultBcryptCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewPasswordHasherWithCost(tt.cost)
			assert.Equal(t, tt.expectedCost, hasher.cost)
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "TestPassword123!"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestComparePassword(t *testing.T) {
	password := "TestPassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hash, password)
	assert.NoError(t, err)

	err = ComparePassword(hash, "WrongPassword")
	assert.Error(t, err)
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "valid password",
			password:    "ValidPass123!",
			expectError: false,
		},
		{
			name:        "invalid password",
			password:    "weak",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHasUppercase(t *testing.T) {
	assert.True(t, hasUppercase("Hello"))
	assert.False(t, hasUppercase("hello"))
}

func TestHasLowercase(t *testing.T) {
	assert.True(t, hasLowercase("Hello"))
	assert.False(t, hasLowercase("HELLO"))
}

func TestHasNumber(t *testing.T) {
	assert.True(t, hasNumber("Hello123"))
	assert.False(t, hasNumber("Hello"))
}

func TestHasSpecialChar(t *testing.T) {
	assert.True(t, hasSpecialChar("Hello!"))
	assert.True(t, hasSpecialChar("Hello@World"))
	assert.False(t, hasSpecialChar("Hello123"))
}
