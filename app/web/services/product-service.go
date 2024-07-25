package services

import (
	"context"
	"time"

	"zen-test/app/exceptions"
	"zen-test/app/helpers"
	"zen-test/app/web/models"
	"zen-test/app/web/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	ProductRepository repositories.ProductRepository
	ImageRepository   repositories.ImageRepositoy
	DB                *gorm.DB
	Validate          *validator.Validate
}

type ProductService interface {
	Create(ctx context.Context, request models.ProductCreateUpdate) models.ProductResponse
	Update(ctx context.Context, request models.ProductCreateUpdate, productId string) models.ProductResponse
	Delete(ctx context.Context, productId string)
	FindById(ctx context.Context, productId string) models.ProductResponse
	FindAll(ctx context.Context) []models.ProductResponse
}

func NewProductService(productRepo repositories.ProductRepository, imageRepo repositories.ImageRepositoy, db *gorm.DB, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepo,
		ImageRepository:   imageRepo,
		DB:                db,
		Validate:          validate,
	}
}

func (s *ProductServiceImpl) Create(ctx context.Context, request models.ProductCreateUpdate) models.ProductResponse {
	err := s.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	productId := uuid.New().String()

	product := models.Product{
		ID:       productId,
		Name:     request.Name,
		Price:    request.Price,
		Stock:    request.Stock,
		Category: request.Category,
	}

	_, err = s.ProductRepository.CreateProduct(ctx, tx, product)
	helpers.PanicIfError(err)

	var images []models.Image
	for _, image := range request.Images {
		image.ID = uuid.New().String()
		image.ProductID = productId
		createdimage, err := s.ImageRepository.CreateImage(ctx, tx, image)
		helpers.PanicIfError(err)
		images = append(images, createdimage)
	}

	product.Images = images
	data, err := s.ProductRepository.CreateProduct(ctx, tx, product)
	helpers.PanicIfError(err)

	return models.ToProductResponse(data)
}

func (s *ProductServiceImpl) Update(ctx context.Context, request models.ProductCreateUpdate, productId string) models.ProductResponse {
	err := s.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	product, err := s.ProductRepository.GetProductById(ctx, tx, productId)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	var updatedImages []models.Image

	for _, updateImage := range request.Images {
		for _, existImage := range product.Images {
			updateImage.ID = existImage.ID
			updateImage.ProductID = existImage.ProductID
			updateImage.CreatedAt = existImage.CreatedAt
			updateImage.UpdatedAt = time.Now()
			_, err := s.ImageRepository.UpdateImage(ctx, tx, existImage)
			helpers.PanicIfError(err)
		}
		updatedImages = append(updatedImages, updateImage)
	}

	product.Name = request.Name
	product.Category = request.Category
	product.Price = request.Price
	product.Stock = request.Stock
	product.Images = updatedImages

	data, err := s.ProductRepository.UpdateProduct(ctx, tx, product)
	helpers.PanicIfError(err)

	return models.ToProductResponse(data)
}

func (s *ProductServiceImpl) Delete(ctx context.Context, productId string) {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	product, err := s.ProductRepository.GetProductById(ctx, tx, productId)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	images, err := s.ImageRepository.FindImages(ctx, tx, product.ID)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	for _, image := range images {
		err := s.ImageRepository.DeleteImage(ctx, tx, image)
		helpers.PanicIfError(err)
	}

	err = s.ProductRepository.DeleteProduct(ctx, tx, product)
	helpers.PanicIfError(err)
}

func (s *ProductServiceImpl) FindAll(ctx context.Context) []models.ProductResponse {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	products, err := s.ProductRepository.FindAllProducts(ctx, tx)
	helpers.PanicIfError(err)
	return models.ToProductResponses(products)
}

func (s *ProductServiceImpl) FindById(ctx context.Context, productId string) models.ProductResponse {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	product, err := s.ProductRepository.GetProductById(ctx, tx, productId)
	helpers.PanicIfError(err)
	return models.ToProductResponse(product)
}
