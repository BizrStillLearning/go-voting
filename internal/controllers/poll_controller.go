package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go-voting/internal/config"
	"go-voting/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func Dashboard(c *gin.Context) {
	var polls []models.Poll
	config.DB.Preload("Options").Preload("Tokens").Find(&polls)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title": "Dashboard Admin",
		"Polls": polls,
	})
}

func ShowCreatePoll(c *gin.Context) {
	c.HTML(http.StatusOK, "create_poll.html", gin.H{"Title": "Buat Polling Baru"})
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return "VOTE-" + string(b)
}

func ViewTokens(c *gin.Context) {
	pollID := c.Param("id")
	var poll models.Poll

	if err := config.DB.Preload("Tokens").First(&poll, pollID).Error; err != nil {
		c.String(http.StatusNotFound, "Data polling tidak ditemukan")
		return
	}

	c.HTML(http.StatusOK, "tokens.html", gin.H{
		"Title": "Daftar Token Akses",
		"Poll":  poll,
	})
}

func ProcessCreatePoll(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	options := c.PostFormArray("options[]")
	voterCountStr := c.PostForm("voter_count")
	expiresAtStr := c.PostForm("expires_at")

	voterCount, _ := strconv.Atoi(voterCountStr)
	pollSlug := slug.Make(title) + "-" + fmt.Sprint(time.Now().Unix())

	expiresAt, err := time.ParseInLocation("2006-01-02T15:04", expiresAtStr, time.Local)
	if err != nil {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		newPoll := models.Poll{
			Title:       title,
			Description: description,
			Slug:        pollSlug,
			IsActive:    true, // Otomatis aktif saat dibuat
			ExpiresAt:   expiresAt,
		}
		if err := tx.Create(&newPoll).Error; err != nil {
			return err
		}

		for _, opt := range options {
			if opt == "" {
				continue
			}
			if err := tx.Create(&models.Option{PollID: newPoll.ID, Value: opt}).Error; err != nil {
				return err
			}
		}

		for i := 0; i < voterCount; i++ {
			tokenStr := generateRandomString(6)
			if err := tx.Create(&models.VoterToken{PollID: newPoll.ID, TokenString: tokenStr}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal memproses pembuatan polling")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func ClosePoll(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Model(&models.Poll{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		c.String(http.StatusInternalServerError, "Gagal menutup polling")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func DeletePoll(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Poll{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Gagal menghapus polling")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}
