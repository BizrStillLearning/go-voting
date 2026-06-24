package main

import (
	"net/http"
	"time"

	"go-voting/internal/config"
	"go-voting/internal/controllers"
	"go-voting/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()

	// Pengaturan CORS untuk LAN
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./public/uploads")

	// Rute Bantuan (Agar tidak 404 jika dibuka langsung di browser)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "E-Voting API Backend Server is Running!",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/poll/:slug", controllers.GetPollPublic)
		api.GET("/poll/:slug/results", controllers.GetLiveResults)
		api.POST("/v/:id", controllers.SubmitVote)
		api.POST("/v/:id/verify", controllers.VerifyToken)

		api.POST("/admin/login", controllers.ProcessLogin)

		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			admin.POST("/poll/create", controllers.CreatePoll)
			admin.GET("/polls", controllers.GetAdminPolls)
			admin.PUT("/poll/:id/close", controllers.ClosePoll)
			admin.DELETE("/poll/:id/delete", controllers.DeletePoll)
			admin.GET("/poll/:id/export", controllers.ExportPollCSV)

			admin.PUT("/password", controllers.ChangePassword)
			admin.POST("/archive", controllers.ArchivePolls)
			admin.DELETE("/reset", controllers.ResetDatabase)
		}
	}

	r.Run(":8080")
}
