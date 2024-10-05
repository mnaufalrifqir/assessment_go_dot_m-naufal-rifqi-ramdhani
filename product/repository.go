package product

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(product Product) (Product, error)
	FindByID(ID uint) (Product, error)
	Update(product Product) (Product, error)
	DecreaseStock(product Product, quantity int) (Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) FindByID(ID uint) (Product, error) {
	var product Product

	err := r.db.Where("id = ?", ID).Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(product Product) (Product, error) {
	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) DecreaseStock(product Product, quantity int) (Product, error) {
	product.Stock = product.Stock - quantity
	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

