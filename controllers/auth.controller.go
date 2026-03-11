package controllers

import (
	"go-image-processing-api/requests"
	"go-image-processing-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Register(req.Username, req.Password)
	if err != nil {
		if err == services.ErrUsernameExist {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  "register success",
		"username": user.Username,
	})
}
