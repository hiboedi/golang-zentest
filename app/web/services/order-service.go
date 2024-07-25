package services

import (
	"context"
	"log"
	"time"

	"zen-test/app/consts"
	"zen-test/app/exceptions"
	"zen-test/app/helpers"
	"zen-test/app/web/models"
	"zen-test/app/web/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	FindAllOrder(ctx context.Context) ([]models.OrderResponse, error)
	CreateOrder(ctx context.Context, request models.OrderItemCreateUpdate, userId string) (models.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status string) error
	AutoCancelUnpaidOrders()
	CancelUnpaidOrders(ctx context.Context)
}

type OrderRepositoryImpl struct {
	OrderRepository   repositories.OrderRepository
	ProductRepository repositories.ProductRepository
	UserRepository    repositories.UserRepository
	DB                *gorm.DB
	Validate          *validator.Validate
}

func NewOrderService(orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository, userRepo repositories.UserRepository, db *gorm.DB, validate *validator.Validate) OrderService {
	return &OrderRepositoryImpl{
		OrderRepository:   orderRepo,
		DB:                db,
		ProductRepository: productRepo,
		UserRepository:    userRepo,
		Validate:          validate,
	}
}

func (s *OrderRepositoryImpl) FindAllOrder(ctx context.Context) ([]models.OrderResponse, error) {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	data, err := s.OrderRepository.FindAllOrder(ctx, tx)
	helpers.PanicIfError(err)

	return models.ToOrderResponses(data), nil
}

func (s *OrderRepositoryImpl) CreateOrder(ctx context.Context, request models.OrderItemCreateUpdate, userId string) (models.OrderResponse, error) {
	err := s.Validate.Struct(request)
	helpers.PanicIfError(err)

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	user, err := s.UserRepository.GetUserById(ctx, tx, userId)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	product, err := s.ProductRepository.GetProductById(ctx, tx, request.ProductID)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	if request.Quantity > product.Stock {
		panic("Product is out of stock")
	}

	taxAmount := CountTax(product.Price, request.Quantity, consts.TaxRate)
	totalPrice := (product.Price * float64(request.Quantity)) - taxAmount

	order := models.Order{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		IsPaid:       false,
		Status:       consts.OrderPaymentStatusUnpaid,
		CustomerName: user.Name,
		Phone:        user.Phone,
		TotalPrice:   totalPrice,
		Address:      user.Address,
	}
	orderCreated, err := s.OrderRepository.CreateOrder(ctx, tx, order)
	helpers.PanicIfError(err)

	product.Stock = product.Stock - uint32(request.Quantity)

	_, err = s.ProductRepository.UpdateProduct(ctx, tx, product)
	helpers.PanicIfError(err)

	orderItem := models.OrderItem{
		ID:        uuid.New().String(),
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  request.Quantity,
	}
	_, err = s.OrderRepository.CreateOrderItem(ctx, tx, orderItem)
	helpers.PanicIfError(err)

	return models.ToOrderResponse(orderCreated), nil
}

func (s *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderId string, status string) error {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	order, err := s.OrderRepository.FindOrder(ctx, tx, orderId)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	order.Status = status
	if status == consts.OrderPaymentStatusPaid {
		order.IsPaid = true
		order.Status = consts.OrderPaymentStatusPaid
	} else if status == consts.OrderPaymentStatusCancel {
		order.IsPaid = false
		order.Status = consts.OrderPaymentStatusCancel
	}

	_, err = s.OrderRepository.UpdateOrder(ctx, tx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderRepositoryImpl) AutoCancelUnpaidOrders() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		ctx := context.Background()
		s.CancelUnpaidOrders(ctx)
	}
}

func (s *OrderRepositoryImpl) CancelUnpaidOrders(ctx context.Context) {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	orders, err := s.OrderRepository.GetUnpaidOrdersOlderThan(ctx, tx, 1*time.Hour)
	if err != nil {
		log.Printf("Error fetching unpaid orders: %v", err)
		return
	}

	for _, order := range orders {
		s.UpdateOrderStatus(ctx, order.ID, consts.OrderPaymentStatusCancel)
	}
}

func CountTax(price float64, qty uint32, taxRate float64) float64 {
	return price * float64(qty) * taxRate
}
