package repositories

import (
	"context"

	"zen-test/app/helpers"
	"zen-test/app/web/models"

	"gorm.io/gorm"
)

type ImageRepositoy interface {
	CreateImage(ctx context.Context, db *gorm.DB, image models.Image) (models.Image, error)
	UpdateImage(ctx context.Context, db *gorm.DB, image models.Image) (models.Image, error)
	DeleteImage(ctx context.Context, db *gorm.DB, image models.Image) error
	FindImages(ctx context.Context, db *gorm.DB, productId string) ([]models.Image, error)
	GetImageById(ctx context.Context, db *gorm.DB, imageId string) (models.Image, error)
}

type ImageRepositoryImpl struct {
}

func NewImageRepository() ImageRepositoy {
	return &ImageRepositoryImpl{}
}

func (r *ImageRepositoryImpl) CreateImage(ctx context.Context, db *gorm.DB, image models.Image) (models.Image, error) {

	err := db.WithContext(ctx).Create(&image).Error
	helpers.PanicIfError(err)

	return image, nil
}

func (r *ImageRepositoryImpl) UpdateImage(ctx context.Context, db *gorm.DB, image models.Image) (models.Image, error) {

	err := db.WithContext(ctx).Model(&models.Image{}).Where("id = ?", image.ID).Updates(&image).Error
	helpers.PanicIfError(err)

	return image, nil
}

func (r *ImageRepositoryImpl) DeleteImage(ctx context.Context, db *gorm.DB, image models.Image) error {
	err := db.WithContext(ctx).Where("id = ?", image.ID).Delete(&image).Error
	helpers.PanicIfError(err)
	return nil
}

func (r *ImageRepositoryImpl) FindImages(ctx context.Context, db *gorm.DB, productId string) ([]models.Image, error) {
	var images []models.Image
	err := db.WithContext(ctx).Model(&models.Image{}).Where("product_id = ?", productId).Find(&images).Error
	helpers.PanicIfError(err)

	return images, nil
}

func (r *ImageRepositoryImpl) GetImageById(ctx context.Context, db *gorm.DB, imageId string) (models.Image, error) {
	var image models.Image
	err := db.WithContext(ctx).Model(&models.Image{}).
		Preload("Images").
		Where("id = ?", imageId).
		Take(&image).
		Error
	helpers.PanicIfError(err)

	return image, nil
}
