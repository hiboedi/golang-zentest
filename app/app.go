package app

import (
	"zen-test/app/database"
	"zen-test/app/web/controllers"
	"zen-test/app/web/repositories"
	"zen-test/app/web/router"
	"zen-test/app/web/services"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AppConfig struct {
	AppName        string
	AppEnv         string
	AppPort        string
	AppHost        string
	AppBasePath    string
	AppVersion     string
	AppDescription string
	AppSchema      []string
}

func Init() (*mux.Router, AppConfig) {

	appConfig := AppConfig{
		AppName:        "Zenstore",
		AppEnv:         "dev",
		AppPort:        "8000",
		AppHost:        "localhost:8000",
		AppBasePath:    "/",
		AppVersion:     "1.0",
		AppDescription: "Sample of Zenstore API Documentation.",
		AppSchema:      []string{"http", "https"},
	}

	db := database.InitializeDB()
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

	database.DBMigrate(db)
	return router, appConfig
}
