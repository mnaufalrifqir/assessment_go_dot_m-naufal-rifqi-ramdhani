package transaction

import (
	"api-dot/user"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount         int
	PaymentStatus  string
	ShippingStatus string
	Code           string
	PaymentURL     string
	User           user.User
	UserID         uint
}
