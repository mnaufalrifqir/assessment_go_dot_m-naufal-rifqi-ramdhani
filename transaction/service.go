package transaction

import (
	"api-dot/payment"
	"api-dot/product"
	"errors"
	"log"
	"strconv"
	"strings"
)

type Service interface {
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	GetTransactions() ([]Transaction, error)
	GetTransactionsUser(userID uint) ([]Transaction, error)
	GetTransactionByID(input GetTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

type service struct {
	repository     Repository
	paymentService payment.Service
	productRepository product.Repository
}

func NewService(repository Repository, paymentService payment.Service, productRepository product.Repository) *service {
	return &service{repository, paymentService, productRepository}
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	tx := s.repository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := Transaction{}
	
	newProduct, err := s.productRepository.FindByID(input.TransactionDetails.ProductID)
	if err != nil {
		tx.Rollback()
		return transaction, err
	}
	
	transaction.Amount = float64(input.TransactionDetails.Quantity) * newProduct.Price
	transaction.PaymentStatus = "pending"
	transaction.ShippingStatus = "unprocessed"
	transaction.UserID = input.User.ID

	newTransaction, err := s.repository.Save(transaction, tx)
	if err != nil {
		tx.Rollback()
		return newTransaction, err
	}

	transactionDetails := TransactionDetails{}
	transactionDetails.ProductID = input.TransactionDetails.ProductID
	transactionDetails.TransactionID = newTransaction.ID
	transactionDetails.Quantity = input.TransactionDetails.Quantity
	transactionDetails.Product = newProduct

	newTransactionDetails, err := s.repository.SaveTransactionDetail(transactionDetails, tx)
	if err != nil {
		tx.Rollback()
		return newTransaction, err
	}

	newTransaction.TransactionDetails = newTransactionDetails

	paymentTransaction := payment.Transaction{
		ID: newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		tx.Rollback()
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	newTransaction.User = input.User

	newTransaction, err = s.repository.Update(newTransaction, tx)
	if err != nil {
		tx.Rollback()
		return newTransaction, err
	}

	tx.Commit()

	return newTransaction, nil
}

func (s *service) GetTransactions() ([]Transaction, error) {
	transactions, err := s.repository.GetAll()
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsUser(userID uint) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByID(input GetTransactionInput) (Transaction, error) {
	transaction, err := s.repository.GetByID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	tx := s.repository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var strID string

	parts := strings.Split(input.OrderID, "_")
	if len(parts) > 1 {
        strID = parts[1]
    } else {
		return errors.New("invalid order id")
	}

	transaction_id, _ := strconv.Atoi(strID)
	log.Println(transaction_id)

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

	updatedTransaction, err := s.repository.Update(transaction, tx)
	if err != nil{
		return err
	}

	product, err := s.productRepository.FindByID(updatedTransaction.TransactionDetails.ProductID)
	if err != nil {
		return err
	}

	if updatedTransaction.PaymentStatus == "paid" {
		product.Stock = product.Stock - updatedTransaction.TransactionDetails.Quantity
		_, err = s.productRepository.Update(product)
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}
