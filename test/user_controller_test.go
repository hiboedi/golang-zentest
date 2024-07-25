package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"zen-test/app/auth"
	"zen-test/app/helpers"
	"zen-test/app/web/models"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createUser(user models.UserCreate, db *gorm.DB) models.User {
	hashPassword, _ := helpers.MakePassword(user.Password)
	userCreated := models.User{
		ID:       uuid.New().String(),
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Address:  user.Address,
		Password: hashPassword,
	}

	err := db.Model(&models.User{}).Create(&userCreated).Error
	helpers.PanicIfError(err)

	return userCreated
}

func mockUser(conditional string) models.UserCreate {
	var user models.UserCreate

	switch conditional {
	case "success":
		user = models.UserCreate{
			Name:     "Budiman",           // min 4 character
			Email:    "example@gmail.com", // should be an email
			Phone:    "08811212112",       // min 8 numbers
			Password: "example",           // min 6 characters
			Address:  "Bandung",
		}

	case "failed": // trigger error validation for create or update user
		user = models.UserCreate{
			Name:     "Bu",
			Email:    "example.com",
			Phone:    "09",
			Password: "ex",
			Address:  "Bandung",
		}

	case "update": // trigger error validation for create or update user
		user = models.UserCreate{
			Phone:    "08811212112",
			Password: "example",
			Address:  "Jakarta", // edited
		}
	default:
		return models.UserCreate{}
	}
	return user
}

func login(conditional string) models.UserLogin {
	var user models.UserLogin

	switch conditional {
	case "success":
		user = models.UserLogin{
			Email:    "example@gmail.com",
			Password: "example",
		}

	case "failed": // trigger error email or password not valid
		user = models.UserLogin{
			Email:    "budi@gmail.com",
			Password: "blablabla",
		}
	default:
		return user
	}
	return user
}

func truncateUser(db *gorm.DB) {
	db.Exec("TRUNCATE users")
}

func TestUserRegister(t *testing.T) {
	db := dbTest()
	truncateUser(db)
	router := routerTest(db)

	requestBody := toRequestBody(mockUser(success))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/users/signup", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, statusOk, responseBody["status"])
	assert.Equal(t, "Budiman", responseBody["data"].(map[string]interface{})["name"])
}

func TestUserRegisterValidationError(t *testing.T) {
	db := dbTest()
	truncateUser(db)
	router := routerTest(db)

	requestBody := toRequestBody(mockUser(failed))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/users/signup", requestBody)
	request.Header.Add("Content-Type", "application/json")

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

func TestLoginSuccess(t *testing.T) {
	db := dbTest()
	truncateUser(db)
	router := routerTest(db)
	user := mockUser(success) // make sure add success as parameter
	createUser(user, db)

	requestBody := toRequestBody(login(success))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/users/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, statusOk, responseBody["status"])
	assert.Equal(t, "Budiman", responseBody["data"].(map[string]interface{})["name"])
}

func TestLoginFail(t *testing.T) {
	db := dbTest()
	truncateUser(db)
	router := routerTest(db)
	user := mockUser(success) // make sure add success as parameter
	createUser(user, db)

	requestBody := toRequestBody(login(failed))
	request := httptest.NewRequest(http.MethodPost, baseURL+"/users/login", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 500, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 500, int(responseBody["code"].(float64)))
	assert.Equal(t, statusInternalServerError, responseBody["status"])
}

func TestUpdateSuccess(t *testing.T) {
	db := dbTest()
	truncateUser(db)
	router := routerTest(db)
	data := mockUser(success) // make sure add success as parameter
	user := createUser(data, db)
	userId := user.ID
	token, _ := auth.CreateToken(userId)

	requestBody := toRequestBody(mockUser(update))
	request := httptest.NewRequest(http.MethodPut, baseURL+"/users/"+userId, requestBody)
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
	assert.Equal(t, "Jakarta", responseBody["data"].(map[string]interface{})["address"])
}
