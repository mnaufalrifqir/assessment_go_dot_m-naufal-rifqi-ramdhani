package transaction

import (
	"api-dot/user"
	"api-dot/product"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount         int
	PaymentStatus  string
	ShippingStatus string
	PaymentURL     string
	User           user.User
	UserID         uint
	TransactionDetails []TransactionDetails
}

type TransactionDetails struct {
	ID             uint   `json:"id"`
	ProductID      uint   `json:"product_id"`
	TransactionID  uint   `json:"transaction_id"`
	Quantity       int    `json:"quantity"`
	Price          int    `json:"price"`
	Product 	   product.Product `json:"product"`
}