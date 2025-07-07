package handler

import (
	"billing-service/internal/model"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"billing-service/internal/service"
)

type WebhookHandler struct {
	Service *service.Services
	PubKey  *rsa.PublicKey
}


type WebhookPaymentPayload struct {
	ID          string `json:"id"`
	InvoiceID   string `json:"invoice_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Signature   string `json:"signature"`
}

func NewWebhookHandler(s *service.Services) (*WebhookHandler, error) {
	pubKey, err := loadPublicKeyFromURL("https://sps.airbapay.kz/acquiring/sign/public.pem")
	if err != nil {
		return nil, err
	}

	return &WebhookHandler{
		Service: s,
		PubKey:  pubKey,
	}, nil
}

func (h *WebhookHandler) HandlePaymentWebhook(c *gin.Context) {
	var payload WebhookPaymentPayload

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	message := payload.ID + payload.InvoiceID + toStr(payload.Amount) + payload.Currency + payload.Status + payload.Description
	if err := verifySignature(h.PubKey, message, payload.Signature); err != nil {
		log.Println("Signature verification failed:", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	err = h.Service.UpdatePaymentStatus(c.Request.Context(), payload.InvoiceID, payload.Status)
	if err != nil {
		log.Println("DB update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func toStr(i int) string {
	return fmt.Sprintf("%d", i)
}

func verifySignature(pubKey *rsa.PublicKey, message, signature string) error {
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sigBytes)
}

func loadPublicKeyFromURL(url string) (*rsa.PublicKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	block, _ := pem.Decode(io.ReadAll(resp.Body))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM format")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA key")
	}

	return rsaPub, nil
}
type WebhookCardPayload struct {
	AccountID string `json:"account_id"`
	CardID    string `json:"card_id"`
	Status    string `json:"status"`
	Signature string `json:"signature"`
}

func (h *WebhookHandler) HandleSaveCardWebhook(c *gin.Context) {
	var payload WebhookCardPayload

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	message := payload.AccountID + payload.CardID + payload.Status
	if err := verifySignature(h.PubKey, message, payload.Signature); err != nil {
		log.Println("Card signature verification failed:", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// ✅ Сохраняем карту в базу данных
	card := model.SavedCard{
		AccountID: payload.AccountID,
		CardID:    payload.CardID,
		Status:    payload.Status,
	}

	err = h.Service.Card.SaveCard(c.Request.Context(), card)
	if err != nil {
		log.Println("Failed to save card:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save card"})
		return
	}

	log.Printf("Card saved for AccountID: %s, CardID: %s, Status: %s\n", payload.AccountID, payload.CardID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
