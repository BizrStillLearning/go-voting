package controllers

import (
	"go-voting/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ProcessLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tidak valid"})
		return
	}

	var admin Admin
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password salah!"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password salah!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   "SUPERADMIN-SECRET-TOKEN-12345",
	})
}

func ChangePassword(c *gin.Context) {
	var input PasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var admin Admin
	if err := config.DB.First(&admin, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan di database"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sandi saat ini salah!"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi sandi baru"})
		return
	}

	if err := config.DB.Model(&admin).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sandi ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sandi berhasil diperbarui!"})
}
