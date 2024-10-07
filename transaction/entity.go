package transaction

import (
	"api-dot/product"
	"api-dot/user"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount         float64
	PaymentStatus  string
	ShippingStatus string
	PaymentURL     string
	User           user.User
	UserID         uint
	TransactionDetails TransactionDetails
}

type TransactionDetails struct {
	gorm.Model
	ProductID      uint  
	TransactionID  uint  
	Quantity       int   
	Product 	   product.Product
}