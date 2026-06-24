package controllers

import (
	"go-voting/internal/config"
	"go-voting/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VoteItem struct {
	PositionID uint `json:"position_id"`
	OptionID   uint `json:"option_id"`
}

type SubmitVoteInput struct {
	Token string     `json:"token"`
	Votes []VoteItem `json:"votes"`
}

type VerifyTokenInput struct {
	Token string `json:"token"`
}

func SubmitVote(c *gin.Context) {
	pollID := c.Param("id")
	var input SubmitVoteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var poll models.Poll
	if err := config.DB.First(&poll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pemilihan tidak ditemukan"})
		return
	}

	if !poll.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Pemilihan telah ditutup"})
		return
	}

	tx := config.DB.Begin()

	var token models.VoterToken

	cleanToken := strings.ToUpper(strings.TrimSpace(input.Token))

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("poll_id = ? AND token_string = ?", poll.ID, cleanToken).First(&token).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau salah"})
		return
	}

	if token.IsUsed {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "Token sudah pernah digunakan!"})
		return
	}

	for _, voteItem := range input.Votes {
		if err := tx.Model(&models.Option{}).Where("id = ? AND position_id = ?", voteItem.OptionID, voteItem.PositionID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merekam suara"})
			return
		}
	}

	if err := tx.Model(&token).Update("is_used", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghanguskan token"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Suara berhasil direkam!"})
}

func GetPollPublic(c *gin.Context) {
	slug := c.Param("slug")
	var poll models.Poll

	if err := config.DB.Preload("Positions.Options").Where("slug = ?", slug).First(&poll).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Polling tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"poll": poll, "is_closed": !poll.IsActive})
}

func GetLiveResults(c *gin.Context) {
	slug := c.Param("slug")
	var poll models.Poll

	if err := config.DB.Preload("Positions.Options").Where("slug = ?", slug).First(&poll).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Polling tidak ditemukan"})
		return
	}

	var usedTokens int64
	config.DB.Model(&models.VoterToken{}).Where("poll_id = ? AND is_used = ?", poll.ID, true).Count(&usedTokens)

	c.JSON(http.StatusOK, gin.H{
		"poll":        poll,
		"total_votes": usedTokens,
	})
}

func VerifyToken(c *gin.Context) {
	pollID := c.Param("id")
	var input VerifyTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	cleanToken := strings.ToUpper(strings.TrimSpace(input.Token))

	var token models.VoterToken
	if err := config.DB.Where("poll_id = ? AND token_string = ?", pollID, cleanToken).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau salah!"})
		return
	}

	if token.IsUsed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak: Token ini sudah pernah digunakan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token aman!"})
}
