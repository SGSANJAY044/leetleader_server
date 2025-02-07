package routes

import (
	"leetleader_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(r *gin.Engine) {
	student := r.Group("/students")
	{
		student.GET("/:roll", handlers.GetStudentDetails)
		student.PUT("/:mail", handlers.UpdateStudentDetails) 
		student.PUT("/update/solved/:username", handlers.UpdateProblemCount)
		student.PUT("/update/streak/:username", handlers.UpdateDailyStreak)
		student.GET("/class/:classID", handlers.GetStudentsByClass)
		student.GET("/dept/:departmentID", handlers.GetStudentsByDept)
		student.GET("/submissions/:username", handlers.GetStudentsSubmissions)
	}
}
