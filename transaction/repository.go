package transaction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	BeginTransaction() *gorm.DB
	GetByID(ID uint) (Transaction, error)
	Save(transaction Transaction, tx *gorm.DB) (Transaction, error)
	Update(transaction  Transaction, tx *gorm.DB) (Transaction, error)
	GetAll() ([]Transaction, error)
	GetByUserID(UserID uint) ([]Transaction, error)
	SaveTransactionDetail(transactionDetail TransactionDetails, tx *gorm.DB) (TransactionDetails, error)
}

type repository struct {
	db *gorm.DB
	rdb *redis.Client
	ctx context.Context
}

func NewRepository(db *gorm.DB, rdb *redis.Client, ctx context.Context) *repository {
	return &repository{db, rdb, ctx}
}

func (r *repository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *repository) GetByID(ID uint) (Transaction, error){
	var transaction Transaction

	cachedTransaction, err := r.rdb.Get(r.ctx, "transaction_"+ fmt.Sprint(ID)).Result()
	if err == redis.Nil {
		err := r.db.Where("id = ?", ID).Preload("TransactionDetails.Product").Preload("User").Find(&transaction).Error
		if err != nil {
			return transaction, err
		}

		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return transaction, err
		}

		err = r.rdb.Set(r.ctx, "transaction_"+ fmt.Sprint(ID), transactionJSON, 0).Err()
		if err != nil {
			return transaction, err
		}
	} else if err != nil {
		return transaction, err
	} else {
		err = json.Unmarshal([]byte(cachedTransaction), &transaction)
		if err != nil {
			return transaction, err
		}
	}

	return transaction, nil
}

func (r *repository) Save(transaction Transaction, tx *gorm.DB) (Transaction, error){
	err := tx.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	err = r.rdb.Del(r.ctx, "transactions").Err()
	if err != nil {
		return transaction, err
	}

	err = r.rdb.Del(r.ctx, "transactions_user_"+ fmt.Sprint(transaction.UserID)).Err()
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction, tx *gorm.DB) (Transaction, error) {
	err := tx.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	err = r.rdb.Del(r.ctx, "transactions").Err()
	if err != nil {
		return transaction, err
	}

	err = r.rdb.Del(r.ctx, "transactions_user_"+ fmt.Sprint(transaction.UserID)).Err()
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetAll() ([]Transaction, error) {
	var transactions []Transaction

	cachedTransactions, err := r.rdb.Get(r.ctx, "transactions").Result()
	if err == redis.Nil {
		err := r.db.Preload("TransactionDetails.Product").Preload("User").Find(&transactions).Error
		if err != nil {
			return transactions, err
		}

		transactionsJSON, err := json.Marshal(transactions)
		if err != nil {
			return transactions, err
		}
		
		err = r.rdb.Set(r.ctx, "transactions", transactionsJSON, 0).Err()
		if err != nil {
			return transactions, err
		}
	} else if err != nil {
		return transactions, err
	} else {
		err = json.Unmarshal([]byte(cachedTransactions), &transactions)
		if err != nil {
			return transactions, err
		}
	}

	return transactions, nil
}

func (r *repository) GetByUserID(UserID uint) ([]Transaction, error) {
	var transactions []Transaction

	cachedTransactions, err := r.rdb.Get(r.ctx, "transactions_user_"+ fmt.Sprint(UserID)).Result()
	if err == redis.Nil {
		err := r.db.Where("user_id = ?", UserID).Preload("TransactionDetails.Product").Preload("User").Find(&transactions).Error
		if err != nil {
			return transactions, err
		}

		transactionsJSON, err := json.Marshal(transactions)
		if err != nil {
			return transactions, err
		}

		err = r.rdb.Set(r.ctx, "transactions_user_"+ fmt.Sprint(UserID), transactionsJSON, 0).Err()
		if err != nil {
			return transactions, err
		}
	} else if err != nil {
		return transactions, err
	} else {
		err = json.Unmarshal([]byte(cachedTransactions), &transactions)
		if err != nil {
			return transactions, err
		}
	}

	return transactions, nil
}

func (r *repository) SaveTransactionDetail(transactionDetail TransactionDetails, tx *gorm.DB) (TransactionDetails, error) {
	err := tx.Create(&transactionDetail).Error
	if err != nil {
		return transactionDetail, err
	}

	return transactionDetail, nil
}
