package main

import (
	"billing-service/internal/airba"
	"billing-service/internal/config"
	"billing-service/internal/handler"
	"billing-service/internal/repository"
	"billing-service/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ failed to load config: %v", err)
	}

	db, err := repository.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	airbaClient := airba.NewClient(
		cfg.Airba.User,
		cfg.Airba.Password,
		cfg.Airba.TerminalID,
		cfg.Airba.BaseURL,
		cfg.Airba.SignatureKey,

	)


	paymentService := service.NewPaymentService(repo, airbaClient)
	webhookHandler := handler.NewWebhookHandler(paymentService)

	router := gin.Default()
	handler.SetupRoutes(router, repo, airbaClient, webhookHandler)

	log.Println("🚀 Server is running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("❌ failed to start server: %v", err)
	}
}
