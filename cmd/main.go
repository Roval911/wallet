package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"wallet/configs"
)

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Инициализация базы данных
	configs.InitDb()

	// Запуск миграций
	configs.RunMigrations()
}
func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(os.Getenv("PORT"))
}
