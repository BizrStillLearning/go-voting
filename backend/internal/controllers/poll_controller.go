package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go-voting/internal/config"
	"go-voting/internal/models"

	"github.com/gin-gonic/gin"
)

type OptionInput struct {
	Value       string `json:"value"`
	Vision      string `json:"vision"`
	Mission     string `json:"mission"`
	PhotoBase64 string `json:"photo"`
}

type PositionInput struct {
	Title   string        `json:"title"`
	Options []OptionInput `json:"options"`
}

type CreatePollInput struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	VoterCount  int             `json:"voter_count"`
	ExpiresAt   *time.Time      `json:"expires_at"`
	Positions   []PositionInput `json:"positions"`
}

func CreatePoll(c *gin.Context) {
	var input CreatePollInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data ditolak: " + err.Error()})
		return
	}

	slugString := fmt.Sprintf("pemilihan-%d", time.Now().Unix())

	poll := models.Poll{
		Title:       input.Title,
		Description: input.Description,
		Slug:        slugString,
		IsActive:    true,
	}

	if input.ExpiresAt != nil {
		poll.ExpiresAt = *input.ExpiresAt
	}

	if err := config.DB.Create(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat polling"})
		return
	}

	for _, posInput := range input.Positions {
		position := models.Position{
			PollID: poll.ID,
			Title:  posInput.Title,
		}
		config.DB.Create(&position)

		for _, optInput := range posInput.Options {
			option := models.Option{
				PositionID: position.ID,
				Value:      optInput.Value,
				Vision:     optInput.Vision,
				Mission:    optInput.Mission,
				PhotoURL:   optInput.PhotoBase64,
			}
			config.DB.Create(&option)
		}
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < input.VoterCount; i++ {
		tokenStr := fmt.Sprintf("VOTE-%X%X", rand.Intn(9999), rand.Intn(9999))
		config.DB.Create(&models.VoterToken{
			PollID:      poll.ID,
			TokenString: tokenStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pemilihan berhasil dibuat!", "slug": slugString})
}

func GetAdminPolls(c *gin.Context) {
	var polls []models.Poll
	if err := config.DB.Preload("Positions.Options").Preload("Tokens").Order("created_at desc").Find(&polls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data polling"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": polls})
}

func ClosePoll(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Model(&models.Poll{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menutup polling"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Polling berhasil ditutup"})
}

func DeletePoll(c *gin.Context) {
	id := c.Param("id")

	var poll models.Poll
	if err := config.DB.First(&poll, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data polling tidak ditemukan"})
		return
	}

	config.DB.Unscoped().Where("poll_id = ?", poll.ID).Delete(&models.VoterToken{})

	var positions []models.Position
	config.DB.Where("poll_id = ?", poll.ID).Find(&positions)
	for _, pos := range positions {
		config.DB.Unscoped().Where("position_id = ?", pos.ID).Delete(&models.Option{})
	}
	config.DB.Unscoped().Where("poll_id = ?", poll.ID).Delete(&models.Position{})

	if err := config.DB.Unscoped().Delete(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus dari database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Polling beserta seluruh data kandidat dan token berhasil dibumihanguskan!"})
}
