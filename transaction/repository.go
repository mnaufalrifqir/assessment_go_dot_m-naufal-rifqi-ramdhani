package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByID(ID uint) (Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction  Transaction) (Transaction, error)
	GetAll() ([]Transaction, error)
	GetByUserID(UserID uint) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByID(ID uint) (Transaction, error){
	var transaction Transaction

	err := r.db.Where("id = ?", ID).Preload("TransactionDetails").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error){
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetAll() ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("TransactionDetails").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserID(UserID uint) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Where("user_id = ?", UserID).Preload("TransactionDetails").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
