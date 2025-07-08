package handler

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AirbaHandler struct {
	client *airba.Client
}

func NewAirbaHandler(client *airba.Client) *AirbaHandler {
	return &AirbaHandler{client: client}
}

// POST /airba/payments
func (h *AirbaHandler) CreatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Авторизация перед выполнением запроса
	if err := h.client.Authorize(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authorize: " + err.Error()})
		return
	}

	resp, err := h.client.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
