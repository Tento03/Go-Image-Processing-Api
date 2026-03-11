package controllers

import (
	"go-image-processing-api/requests"
	"go-image-processing-api/services"
	"net/http"
	"os"

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

func Login(c *gin.Context) {
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := services.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secured := os.Getenv("APP_ENV") == "production"

	c.SetCookie("accessToken", accessToken, 15*60, "/", "", secured, true)
	c.SetCookie("refreshToken", refreshToken, 7*24*60*60, "/", "", secured, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login success",
		"data": gin.H{
			"username": req.Username,
		},
	})
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	newAccessToken, newRefreshToken, err := services.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secured := os.Getenv("APP_ENV") == "production"

	c.SetCookie("accessToken", newAccessToken, 15*60, "/", "", secured, true)
	c.SetCookie("refreshToken", newRefreshToken, 7*24*60*60, "", "", secured, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "token refreshed"})
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	_ = services.Logout(refreshToken)

	secured := os.Getenv("APP_ENV") == "production"

	c.SetCookie("accessToken", "", -1, "/", "", secured, true)
	c.SetCookie("refreshToken", "", -1, "/", "", secured, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "logout success"})
}
