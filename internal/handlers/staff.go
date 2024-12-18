package handlers

import (
	"net/http"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// UpdateStaffDetails updates a staff's details
func UpdateStaffDetails(c *gin.Context) {
	// Extract the staff ID from the URL parameter
	Mail := c.Param("mail")

	// Define the input struct for update
	var input struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	// Bind the JSON request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the staff record from the database
	var staff models.Staff
	if err := database.DB.Where("mail = ?", Mail).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	// Update staff fields if provided
	if input.Name != "" {
		staff.Name = input.Name
	}
	if input.Phone != "" {
		staff.Phone = input.Phone
	}

	// Save the updated staff record to the database
	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff details updated successfully", "staff": staff})
}

func GetStaffDetails(c *gin.Context) {
	staffID := c.Param("id") // Get staff ID from route parameter

	var staff models.Staff
	err := database.DB.Where("staff_id = ?", staffID).First(&staff).Error

	// Handle the error based on the query result
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Respond with a "not found" message if the staff is not found
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Staff not found",
			})
		} else {
			// Respond with a generic error message if any other error occurs
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Database error",
			})
		}
		return
	}

	// Respond with the staff details on success
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   staff,
	})
}