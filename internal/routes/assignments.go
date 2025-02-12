package routes

import (
	"leetleader_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AssignmentRoutes(r *gin.Engine) {
	assignment := r.Group("/assignment")
	{ 
		assignment.POST("/assign/todaytask", handlers.AssignQuestion)
		assignment.POST("/assign/todaytasks", handlers.AssignQuestions)
		assignment.GET("/todaytasks/:student_id", handlers.GetTodaysAssignments)
	}
}
