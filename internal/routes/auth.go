package routes

import (
	"leetleader_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		// Student routes
		auth.POST("/student/signup", handlers.SignupStudent)
		auth.POST("/student/login", handlers.LoginStudent)

		// Staff routes
		auth.POST("/staff/signup", handlers.SignupStaff)
		auth.POST("/staff/login", handlers.LoginStaff)
	}
}
