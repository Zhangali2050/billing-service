package handler

import (
	"billing-service/internal/service"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Payload для платежного вебхука
type AirbaWebhookPayload struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"invoice_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	Signature   string  `json:"signature"`
}

// Payload для вебхука карты
type AirbaCardWebhookPayload struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	MaskedPan string `json:"masked_pan"`
	Name      string `json:"name"`
	Expire    string `json:"expire"`
	Signature string `json:"sign"`
}

type WebhookHandler struct {
	PaymentService *service.PaymentService
	PublicKeyURL   string
}

func NewWebhookHandler(paymentService *service.PaymentService) *WebhookHandler {
	return &WebhookHandler{
		PaymentService: paymentService,
		PublicKeyURL:   "https://sps.airbapay.kz/acquiring/sign/public.pem",
	}
}

// POST /webhook/payment
func (h *WebhookHandler) HandlePaymentWebhook(c *gin.Context) {
	var payload AirbaWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	ok, err := h.verifySignature(payload)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	err = h.PaymentService.UpdatePaymentStatus(
		c.Request.Context(),
		payload.InvoiceID,
		payload.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// POST /webhook/card
func (h *WebhookHandler) HandleCardWebhook(c *gin.Context) {
	var payload AirbaCardWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card payload"})
		return
	}

	ok, err := h.verifyCardSignature(payload)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid card signature"})
		return
	}

	// Симулируем сохранение карты
	log.Printf("✅ Received card for user %s: %s (%s)", payload.AccountID, payload.MaskedPan, payload.Expire)
	c.JSON(http.StatusOK, gin.H{"status": "card saved (simulated)"})
}

// Подпись платежа: id+invoice_id+amount+currency+status+description
func (h *WebhookHandler) verifySignature(p AirbaWebhookPayload) (bool, error) {
	pubKey, err := h.loadPublicKey()
	if err != nil {
		return false, err
	}
	message := fmt.Sprintf("%s%s%.2f%s%s%s",
		p.ID, p.InvoiceID, p.Amount, p.Currency, p.Status, p.Description)

	return h.verifyWithPublicKey(pubKey, message, p.Signature)
}

// Подпись карты: id+account_id+masked_pan+name+expire
func (h *WebhookHandler) verifyCardSignature(p AirbaCardWebhookPayload) (bool, error) {
	pubKey, err := h.loadPublicKey()
	if err != nil {
		return false, err
	}
	message := fmt.Sprintf("%s%s%s%s%s", p.ID, p.AccountID, p.MaskedPan, p.Name, p.Expire)

	return h.verifyWithPublicKey(pubKey, message, p.Signature)
}

// Получает RSA-публичный ключ
func (h *WebhookHandler) loadPublicKey() (*rsa.PublicKey, error) {
	resp, err := http.Get(h.PublicKeyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pemBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPubKey, nil
}

// Проверка подписи
func (h *WebhookHandler) verifyWithPublicKey(pubKey *rsa.PublicKey, message, base64Sig string) (bool, error) {
	sigBytes, err := base64.StdEncoding.DecodeString(base64Sig)
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes)
	return err == nil, nil
}
