package handler

import (
	"billing-service/internal/airba"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AirbaHandler содержит зависимости для Airba Pay
type AirbaHandler struct {
	Client *airba.Client
}

// NewAirbaHandler конструктор
func NewAirbaHandler(client *airba.Client) *AirbaHandler {
	return &AirbaHandler{Client: client}
}

// POST /airba/auth
func (h *AirbaHandler) Auth(c *gin.Context) {
	err := h.Client.Authenticate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth failed", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": h.Client.AccessToken})
}

// POST /airba/payment
func (h *AirbaHandler) CreatePayment(c *gin.Context) {
	var req any
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.Client.GenericPost("/api/v2/payments", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment failed", "details": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", resp)
}

// POST /airba/cards
func (h *AirbaHandler) AddCard(c *gin.Context) {
	var req any
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.Client.GenericPost("/api/v1/cards", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "add card failed", "details": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", resp)
}

// GET /airba/cards/:accountId
func (h *AirbaHandler) GetCards(c *gin.Context) {
	accountId := c.Param("accountId")
	path := "/api/v1/cards/" + accountId

	resp, err := h.Client.GenericPost(path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get cards failed", "details": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", resp)
}

// DELETE /airba/cards/:id
func (h *AirbaHandler) DeleteCard(c *gin.Context) {
	cardId := c.Param("id")
	req, err := http.NewRequest("DELETE", h.Client.BaseURL+"/api/v1/cards/"+cardId, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", "Bearer "+h.Client.AccessToken)

	resp, err := h.Client.HTTPClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
}
