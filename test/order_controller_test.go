package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zen-test/app/auth"
	"zen-test/app/consts"
	"zen-test/app/helpers"
	"zen-test/app/web/models"
	"zen-test/app/web/services"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createOrder(order models.OrderItemCreateUpdate, user models.User, product models.Product, db *gorm.DB) models.Order {

	taxAmount := services.CountTax(product.Price, order.Quantity, consts.TaxRate)
	totalPrice := (product.Price * float64(order.Quantity)) - taxAmount
	orderId := uuid.New().String()
	orderCreated := models.Order{
		ID:           orderId,
		UserID:       user.ID,
		IsPaid:       false,
		Status:       consts.OrderPaymentStatusUnpaid,
		CustomerName: user.Name,
		Phone:        user.Phone,
		TotalPrice:   totalPrice,
		Address:      user.Address,
	}

	err := db.Model(&models.Order{}).Create(&orderCreated).Error
	helpers.PanicIfError(err)

	return orderCreated
}

func mockOrder(conditional string, productId string) models.OrderItemCreateUpdate {
	var orderItem models.OrderItemCreateUpdate

	switch conditional {
	case "success":
		orderItem = models.OrderItemCreateUpdate{
			Quantity:  8,
			ProductID: productId,
		}

	case "failed": // trigger error validation for create or update orderItem
		orderItem = models.OrderItemCreateUpdate{
			Quantity:  1,              // min 1
			ProductID: "just-example", // wrong product id
		}

	case "update": // trigger error validation for create or update orderItem
		orderItem = models.OrderItemCreateUpdate{
			Quantity:  10,
			ProductID: productId,
		}
	default:
		return models.OrderItemCreateUpdate{}
	}
	return orderItem
}

func truncateOrder(db *gorm.DB) {
	db.Exec("TRUNCATE orders")
}

func TestCreateOrderSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)
	truncateOrder(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)
	product := createProduct(mockProduct(success), db)

	requestBody := toRequestBody(mockOrder(success, product.ID))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/orders", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, statusOk, responseBody["status"])
	assert.Equal(t, "Bandung", responseBody["data"].(map[string]interface{})["address"])
	assert.Equal(t, "08811212112", responseBody["data"].(map[string]interface{})["phone"])
}

func TestFindAllOrdersSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)
	truncateOrder(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	product := createProduct(mockProduct(success), db)
	createOrder(mockOrder(success, product.ID), user, product, db)

	requestBody := toRequestBody(mockProduct(update))
	request := httptest.NewRequest(http.MethodGet, baseURL+"/orders", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, statusOk, responseBody["status"])
}
