package controllers

import (
	"net/http"

	"zen-test/app/helpers"
	"zen-test/app/web"
	"zen-test/app/web/models"
	"zen-test/app/web/services"

	"github.com/gorilla/mux"
)

type ProductControllerImpl struct {
	ProductService services.ProductService
}

type ProductController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
}

func NewProductController(productService services.ProductService) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
	}
}

// Create Product godoc
// @Summary create Product for the store
// @Description create Product for the store
// @Tags Product
// @Accept json
// @Produce json
// @Param Product body models.ProductDto true "Product create"
// @Success 200 {object} web.WebResponse{data=models.ProductResponse}
// @Failure 401 {object} web.WebResponse
// @Router /products [post]
// @Security BearerAuth
func (c *ProductControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	productCreateRequest := models.ProductCreateUpdate{}
	helpers.ToRequestBody(r, &productCreateRequest)

	productResponse := c.ProductService.Create(r.Context(), productCreateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   productResponse,
	}

	helpers.WriteResponseBody(w, webResponse)
}

// Update Product godoc
// @Summary Update Product from the store
// @Description Update Product from the store
// @Tags Product
// @Accept json
// @Produce json
// @Param Product body models.ProductDto true "Product Update"
// @Param productId path string true "Product ID"
// @Success 200 {object} web.WebResponse{data=models.ProductResponse}
// @Failure 401 {object} web.WebResponse
// @Router /products/{productId} [put]
// @Security BearerAuth
func (c *ProductControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	productUpdateRequest := models.ProductCreateUpdate{}
	helpers.ToRequestBody(r, &productUpdateRequest)

	vars := mux.Vars(r)
	productId := vars["productId"]

	productResponse := c.ProductService.Update(r.Context(), productUpdateRequest, productId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   productResponse,
	}

	helpers.WriteResponseBody(w, webResponse)
}

// Delete Product godoc
// @Summary Delete Product from the store
// @Description Delete Product from the store
// @Tags Product
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} web.WebResponse
// @Failure 401 {object} web.WebResponse
// @Router /products/{productId} [delete]
// @Security BearerAuth
func (c *ProductControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]

	c.ProductService.Delete(r.Context(), productId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
	}
	helpers.WriteResponseBody(w, webResponse)
}

// FindById Product godoc
// @Summary FindById Product from the store
// @Description FindById Product from the store
// @Tags Product
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} web.WebResponse{data=models.ProductResponse}
// @Failure 401 {object} web.WebResponse
// @Router /products/{productId} [get]
// @Security BearerAuth
func (c *ProductControllerImpl) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]

	productResponse := c.ProductService.FindById(r.Context(), productId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   productResponse,
	}
	helpers.WriteResponseBody(w, webResponse)
}

// FindAll Products godoc
// @Summary FindAll Products from the store
// @Description FindAll Products from the store
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse{data=models.ProductResponse}
// @Failure 401 {object} web.WebResponse
// @Router /products [get]
// @Security BearerAuth
func (c *ProductControllerImpl) FindAll(w http.ResponseWriter, r *http.Request) {
	productResponse := c.ProductService.FindAll(r.Context())
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   productResponse,
	}
	helpers.WriteResponseBody(w, webResponse)
}
