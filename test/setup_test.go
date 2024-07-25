package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"zen-test/app/database"
	"zen-test/app/helpers"
	"zen-test/app/middleware"
	"zen-test/app/web/controllers"
	"zen-test/app/web/repositories"
	"zen-test/app/web/router"
	"zen-test/app/web/services"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	baseURL                   string = "http://localhost:8000"
	success                   string = "success"
	failed                    string = "failed"
	update                    string = "update"
	statusOk                  string = "Ok"
	statusBadRequest          string = "Bad Request"
	statusInternalServerError string = "Internal Server Error"
)

func toRequestBody(any interface{}) io.Reader {
	resultJson, err := json.Marshal(any)
	if err != nil {
		log.Fatalf("Failed to marshal : %v", err)
	}

	return bytes.NewReader(resultJson)
}

func dbTest() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var dbConfig database.DBConfig

	dbConfig.DBName = helpers.GetEnv("DATABASE_NAME_TEST", "exam")
	dbConfig.DBPort = helpers.GetEnv("DATABASE_PORT_TEST", "5432")
	dbConfig.DBUser = helpers.GetEnv("DATABASE_USER_TEST", "boedi")
	dbConfig.DBPassword = helpers.GetEnv("DATABASE_PASSWORD_TEST", "")
	dbConfig.DBHost = helpers.GetEnv("DATABASE_HOST_TEST", "localhost")
	dbConfig.DBDriver = helpers.GetEnv("DATABASE_DRIVER_TEST", "postgres")

	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		helpers.PanicIfError(err)
		database.DBMigrate(db)
		return db
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBPort, dbConfig.DBName)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

		if err != nil {
			panic("Failed on connecting to the database server")
		}
		database.DBMigrate(db)
		return db
	}
}

func routerTest(db *gorm.DB) http.Handler {
	validate := validator.New()

	userRepo := repositories.NewUserRepository()
	productRepo := repositories.NewProductRepository()
	orderRepo := repositories.NewOrderRepository()
	imageRepo := repositories.NewImageRepository()

	userService := services.NewUserService(userRepo, db, validate)
	productservice := services.NewProductService(productRepo, imageRepo, db, validate)
	orderService := services.NewOrderService(orderRepo, productRepo, userRepo, db, validate)

	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productservice)
	orderController := controllers.NewOrderController(orderService)

	go orderService.AutoCancelUnpaidOrders()

	router := router.InitializeRouter(userController, productController, orderController)

	return middleware.AuthMiddleware(router)
}
