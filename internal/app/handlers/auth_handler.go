package handlers

import (
	"base-gin/internal/app/middleware"
	"base-gin/internal/domain/models"
	"base-gin/internal/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserServiceInterface defines the interface for UserService
type UserServiceInterface interface {
	Create(req models.UserCreateRequest) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByIDWithRoles(id string) (*models.User, error)
	ValidateCredentials(email, password string) (*models.User, error)
	ChangePassword(userID, oldPassword, newPassword string) error
}

// JWTManagerInterface defines the interface for JWTManager
type JWTManagerInterface interface {
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(token string, expectedType auth.TokenType) (*auth.JWTClaims, error)
	GetTokenDuration(tokenType auth.TokenType) time.Duration
}

// RoleServiceInterface defines the interface for RoleService
type RoleServiceInterface interface {
	AssignDefaultRole(userID uuid.UUID) error
}

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService UserServiceInterface
	roleService RoleServiceInterface
	jwtManager  JWTManagerInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService UserServiceInterface, roleService RoleServiceInterface, jwtManager JWTManagerInterface) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		roleService: roleService,
		jwtManager:  jwtManager,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	// Create user
	userReq := models.UserCreateRequest{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.userService.Create(userReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	// Assign default "user" role to new user
	if h.roleService != nil {
		if err := h.roleService.AssignDefaultRole(user.ID); err != nil {
			// Log error but don't fail registration
			// The user is created but without a role - admin can assign later
			// In production, use proper logging: log.Printf("Failed to assign default role: %v", err)
		}
	}

	// Get user with roles for response
	userWithRoles, err := h.userService.GetByIDWithRoles(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "user_fetch_failed",
			Message: "Failed to fetch user details",
		})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	response := models.AuthResponse{
		User:         userWithRoles.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetTokenDuration(auth.AccessToken).Seconds()),
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	// Validate credentials
	user, err := h.userService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "authentication_failed",
			Message: err.Error(),
		})
		return
	}

	// Get user with roles for response
	userWithRoles, err := h.userService.GetByIDWithRoles(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "user_fetch_failed",
			Message: "Failed to fetch user details",
		})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	response := models.AuthResponse{
		User:         userWithRoles.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetTokenDuration(auth.AccessToken).Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	// Validate refresh token and get claims
	claims, err := h.jwtManager.ValidateToken(req.RefreshToken, auth.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// Get user with roles
	user, err := h.userService.GetByIDWithRoles(claims.UserID.String())
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "user_not_found",
			Message: "User associated with token not found",
		})
		return
	}

	// Generate new token pair
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate refresh token",
		})
		return
	}

	response := models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetTokenDuration(auth.AccessToken).Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userInterface, exists := middleware.GetCurrentUser(c)
	if !exists || userInterface == nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}
	c.JSON(http.StatusOK, user.ToResponse())
}

// ChangePassword handles password change
// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "password_change_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Password changed successfully",
	})
}
