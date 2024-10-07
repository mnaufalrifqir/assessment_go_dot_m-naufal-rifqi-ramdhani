package transaction

import "api-dot/product"

type TransactionFormatter struct {
	ID             uint   `json:"id"`
	Amount         float64    `json:"amount"`
	PaymentStatus  string `json:"payment_status"`
	ShippingStatus string `json:"shipping_status"`
	PaymentURL     string `json:"payment_url"`
	User           TransactionUserFormatter `json:"user"`
	TransactionDetails TransactionDetailFormatter `json:"transaction_details"`
}

type TransactionDetailFormatter struct {
	ID       uint        `json:"id"`
	Quantity int         `json:"quantity"`
	Product  product.ProductFormatter `json:"product"`
}

type TransactionUserFormatter struct {
	ID             uint   `json:"id"`
	Name		   string `json:"name"`
	Email		   string `json:"email"`
	Role		   string `json:"role"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		ID:             transaction.ID,
		Amount:         transaction.Amount,
		PaymentStatus:  transaction.PaymentStatus,
		ShippingStatus: transaction.ShippingStatus,
		PaymentURL:     transaction.PaymentURL,
		User:           TransactionUserFormatter{
			ID: transaction.User.ID,
			Name: transaction.User.Name,
			Email: transaction.User.Email,
			Role: transaction.User.Role,
		},
		TransactionDetails: TransactionDetailFormatter{
			ID: transaction.TransactionDetails.ID,
			Quantity: transaction.TransactionDetails.Quantity,
			Product: product.ProductFormatter{
				ID: transaction.TransactionDetails.Product.ID,
				Name: transaction.TransactionDetails.Product.Name,
				Price: transaction.TransactionDetails.Product.Price,
				Stock: transaction.TransactionDetails.Product.Stock,
			},
		},
	}

	return formatter
}

func FormatTransactions(transactions []Transaction) []TransactionFormatter {
	transactionsFormatter := []TransactionFormatter{}

	for _, transaction := range transactions {
		formatter := FormatTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}
