package main

import (
	"log"
	_ "subscriptions/docs"
	"subscriptions/internal/config"
	api "subscriptions/internal/http"
	"subscriptions/internal/http/middleware"
	"subscriptions/internal/service"
	"subscriptions/internal/storage/postgres"
	"subscriptions/pkg/database"
	"subscriptions/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscriptions REST API
// @version 1.0
// @description Пример REST API сервиса для управления пользовательскими подписками
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Инициализация логгера
	appLogger := logger.New(cfg.LogLevel)

	// Подключение к БД
	db, err := database.NewPostgresConnection(cfg.DB)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	appLogger.Info("Database connected succesfully!")

	// Инициализация хранилища (в задании - базы данных, но может быть и другой вариант)
	subscriptionDB := postgres.NewSubscriptionStorage(db, appLogger)

	// Инициализация сервиса
	subscriptionService := service.NewSubscriptionService(subscriptionDB, appLogger)

	// Инициализация HTTP обработчиков
	handler := api.NewHandler(subscriptionService, appLogger)

	// Инициализация gin роутера
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware(appLogger))
	router.Use(gin.Recovery())

	// Настройка роутера - предположим, что потом будут и другие версии (так что пока v1)
	api := router.Group("/api/v1")
	handler.InitRoutes(api)

	// Настройка Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// health-check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Запуск сервера
	appLogger.Info("Starting server on port " + cfg.HTTP.Port)
	if err := router.Run(":" + cfg.HTTP.Port); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
