package handler

import (
	"billing-service/internal/airba"
	"billing-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// SetupRoutes инициализирует все маршруты
func SetupRoutes(
	r *gin.Engine,
	repo *repository.Repository,
	airbaClient *airba.Client,
	webhookHandler *WebhookHandler, // <-- теперь передаём извне
	paymentStatusService *service.PaymentStatusService,
) {
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

	// Webhooks
	r.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
	r.POST("/webhook/save-card", webhookHandler.HandleSaveCardWebhook)

	// Карты (отдельный handler)
	cardsHandler := NewCardsHandler(airbaClient)
	r.POST("/cards", cardsHandler.AddCard)
	r.GET("/cards/:accountId", cardsHandler.ListCards)
	r.DELETE("/cards/:id", cardsHandler.DeleteCard)

	r.GET("/status/:invoiceId", func(c *gin.Context) {
		invoiceID := c.Param("invoiceId")
		status, err := paymentStatusService.GetPaymentStatus(invoiceID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, status)
	})

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
}
