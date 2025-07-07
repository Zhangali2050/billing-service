package main

import (
	"billing-service/internal/handler"
	"billing-service/internal/repository"
	"billing-service/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Получаем строку подключения из переменных окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Подключаемся к PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Инициализируем репозиторий и сервисы
	repo := repository.NewRepository(db)
	apiService := service.NewAirbaPayService()

	// Создаём роутер Gin
	router := gin.Default()

	// Статический маршрут для index.html
	router.StaticFile("/docs", "./static/index.html")

	// Авторизация API-ключом для всех защищённых маршрутов
	protected := router.Group("/", func(c *gin.Context) {
		key := c.GetHeader("X-Api-Key")
		if key != "sandbox_123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	})

	// Регистрация маршрутов
	handler.SetupRoutes(protected, repo, apiService)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on http://localhost:%s/", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
