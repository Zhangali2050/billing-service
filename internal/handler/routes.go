package handler

import (
	"billing-service/internal/airba"
	"billing-service/internal/repository"
	"billing-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	repo *repository.Repository,
	airbaClient *airba.Client,
	webhookHandler *WebhookHandler,
) {
	r.Use(apiKeyMiddleware())

	// Сервисы
	paymentService := service.NewPaymentService(repo, airbaClient)
	cardService := service.NewCardService(airbaClient)
	refundService := service.NewRefundService(airbaClient)

	// Обработчики
	paymentHandler := NewPaymentHandler(paymentService)
	cardsHandler := NewCardsHandler(cardService, refundService)

	// Роуты для API
	r.POST("/invoice", paymentHandler.CreateInvoice)
	r.GET("/payments", paymentHandler.GetPaymentHistory)

	r.POST("/cards", cardsHandler.AddCard)
	r.GET("/cards/:accountId", cardsHandler.ListCards)
	r.DELETE("/cards/:id", cardsHandler.DeleteCard)

	r.POST("/refund", cardsHandler.Refund)

	// Вебхуки от AirbaPay
	r.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
	r.POST("/webhook/card", webhookHandler.HandleCardWebhook)

	// Статичный HTML для теста (браузерный UI)
	r.StaticFile("/", "./static/index.html")
}

// // Простая проверка API-ключа
// func apiKeyMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		apiKey := c.GetHeader("X-Api-Key")
// 		if apiKey != "sandbox_123" {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
// 			return
// 		}
// 		c.Next()
// 	}
// }

func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем доступ к корневому HTML (index.html) без API-ключа
		if c.Request.Method == "GET" && c.Request.URL.Path == "/" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-Api-Key")
		if apiKey != "sandbox_123" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
