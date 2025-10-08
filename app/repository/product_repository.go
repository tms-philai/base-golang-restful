package repository

import (
	"base-golang-restful-app/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	FindByCategory(ctx context.Context, category string) ([]models.Product, error)
	FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.Product, int64, error)
	Search(ctx context.Context, query string) ([]models.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) error
}

type productRepository struct {
	BaseRepository
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		BaseRepository: NewBaseRepository(db),
		db:             db,
	}
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByCategory(ctx context.Context, category string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Where("category = ?", category).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindAllWithPagination(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Search(ctx context.Context, query string) ([]models.Product, error) {
	var products []models.Product
	searchPattern := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR description ILIKE ? OR category ILIKE ? OR sku ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, stock int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id.String()).
		Update("stock", stock).Error
}

func (r *productRepository) Create(ctx context.Context, entity interface{}) error {
	product, ok := entity.(*models.Product)
	if !ok {
		return gorm.ErrInvalidData
	}

	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).First(entity, "id = ?", id.String()).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID, entity interface{}) error {
	return r.db.WithContext(ctx).Delete(entity, "id = ?", id.String()).Error
}
