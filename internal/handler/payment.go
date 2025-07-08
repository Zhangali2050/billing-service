package handler

import (
	"billing-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type CreateInvoiceRequest struct {
	Role     string  `json:"role" binding:"required"`
	UserID   string  `json:"user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
}

func (h *PaymentHandler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	input := service.CreatePaymentInput{
		Role:     req.Role,
		UserID:   userUUID,
		Amount:   req.Amount,
		Quantity: req.Quantity,
	}

	err = h.paymentService.CreatePayment(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "invoice created",
	})
}

type PaymentHistoryRequest struct {
	Role   string `form:"role" binding:"required"`
	UserID string `form:"user_id" binding:"required"`
}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	var req PaymentHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	payments, err := h.paymentService.GetPayments(c.Request.Context(), userUUID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
