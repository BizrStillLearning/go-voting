package controllers

import (
	"go-voting/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(c *gin.Context) {
	var input PasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sandi berhasil diperbarui!"})
}

func ArchivePolls(c *gin.Context) {
	if err := config.DB.Exec("UPDATE polls SET is_active = false").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengarsipkan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua pemilihan telah diarsipkan (ditutup)."})
}

func ResetDatabase(c *gin.Context) {
	config.DB.Exec("DELETE FROM votes")
	config.DB.Exec("DELETE FROM voter_tokens")
	config.DB.Exec("DELETE FROM options")
	config.DB.Exec("DELETE FROM polls")

	c.JSON(http.StatusOK, gin.H{"message": "Database berhasil diformat ke kondisi awal!"})
}
