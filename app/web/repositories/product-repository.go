package repositories

import (
	"context"

	"zen-test/app/helpers"
	"zen-test/app/web/models"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, db *gorm.DB, product models.Product) (models.Product, error)
	UpdateProduct(ctx context.Context, db *gorm.DB, product models.Product) (models.Product, error)
	DeleteProduct(ctx context.Context, db *gorm.DB, product models.Product) error
	GetProductById(ctx context.Context, db *gorm.DB, productId string) (models.Product, error)
	FindAllProducts(ctx context.Context, db *gorm.DB) ([]models.Product, error)
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, db *gorm.DB, product models.Product) (models.Product, error) {

	err := db.WithContext(ctx).Save(&product).Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r *ProductRepositoryImpl) UpdateProduct(ctx context.Context, db *gorm.DB, product models.Product) (models.Product, error) {

	err := db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", product.ID).Updates(&product).Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r *ProductRepositoryImpl) DeleteProduct(ctx context.Context, db *gorm.DB, product models.Product) error {
	err := db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", product.ID).Delete(&product).Error
	helpers.PanicIfError(err)
	return nil
}

func (r *ProductRepositoryImpl) GetProductById(ctx context.Context, db *gorm.DB, productId string) (models.Product, error) {
	var product models.Product
	err := db.WithContext(ctx).Model(&models.Product{}).
		Preload("Images").
		Where("id = ?", productId).
		Take(&product).
		Error
	helpers.PanicIfError(err)

	return product, nil
}

func (r *ProductRepositoryImpl) FindAllProducts(ctx context.Context, db *gorm.DB) ([]models.Product, error) {
	var products []models.Product

	err := db.WithContext(ctx).Model(&models.Product{}).
		Preload("Images").
		Find(&products).
		Error
	helpers.PanicIfError(err)

	return products, nil
}
