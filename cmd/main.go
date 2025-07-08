package main

import (
	"billing-service/internal/airba"
	"billing-service/internal/config"
	"billing-service/internal/handler"
	"billing-service/internal/model" // если структура Repository лежит тут
	"billing-service/internal/repository"
	"billing-service/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// 🔧 Создаём Airba-клиент
	airbaClient := airba.NewClient(cfg.AirbaAPIKey, cfg.AirbaBaseURL)

	// 🔧 Репозиторий
	repo := &model.Repository{DB: db}

	// 🔧 Создаём сервис статуса платежа
	paymentStatusService := service.NewPaymentStatusService(airbaClient, repo)

	// 🔧 Основной сервис платежей
	paymentService := service.NewPaymentService(repo)

	// 🔧 Webhook handler
	webhookHandler, err := handler.NewWebhookHandler(paymentService)
	if err != nil {
		log.Fatalf("failed to create webhook handler: %v", err)
	}

	router := gin.Default()

	// 🔧 Настроить маршруты
	handler.SetupRoutes(router, repo, airbaClient, webhookHandler, paymentStatusService)


	log.Println("📦 Server is running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
