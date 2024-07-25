package router

import (
	"zen-test/app/middleware"
	"zen-test/app/web/controllers"

	"github.com/gorilla/mux"
)

func InitializeRouter(
	userController controllers.UserController,
	productController controllers.ProductController,
	orderController controllers.OrderController,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users/login", userController.Login).Methods("POST")
	router.HandleFunc("/users/signup", userController.SignUp).Methods("POST")
	router.HandleFunc("/users/{userId}", userController.Update).Methods("PUT")
	router.HandleFunc("/users/logout", userController.Logout).Methods("POST")
	router.HandleFunc("/users/refresh-token", userController.RefreshToken).Methods("POST")

	router.HandleFunc("/products", productController.Create).Methods("POST")
	router.HandleFunc("/products", productController.FindAll).Methods("GET")
	router.HandleFunc("/products/{productId}", productController.Update).Methods("PUT")
	router.HandleFunc("/products/{productId}", productController.FindById).Methods("GET")
	router.HandleFunc("/products/{productId}", productController.Delete).Methods("DELETE")

	router.HandleFunc("/orders", orderController.FindAllOrder).Methods("GET")
	router.HandleFunc("/orders", orderController.CreateOrder).Methods("POST")

	router.Use(middleware.RecoverMiddleware)

	return router
}
