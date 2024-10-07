package database

import (
	"api-dot/product"
	"api-dot/transaction"
	"api-dot/user"
	"api-dot/utils"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectDB() {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		utils.GetConfig("DB_HOST"),
		utils.GetConfig("DB_USER"),
		utils.GetConfig("DB_PASSWORD"),
		utils.GetConfig("DB_NAME"),
		utils.GetConfig("DB_PORT"),
		utils.GetConfig("DB_SSL_MODE"),
		utils.GetConfig("DB_TIMEZONE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func InitialMigration() {
	err := DB.AutoMigrate(&user.User{}, &product.Product{}, &transaction.Transaction{}, &transaction.TransactionDetails{})


	if err != nil {
		log.Printf("Error when migrating the database: %v", err)
	}

	SeedDatabase()
}

func SeedDatabase() {
	var countProduct, countUser int64

	DB.Model(&user.User{}).Count(&countUser)
	DB.Model(&product.Product{}).Count(&countProduct)

	if countUser == 0 {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)
		if err != nil {
			log.Printf("Error seeding dummy products: %v", err)
		}

		user := user.User{
			Name:         "Admin",
			Email:        "admin@example.com",
			PasswordHash: string(passwordHash),
			Role:         "admin",
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Error seeding dummy products: %v", err)
		} else {
			log.Println("Dummy products added successfully!")
		}
	}

	if countProduct == 0 {
		products := []product.Product{
			{
				Name:  "Product 1",
				Price: 10000,
				Stock: 10,
			},
			{
				Name:  "Product 2",
				Price: 15000,
				Stock: 5,
			},
			{
				Name:  "Product 3",
				Price: 7500,
				Stock: 20,
			},
		}

		for _, product := range products {
			if err := DB.Create(&product).Error; err != nil {
				log.Printf("Error seeding dummy products: %v", err)
			}
		}
	}
}
