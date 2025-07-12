package handler

import (
	"billing-service/internal/airba"
	"billing-service/internal/repository"
	"billing-service/internal/service"
	"fmt"
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
	accessHandler := NewAccessHandler(paymentService) // ✅ новый хендлер

	// Роуты для API
	r.POST("/invoice", paymentHandler.CreateInvoice)
	r.GET("/payments", paymentHandler.GetPaymentHistory)

	// ✅ Новые маршруты подлоги
	r.POST("/api/payment/post/invoice", accessHandler.GrantAccess)
	r.GET("/api/payment/get/access", accessHandler.GetAccess)

	r.POST("/cards", cardsHandler.AddCard)
	r.GET("/cards/:accountId", cardsHandler.ListCards)
	r.DELETE("/cards/:id", cardsHandler.DeleteCard)

	r.POST("/api/payment/post/invoice/create", paymentHandler.CreateInvoiceDetailed)


	r.POST("/refund", cardsHandler.Refund)

	// Вебхуки от AirbaPay
	r.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
	r.POST("/webhook/card", webhookHandler.HandleCardWebhook)


}

func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		fmt.Println("HEADER:", c.GetHeader("X-Api-Key"))
	}
}
