package handler

import (
	"billing-service/internal/service"
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"time"
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

type GrantAccessRequest struct {
	UserID   int64     `json:"user_id" binding:"required"`
	Role string    `json:"user_role" binding:"required"`
	Amount   float64   `json:"amount" binding:"required"`
	Count    int       `json:"count" binding:"required"`
	Until    time.Time `json:"until" binding:"required"`
}

type AccessQueryRequest struct {
	UserID int64  `form:"user_id" binding:"required"`
	Role   string `form:"role" binding:"required"` // было UserRole
}

type CreateInvoiceDetailedRequest struct {
	UserID      int64     `json:"user_id" binding:"required"`
	UserRole    string    `json:"user_role" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Count       int       `json:"count" binding:"required"`
	Until       time.Time `json:"until" binding:"required"`
	OverallPrice float64  `json:"overallprice" binding:"required"`
}


type AccessHandler struct {
	service *service.PaymentService
}

func NewAccessHandler(s *service.PaymentService) *AccessHandler {
	return &AccessHandler{service: s}
}

func (h *AccessHandler) GrantAccess(c *gin.Context) {
	var req GrantAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.GrantAccess(c.Request.Context(), service.AccessData{
		UserID: req.UserID,
		Role:   req.Role,
		Amount: req.Amount,
		Count:  req.Count,
		Until:  req.Until,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grant access"})
		return
	}




	c.JSON(http.StatusOK, gin.H{"status": "access granted"})
}

func (h *AccessHandler) GetAccess(c *gin.Context) {
	var query AccessQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetAccess(c.Request.Context(), query.UserID, query.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get access info"})
		return
	}

	c.JSON(http.StatusOK, result)
}


func (h *PaymentHandler) CreateInvoiceDetailed(c *gin.Context) {
	var req CreateInvoiceDetailedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Здесь ты можешь сохранить данные как нужно. Например, просто лог:
	fmt.Printf("Создание инвойса: %+v\n", req)

	// Или сохранить в таблицу `payments`, если хочешь — скажи, добавим сохранение.

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}