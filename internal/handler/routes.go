package handler

import (
	"billing-service/internal/airba"
	"billing-service/internal/repository"
	"billing-service/internal/handler"
	"billing-service/internal/service"
	"crypto/rsa"
	"log"
	"os"
	"github.com/gin-gonic/gin"
)

// SetupRoutes инициализирует все маршруты
func SetupRoutes(r *gin.Engine, repo *repository.Repository, airbaClient *airba.Client) {
	// Middleware авторизации по X-Api-Key
	r.Use(apiKeyMiddleware())

	// Статическая документация
	r.StaticFile("/docs", "./static/index.html")

	// Основные роуты (биллинг)
	setupBillingRoutes(r, repo)

	// AirbaPay
	paymentHandler := NewAirbaHandler(airbaClient)
	r.POST("/airba/payments", paymentHandler.CreatePayment)
	r.POST("/airba/payments/:cardId", paymentHandler.CreatePaymentWithCard)
	r.POST("/airba/cards", paymentHandler.AddCard)
	r.GET("/airba/cards/:accountId", paymentHandler.GetCards)
	r.DELETE("/airba/cards/:id", paymentHandler.DeleteCard)
	r.POST("/airba/webhook", paymentHandler.Webhook)

	// Карты (отдельный handler)
	cardsHandler := NewCardsHandler(airbaClient)
	r.POST("/cards", cardsHandler.AddCard)
	r.GET("/cards/:accountId", cardsHandler.ListCards)
	r.DELETE("/cards/:id", cardsHandler.DeleteCard)
	r.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
	r.POST("/webhook/save-card", webhookHandler.HandleSaveCardWebhook)

}

func setupBillingRoutes(r *gin.Engine, repo *repository.Repository) {
	r.POST("/roles", func(c *gin.Context) {
		var role repository.Role
		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := repo.CreateRole(c, role); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(201)
	})

	r.POST("/invoice", func(c *gin.Context) {
		var p repository.Payment
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := repo.CreatePayment(c, p); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(201)
	})

	r.POST("/payment", func(c *gin.Context) {
		var p repository.Payment
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := repo.CreatePayment(c, p); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(201)
	})

	r.GET("/payments", func(c *gin.Context) {
		payments, err := repo.GetAllPayments(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, payments)
	})

	r.GET("/payments/:id/:role", func(c *gin.Context) {
		id := c.Param("id")
		role := c.Param("role")
		payments, err := repo.GetPaymentsByID(c, id, role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, payments)
	})
}

func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Api-Key")
		if apiKey != "sandbox_123" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
	// Загружаем публичный ключ для проверки подписи webhook
	pubKeyPath := "./airba_public_key.pem"
	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("не удалось прочитать публичный ключ: %v", err)
	}

	publicKey, err := handler.ParseRSAPublicKey(pubKeyBytes)
	if err != nil {
		log.Fatalf("не удалось распарсить публичный ключ: %v", err)
	}

	if webhookHandler, err := handler.NewWebhookHandler(paymentService); err == nil {
		r.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
	} else {
		log.Fatalf("failed to load webhook handler: %v", err)
	}

	r.POST("/webhook/save-card", webhookHandler.HandleSaveCardWebhook)


}
