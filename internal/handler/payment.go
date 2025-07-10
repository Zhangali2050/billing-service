package handler

import (
	"billing-service/internal/service"
	"github.com/gin-gonic/gin"
	"fmt"
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
	UserID   int64   `json:"user_id" binding:"required"`  // теперь число
	Amount   float64 `json:"amount" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
}

func (h *PaymentHandler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := service.CreatePaymentInput{
		Role:     req.Role,
		UserID:   req.UserID,
		Amount:   req.Amount,
		Quantity: req.Quantity,
	}

	resp, err := h.paymentService.CreateAndSavePayment(c.Request.Context(), input)
	if err != nil {
		fmt.Println("CreateAndSavePayment error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}


type PaymentHistoryRequest struct {
	Role   string `form:"role" binding:"required"`
	UserID int64  `form:"user_id" binding:"required"` // ✅ тоже заменяем здесь
}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	var req PaymentHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ❌ больше не нужно uuid.Parse

	payments, err := h.paymentService.GetPayments(c.Request.Context(), req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
