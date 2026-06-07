package main

import (
	"fmt"

	"go-voting/internal/config"
	"go-voting/internal/controllers"
	"go-voting/internal/middlewares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()

	store := cookie.NewStore([]byte("rahasia-super-aman"))
	r.Use(sessions.Sessions("govoting_session", store))

	r.Static("/assets", "./assets")
	r.Static("/node_modules", "./node_modules")
	r.LoadHTMLGlob("templates/admin/*")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/login", controllers.ShowLogin)
	r.POST("/login", controllers.ProcessLogin)
	r.GET("/logout", controllers.Logout)

	r.GET("/v/:slug", controllers.ShowVotingPage)
	r.POST("/vote/:poll_id", controllers.ProcessVote)

	r.GET("/v/:slug/results", controllers.ShowResults)
	r.GET("/api/poll/:slug/results", controllers.GetPollResultsAPI)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middlewares.AuthRequired())
	{
		adminRoutes.GET("/dashboard", controllers.Dashboard)
		adminRoutes.GET("/poll/create", controllers.ShowCreatePoll)
		adminRoutes.POST("/poll/create", controllers.ProcessCreatePoll)

		adminRoutes.POST("/poll/:id/close", controllers.ClosePoll)
		adminRoutes.POST("/poll/:id/delete", controllers.DeletePoll)
	}

	fmt.Println("Server backend berjalan di http://localhost:8080")
	r.Run(":8080")
}
