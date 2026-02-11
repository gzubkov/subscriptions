package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "subscriptions/docs"
	"subscriptions/internal/config"
	api "subscriptions/internal/http"
	"subscriptions/internal/http/middleware"
	"subscriptions/internal/service"
	"subscriptions/internal/storage/postgres"
	"subscriptions/pkg/database"
	"subscriptions/pkg/logger"
	"syscall"
	"time"

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

	// Force log's color
	gin.ForceConsoleColor()

	// Инициализация gin роутера
	router := gin.Default()

	// Set mode
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(middleware.LoggerMiddleware(appLogger))
	router.Use(gin.Recovery())

	// Настройка CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Настройка роутера - предположим, что потом будут и другие версии (так что пока v1)
	api := router.Group("/api/v1")
	handler.InitRoutes(api)

	// Настройка Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// health-check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Запуск сервера с graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: router,
	}

	go func() {
		appLogger.Info("Starting server on port " + cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exiting")
}
