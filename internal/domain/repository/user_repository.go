package repository

import (
	"base-gin/internal/domain/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByIDWithRoles(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	Search(ctx context.Context, query string) ([]models.User, error)
}

type userRepository struct {
	BaseRepository
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Search(ctx context.Context, query string) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, entity interface{}) error {
	user, ok := entity.(*models.User)
	if !ok {
		return gorm.ErrInvalidData
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).First(entity, "id = ?", id.String()).Error
}

// FindByIDWithRoles loads user with roles and their permissions
func (r *userRepository) FindByIDWithRoles(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&user, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).Delete(entity, "id = ?", id.String()).Error
}
