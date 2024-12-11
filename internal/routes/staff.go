package routes

import (
	"leetleader_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(r *gin.Engine) {
	staff := r.Group("/staffs")
	{
		staff.PUT("/:id", handlers.UpdateStaffDetails) // Update staff details
	}
}
