package main

import (
	"billing-service/internal/handler"
	"billing-service/internal/repository"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db: ", err)
	}

	repo := repository.NewRepository(db)
	r := gin.Default()

	// 🔐 Защита всех маршрутов API
	r.Use(func(c *gin.Context) {
		apiKey := c.GetHeader("X-Api-Key")
		if apiKey != "sandbox_123" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	})

	// 📄 Документация без авторизации
	r.Static("/docs", "./static")

	// API
	handler.SetupRoutes(r, repo)

	log.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

