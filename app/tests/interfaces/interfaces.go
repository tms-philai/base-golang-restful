package interfaces

import "base-golang-restful-app/models"

// UserServiceInterface defines the interface for user service
type UserServiceInterface interface {
	Create(req models.UserCreateRequest) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id string, req models.UserUpdateRequest) (*models.User, error)
	Delete(id string) error
	List(page, limit int) ([]models.User, models.PaginationMetadata, error)
	ChangePassword(userID, currentPassword, newPassword string) error
	ValidateCredentials(username, password string) (*models.User, error)
}

// JWTServiceInterface defines the interface for JWT service
type JWTServiceInterface interface {
	GenerateTokenPair(user *models.User) (*models.TokenPair, error)
	ValidateToken(tokenString string) (*models.JWTClaims, error)
	RefreshToken(refreshTokenString string, user *models.User) (*models.TokenPair, error)
	ExtractTokenFromHeader(authHeader string) (string, error)
}

// ProductServiceInterface defines the interface for product service
type ProductServiceInterface interface {
	Create(req models.ProductCreateRequest, createdBy string) (*models.Product, error)
	GetByID(id string) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(id string, req models.ProductUpdateRequest) (*models.Product, error)
	Delete(id string) error
	List(query models.ProductListQuery) ([]models.Product, models.PaginationMetadata, error)
	GetCategories() []string
	UpdateStock(id string, stock int) error
	CheckStock(id string, quantity int) (bool, error)
}
