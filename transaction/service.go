package transaction

import (
	"api-dot/payment"
	"strconv"
)

type Service interface {
	ProcessPayment(input TransactionNotificationInput) error
}

type service struct {
	repository     Repository
	paymentService payment.Service
}

func NewService(repository Repository, paymentService payment.Service) *service {
	return &service{repository, paymentService}
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(uint(transaction_id))
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.PaymentStatus = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.PaymentStatus = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.PaymentStatus = "cancelled"
	}

	return nil
}
