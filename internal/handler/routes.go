package handler

import (
	"billing-service/internal/model"
	"billing-service/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(r *gin.Engine, repo *repository.Repository) {
	r.POST("/roles", func(c *gin.Context) {
		var role model.Role
		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := repo.CreateRole(c.Request.Context(), role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.POST("/invoice", func(c *gin.Context) {
		var p model.Payment
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := repo.CreatePayment(c.Request.Context(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.POST("/payment", func(c *gin.Context) {
		var p model.Payment
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := repo.CreatePayment(c.Request.Context(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.GET("/payments", func(c *gin.Context) {
		payments, err := repo.GetAllPayments(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, payments)
	})

	r.GET("/payments/:id/:role", func(c *gin.Context) {
		id := c.Param("id")
		role := c.Param("role")
		payments, err := repo.GetPaymentsByID(c.Request.Context(), id, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, payments)
	})

}
