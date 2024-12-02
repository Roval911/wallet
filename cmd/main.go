package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"wallet/controllers"
	"wallet/db"
	_ "wallet/docs"
)

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Инициализация базы данных
	db.InitDb()

	// Запуск миграций
	db.RunMigrations()
}

// @title Wallet Service API
// @version 1.0
// @description API для управления операциями кошелька
// @host localhost:8080
// @BasePath /api/v1
func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/v1/wallets", controllers.CreateWallet)
	router.GET("/api/v1/wallets/:walletId", controllers.GetWalletBalance)
	router.POST("/api/v1/wallet", controllers.UpdateWallet)

	router.Run(os.Getenv("PORT"))
}
