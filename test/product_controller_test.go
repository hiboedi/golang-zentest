package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zen-test/app/auth"
	"zen-test/app/web/models"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createProduct(product models.ProductCreateUpdate, db *gorm.DB) models.Product {
	productId := uuid.New().String()
	var images []models.Image
	for _, image := range images {
		image.ID = uuid.New().String()
		db.Save(&image)
		images = append(images, image)
	}
	productCreated := models.Product{
		ID:       productId,
		Name:     product.Name,
		Price:    product.Price,
		Stock:    product.Stock,
		Category: product.Category,
		Images:   images,
	}
	err := db.Save(&productCreated).Error
	if err != nil {
		return models.Product{}
	}
	return productCreated
}

func mockProduct(conditional string) models.ProductCreateUpdate {
	var product models.ProductCreateUpdate

	switch conditional {
	case "success":
		product = models.ProductCreateUpdate{
			Category: "laptop",          // require min 4 character
			Name:     "Huawei Matebook", // require min 4 character
			Price:    80000,             // require min 1
			Stock:    90,                // require min 1
			Images: []models.Image{
				{
					URL: "image-1",
				},
				{
					URL: "image-2",
				},
			},
		}

	case "failed": // trigger error validation for create or update product
		product = models.ProductCreateUpdate{
			Category: "lap",
			Name:     "Hua",
			Price:    0,
			Stock:    0,
			Images: []models.Image{
				{
					URL: "image-1",
				},
				{
					URL: "image-2",
				},
			},
		}

	case "update": // trigger error validation for create or update product
		product = models.ProductCreateUpdate{
			Category: "Gadget",      // edited
			Name:     "Samsung A54", // edited
			Price:    80000,
			Stock:    90,
		}
	default:
		return models.ProductCreateUpdate{}
	}
	return product
}

func truncateProduct(db *gorm.DB) {
	db.Exec("TRUNCATE products")
}

func TestCreateProductSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	requestBody := toRequestBody(mockProduct(success))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/products", requestBody)
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
	assert.Equal(t, "Huawei Matebook", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateProductFail(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	requestBody := toRequestBody(mockProduct(failed))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/products", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, statusBadRequest, responseBody["status"])
}

func TestUpdateProductSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	product := createProduct(mockProduct(success), db) // make sure add success as parameter
	productId := product.ID

	requestBody := toRequestBody(mockProduct(update))
	request := httptest.NewRequest(http.MethodPut, baseURL+"/products/"+productId, requestBody)
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
	assert.Equal(t, "Samsung A54", responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["category"])
}

func TestDeleteProductSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	product := createProduct(mockProduct(success), db) // make sure add success as parameter
	productId := product.ID

	requestBody := toRequestBody(mockProduct(update))
	request := httptest.NewRequest(http.MethodDelete, baseURL+"/products/"+productId, requestBody)
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

func TestFindAllProductSuccess(t *testing.T) {
	db := dbTest()
	router := routerTest(db)
	truncateUser(db)
	truncateProduct(db)

	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	createProduct(mockProduct(success), db)

	requestBody := toRequestBody(mockProduct(update))
	request := httptest.NewRequest(http.MethodGet, baseURL+"/products", requestBody)
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
