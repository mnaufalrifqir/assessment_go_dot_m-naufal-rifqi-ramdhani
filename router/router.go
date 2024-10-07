package router

import (
	"api-dot/auth"
	"api-dot/database"
	"api-dot/handler"
	"api-dot/helper"
	"api-dot/payment"
	"api-dot/product"
	"api-dot/transaction"
	"api-dot/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := uint(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

func SetupRouter(mode string) *gin.Engine {
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	database.ConnectDB()
	database.InitialMigration()
	database.InitRedis()

	userRepository := user.NewRepository(database.DB)
	productRepository := product.NewRepository(database.DB, database.RDB, database.Ctx)
	transactionRepository := transaction.NewRepository(database.DB, database.RDB, database.Ctx)

	userService := user.NewService(userRepository)
	productService := product.NewService(productRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, paymentService, productRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router.Use(cors.Default())
	router.Use(logger.SetLogger())

	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)

	api.GET("/products", productHandler.GetAllProduct)
	api.GET("/product/:id", productHandler.GetProduct)
	api.POST("/product", authMiddleware(authService, userService), productHandler.CreateProduct)
	api.PUT("/product/:id", authMiddleware(authService, userService), productHandler.UpdateProduct)
	api.DELETE("/product/:id", authMiddleware(authService, userService), productHandler.DeleteProduct)

	api.POST("/transaction", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetTransactions)
	api.GET("/transaction/:id", authMiddleware(authService, userService), transactionHandler.GetTransactionByID)
	api.GET("/transactions/user", authMiddleware(authService, userService), transactionHandler.GetTransactionsUser)
	api.POST("/transactions/notification", transactionHandler.GetNotification)


	return router
}
