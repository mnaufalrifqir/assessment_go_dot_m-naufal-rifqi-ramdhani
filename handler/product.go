package handler

import (
	"api-dot/helper"
	"api-dot/product"
	"net/http"

	"github.com/gin-gonic/gin"
)
type productHandler struct {
	productService product.Service
}

func NewProductService(productService product.Service) *productHandler {
	return &productHandler{productService}
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var input product.CreateProductInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newProduct, err := h.productService.CreateProduct(input)
	if err != nil {
		response := helper.APIResponse("Create product failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := product.FormatProduct(newProduct)

	response := helper.APIResponse("Create product success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) GetProduct(c *gin.Context) {
	var input product.GetProductDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newProduct, err := h.productService.GetProduct(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := product.FormatProduct(newProduct)

	response := helper.APIResponse("Product detail", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) GetAllProduct(c *gin.Context) {
	newProducts, err := h.productService.GetAllProduct()
	if err != nil {
		response := helper.APIResponse("Failed to get products", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var formatter []product.ProductFormatter
	for _, newProduct := range newProducts {
		formatter = append(formatter, product.FormatProduct(newProduct))
	}

	response := helper.APIResponse("List of products", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	var inputID product.GetProductDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData product.CreateProductInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := product.FormatProduct(updatedProduct)

	response := helper.APIResponse("Update product success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	var input product.GetProductDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to delete product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deletedProduct, err := h.productService.DeleteProduct(input)
	if err != nil {
		response := helper.APIResponse("Failed to delete product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := product.FormatProduct(deletedProduct)

	response := helper.APIResponse("Product has been deleted", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}