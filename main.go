package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Poll struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	Options     []Option  `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"options"`
}

type Option struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PollID    uint   `json:"poll_id"`
	Value     string `gorm:"type:varchar(255);not null" json:"value"`
	VoteCount int    `gorm:"default:0" json:"vote_count"`
}

type Vote struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PollID     uint      `gorm:"not null" json:"poll_id"`
	OptionID   uint      `gorm:"not null" json:"option_id"`
	Identifier string    `gorm:"type:varchar(255);not null" json:"identifier"`
	CreatedAt  time.Time `json:"created_at"`
}

var DB *gorm.DB

func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go-voting?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	DB.AutoMigrate(&Poll{}, &Option{}, &Vote{})
	fmt.Println("Database Connected & Migrated!")
}

func main() {
	InitDB()
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"multiply": func(a, b int) int { return a * b },
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
	})

	r.Static("/assets", "./assets")
	r.Static("/node_modules", "./node_modules")
	r.LoadHTMLGlob("templates/*")

	// Page: Home
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"Title": "GoVoting - Buat Polling"})
	})

	r.POST("/create-poll", func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		options := c.PostFormArray("options[]")
		pollSlug := slug.Make(title) + "-" + fmt.Sprint(time.Now().Unix())

		err := DB.Transaction(func(tx *gorm.DB) error {
			newPoll := Poll{Title: title, Description: description, Slug: pollSlug}
			if err := tx.Create(&newPoll).Error; err != nil {
				return err
			}

			for _, opt := range options {
				if opt == "" {
					continue
				}
				if err := tx.Create(&Option{PollID: newPoll.ID, Value: opt}).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Gagal simpan data")
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/v/"+pollSlug)
	})

	r.GET("/v/:slug", func(c *gin.Context) {
		var poll Poll
		if err := DB.Preload("Options").Where("slug = ?", c.Param("slug")).First(&poll).Error; err != nil {
			c.String(http.StatusNotFound, "Polling tidak ditemukan")
			return
		}
		c.HTML(http.StatusOK, "voting.html", gin.H{"Poll": poll})
	})

	r.POST("/vote/:poll_id", func(c *gin.Context) {
		pollID, _ := strconv.Atoi(c.Param("poll_id"))
		optionID, _ := strconv.Atoi(c.PostForm("option_id"))

		var poll Poll
		DB.First(&poll, pollID)

		DB.Model(&Option{}).Where("id = ? AND poll_id = ?", optionID, pollID).UpdateColumn("vote_count", gorm.Expr("vote_count + 1"))

		DB.Create(&Vote{
			PollID:     uint(pollID),
			OptionID:   uint(optionID),
			Identifier: c.ClientIP(),
		})

		c.Redirect(http.StatusSeeOther, "/v/"+poll.Slug+"/results")
	})

	r.GET("/v/:slug/results", func(c *gin.Context) {
		var poll Poll
		if err := DB.Preload("Options").Where("slug = ?", c.Param("slug")).First(&poll).Error; err != nil {
			c.String(http.StatusNotFound, "Polling tidak ditemukan")
			return
		}

		var totalVotes int
		for _, opt := range poll.Options {
			totalVotes += opt.VoteCount
		}

		c.HTML(http.StatusOK, "results.html", gin.H{
			"Poll":       poll,
			"TotalVotes": totalVotes,
		})
	})

	r.Run(":8080")
}
