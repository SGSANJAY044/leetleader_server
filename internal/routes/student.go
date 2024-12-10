package routes

import (
	"leetleader_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(r *gin.Engine) {
	student := r.Group("/students")
	{
		student.PUT("/:id", handlers.UpdateStudentDetails) // Update student details
	}
}
