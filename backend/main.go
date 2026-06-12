package main

import (
	"time"

	"go-voting/internal/config"
	"go-voting/internal/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./public/uploads")

	api := r.Group("/api")
	{
		api.GET("/poll/:slug", controllers.GetPollPublic)
		api.GET("/poll/:slug/results", controllers.GetLiveResults)
		api.POST("/vote/:id", controllers.SubmitVote)
		api.POST("/vote/:id/verify", controllers.VerifyToken)

		admin := api.Group("/admin")
		{
			admin.POST("/login", controllers.ProcessLogin)
			admin.POST("/poll/create", controllers.CreatePoll)
			admin.GET("/polls", controllers.GetAdminPolls)
			admin.PUT("/poll/:id/close", controllers.ClosePoll)
			admin.DELETE("/poll/:id/delete", controllers.DeletePoll)

			admin.PUT("/password", controllers.ChangePassword)
			admin.POST("/archive", controllers.ArchivePolls)
			admin.DELETE("/reset", controllers.ResetDatabase)
		}
	}

	r.Run(":8080")
}
