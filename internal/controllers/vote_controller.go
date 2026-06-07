package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go-voting/internal/config"
	"go-voting/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowVotingPage(c *gin.Context) {
	slug := c.Param("slug")
	var poll models.Poll

	if err := config.DB.Preload("Options").Where("slug = ?", slug).First(&poll).Error; err != nil {
		c.String(http.StatusNotFound, "Halaman polling tidak ditemukan.")
		return
	}

	isClosed := !poll.IsActive || time.Now().After(poll.ExpiresAt)

	c.HTML(http.StatusOK, "vote.html", gin.H{
		"Title":    poll.Title,
		"Poll":     poll,
		"IsClosed": isClosed,
	})
}

func ProcessVote(c *gin.Context) {
	pollID := c.Param("poll_id")
	optionID := c.PostForm("option_id")
	tokenString := c.PostForm("token")

	var poll models.Poll
	config.DB.First(&poll, pollID)

	if !poll.IsActive || time.Now().After(poll.ExpiresAt) {
		errMsg := url.QueryEscape("Voting telah ditutup atau sudah kadaluwarsa!")
		c.Redirect(http.StatusSeeOther, "/v/"+poll.Slug+"?error="+errMsg)
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var token models.VoterToken

		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("token_string = ? AND poll_id = ? AND is_used = ?", tokenString, pollID, false).
			First(&token).Error; err != nil {
			return err
		}

		if err := tx.Model(&token).Update("is_used", true).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Option{}).Where("id = ?", optionID).UpdateColumn("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
			return err
		}

		newVote := models.Vote{
			PollID:   token.PollID,
			OptionID: uint(mustParseUint(optionID)),
			TokenID:  token.ID,
		}
		if err := tx.Create(&newVote).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		errMsg := url.QueryEscape("Token tidak valid, salah, atau sudah pernah digunakan!")
		c.Redirect(http.StatusSeeOther, "/v/"+poll.Slug+"?error="+errMsg)
		return
	}

	successMsg := url.QueryEscape("Suara Anda berhasil direkam!")
	c.Redirect(http.StatusSeeOther, "/v/"+poll.Slug+"/results?success="+successMsg)
}

func ShowResults(c *gin.Context) {
	slug := c.Param("slug")
	var poll models.Poll

	if err := config.DB.Preload("Options").Where("slug = ?", slug).First(&poll).Error; err != nil {
		c.String(http.StatusNotFound, "Polling tidak ditemukan")
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Title": "Hasil: " + poll.Title,
		"Poll":  poll,
	})
}

func GetPollResultsAPI(c *gin.Context) {
	slug := c.Param("slug")
	var poll models.Poll

	if err := config.DB.Preload("Options").Where("slug = ?", slug).First(&poll).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	var totalVotes int
	for _, opt := range poll.Options {
		totalVotes += opt.VoteCount
	}

	c.JSON(http.StatusOK, gin.H{
		"total_votes": totalVotes,
		"options":     poll.Options,
	})
}

func mustParseUint(s string) int {
	var res int
	fmt.Sscanf(s, "%d", &res)
	return res
}
