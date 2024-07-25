package repositories

import (
	"context"
	"time"

	"zen-test/app/consts"
	"zen-test/app/helpers"
	"zen-test/app/web/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, db *gorm.DB, order models.Order) (models.Order, error)
	UpdateOrder(ctx context.Context, db *gorm.DB, order models.Order) (models.Order, error)
	CreateOrderItem(ctx context.Context, db *gorm.DB, orderItem models.OrderItem) (models.OrderItem, error)
	FindAllOrder(ctx context.Context, db *gorm.DB) ([]models.Order, error)
	GetUnpaidOrdersOlderThan(ctx context.Context, tx *gorm.DB, duration time.Duration) ([]models.Order, error)
	FindOrder(ctx context.Context, db *gorm.DB, orderId string) (models.Order, error)
}

type orderRepositoryImpl struct {
}

func NewOrderRepository() OrderRepository {
	return &orderRepositoryImpl{}
}

func (r *orderRepositoryImpl) CreateOrder(ctx context.Context, db *gorm.DB, order models.Order) (models.Order, error) {

	err := db.WithContext(ctx).
		Model(&models.Order{}).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Images").
		Create(&order).Error
	helpers.PanicIfError(err)

	return order, nil
}

func (r *orderRepositoryImpl) UpdateOrder(ctx context.Context, db *gorm.DB, order models.Order) (models.Order, error) {

	err := db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", order.ID).Updates(&order).Error
	helpers.PanicIfError(err)

	return order, nil
}

func (r *orderRepositoryImpl) FindOrder(ctx context.Context, db *gorm.DB, orderId string) (models.Order, error) {
	var order models.Order

	err := db.WithContext(ctx).
		Model(&models.Order{}).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Images").
		Where("id = ?", orderId).
		Take(&order).Error
	helpers.PanicIfError(err)

	return order, nil
}

func (r *orderRepositoryImpl) CreateOrderItem(ctx context.Context, db *gorm.DB, orderItem models.OrderItem) (models.OrderItem, error) {

	err := db.WithContext(ctx).Create(&orderItem).Error
	helpers.PanicIfError(err)

	return orderItem, nil
}

func (r *orderRepositoryImpl) GetUnpaidOrdersOlderThan(ctx context.Context, tx *gorm.DB, duration time.Duration) ([]models.Order, error) {
	var orders []models.Order
	cutoff := time.Now().Add(-duration)
	err := tx.Where("status = ? AND created_at < ?", consts.OrderPaymentStatusUnpaid, cutoff).Find(&orders).Error
	return orders, err
}

func (r *orderRepositoryImpl) FindAllOrder(ctx context.Context, db *gorm.DB) ([]models.Order, error) {
	var Orders []models.Order

	err := db.WithContext(ctx).
		Model(&models.Order{}).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Images").
		Find(&Orders).Error
	helpers.PanicIfError(err)

	return Orders, nil
}
