package handler

import (
	"billing-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CardsHandler struct {
	cardService   *service.CardService
	refundService *service.RefundService
}

func NewCardsHandler(card *service.CardService, refund *service.RefundService) *CardsHandler {
	return &CardsHandler{
		cardService:   card,
		refundService: refund,
	}
}

// POST /cards { account_id: string }
func (h *CardsHandler) AddCard(c *gin.Context) {
	var body struct {
		AccountID string `json:"account_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, err := h.cardService.AddCard(c, body.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"redirect_url": url})
}

// GET /cards/:accountId
func (h *CardsHandler) ListCards(c *gin.Context) {
	accountId := c.Param("accountId")
	cards, err := h.cardService.GetCards(c, accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cards)
}

// DELETE /cards/:id
func (h *CardsHandler) DeleteCard(c *gin.Context) {
	cardId := c.Param("id")
	err := h.cardService.DeleteCard(c, cardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /refund { payment_id: string, amount: float64 }
func (h *CardsHandler) Refund(c *gin.Context) {
	var body struct {
		PaymentID string  `json:"payment_id"`
		Amount    float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.refundService.Refund(c, body.PaymentID, body.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
