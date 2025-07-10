package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/sign-in", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "mock_access_token",
			"token_type":   "Bearer",
			"expires_in":   "3600",
		})
	})

	http.HandleFunc("/api/v2/payments", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"id":           "mock_payment_id",
			"invoice_id":   "mock_invoice_123",
			"redirect_url": "https://mock.airba/redirect/payment123",
		})
	})

	http.HandleFunc("/acquiring/sign/public.pem", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./mock_public.pem") // Смотри ниже, как сгенерировать
	})

	println("✅ Mock AirbaPay server running on http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}
