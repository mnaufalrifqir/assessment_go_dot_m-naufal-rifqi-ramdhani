package handler

import (
	"api-dot/helper"
	"api-dot/transaction"
	"api-dot/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.transactionService.ProcessPayment(input)
	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	currentUser := c.MustGet("currentUser").(user.User)

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to process transaction", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = currentUser

	newTransaction, err := h.transactionService.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to process transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to process transaction", http.StatusOK, "success", newTransaction)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	if currentUser.Role != "admin" {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	transactions, err := h.transactionService.GetTransactions()
	if err != nil {
		response := helper.APIResponse("Failed to get transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatTransactions(transactions)

	response := helper.APIResponse("List of transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetTransactionsUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.transactionService.GetTransactionsUser(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatTransactions(transactions)

	response := helper.APIResponse("List of transactions", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetTransactionByID(c *gin.Context) {
	var input transaction.GetTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactionData, err := h.transactionService.GetTransactionByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := transaction.FormatTransaction(transactionData)

	response := helper.APIResponse("Transaction detail", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}