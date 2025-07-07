package handler

import (
	"billing-service/internal/airba"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CardsHandler обрабатывает сохранённые карты клиента
type CardsHandler struct {
	Client *airba.Client
}

// NewCardsHandler создает новый обработчик карт
func NewCardsHandler(client *airba.Client) *CardsHandler {
	return &CardsHandler{Client: client}
}

// POST /cards — сохранить карту
func (h *CardsHandler) AddCard(c *gin.Context) {
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

// GET /cards/:accountId — получить сохраненные карты клиента
func (h *CardsHandler) ListCards(c *gin.Context) {
	accountId := c.Param("accountId")
	resp, err := h.Client.GenericGet("/api/v1/cards/" + accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get cards failed", "details": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", resp)
}

// DELETE /cards/:id — удалить карту
func (h *CardsHandler) DeleteCard(c *gin.Context) {
	cardId := c.Param("id")
	err := h.Client.GenericDelete("/api/v1/cards/" + cardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete card failed", "details": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
