package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ProcessLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tidak valid"})
		return
	}

	if input.Username == "superadmin" && input.Password == "superadmin" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   "SUPERADMIN-SECRET-TOKEN-12345",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password salah!"})
}
