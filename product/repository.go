package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll()([]Product, error)
	FindByID(ID uint) (Product, error)
	Save(product Product) (Product, error)
	Update(product Product) (Product, error)
	DecreaseStock(product Product, quantity int) (Product, error)
	DeleteByID(ID uint) (Product, error)
}

type repository struct {
	db *gorm.DB
	rdb *redis.Client
	ctx context.Context
}

func NewRepository(db *gorm.DB, rdb *redis.Client, ctx context.Context) *repository {
	return &repository{db, rdb, ctx}
}

func (r *repository) FindAll() ([]Product, error) {
	var products []Product

	cachedProducts, err := r.rdb.Get(r.ctx, "products").Result()
	if err == redis.Nil {
		err := r.db.Find(&products).Error
		if err != nil {
			return products, err
		}

		productsJSON, err := json.Marshal(products)
		if err != nil {
			return products, err
		}

		err = r.rdb.Set(r.ctx, "products", productsJSON, 0).Err()
		if err != nil {
			return products, err
		}
	} else if err != nil {
		return products, err
	} else {
		err = json.Unmarshal([]byte(cachedProducts), &products)
		if err != nil {
			return products, err
		}
	}
	
	return products, nil
}

func (r *repository) FindByID(ID uint) (Product, error) {
	var product Product

	cachedProduct, err := r.rdb.Get(r.ctx, "product_"+ fmt.Sprint(ID)).Result()
	if err == redis.Nil {
		err := r.db.Where("id = ?", ID).Find(&product).Error
		if err != nil {
			return product, err
		}

		productJSON, err := json.Marshal(product)
		if err != nil {
			return product, err
		}

		err = r.rdb.Set(r.ctx, "product_"+ fmt.Sprint(ID), productJSON, 0).Err()
		if err != nil {
			return product, err
		}
	} else if err != nil {
		return product, err
	} else {
		err = json.Unmarshal([]byte(cachedProduct), &product)
		if err != nil {
			return product, err
		}
	}

	return product, nil
}

func (r *repository) Save(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	
	if err != nil {
		return product, err
	}

	err = r.rdb.Del(r.ctx, "products").Err()
    if err != nil {
        fmt.Println("Failed to delete cache:", err)
    }

	return product, nil
}

func (r *repository) Update(product Product) (Product, error) {
	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	err = r.rdb.Del(r.ctx, "products").Err()
    if err != nil {
        fmt.Println("Failed to delete cache:", err)
    }

	err = r.rdb.Del(r.ctx, "product_"+ fmt.Sprint(product.ID)).Err()
	if err != nil {
		fmt.Println("Failed to delete cache:", err)
	}

	return product, nil
}

func (r *repository) DecreaseStock(product Product, quantity int) (Product, error) {
	product.Stock = product.Stock - quantity
	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	err = r.rdb.Del(r.ctx, "product_"+ fmt.Sprint(product.ID)).Err()

	if err != nil {
		fmt.Println("Failed to delete cache:", err)
	}

	return product, nil
}

func (r *repository) DeleteByID(ID uint) (Product, error) {
	var product Product

	err := r.db.Where("id = ?", ID).Delete(&product).Error
	if err != nil {
		return product, err
	}

	err = r.rdb.Del(r.ctx, "product_"+ fmt.Sprint(product.ID)).Err()

	if err != nil {
		fmt.Println("Failed to delete cache:", err)
	}

	err = r.rdb.Del(r.ctx, "products").Err()
	if err != nil {
		fmt.Println("Failed to delete cache:", err)
	}

	return product, nil
}

