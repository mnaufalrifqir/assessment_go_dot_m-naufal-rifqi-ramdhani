package transaction

import "api-dot/user"

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

type CreateTransactionInput struct {
	TransactionDetails CreateTransactionDetailInput `json:"transaction_details" binding:"required"`
	User               user.User
}

type CreateTransactionDetailInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type GetTransactionInput struct {
	ID uint `uri:"id" binding:"required"`
}
