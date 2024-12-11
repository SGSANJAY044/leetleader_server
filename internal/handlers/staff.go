package handlers

import (
	"net/http"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"

	"github.com/gin-gonic/gin"
)

// UpdateStaffDetails updates a staff's details
func UpdateStaffDetails(c *gin.Context) {
	// Extract the staff ID from the URL parameter
	staffID := c.Param("id")

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
	if err := database.DB.First(&staff, staffID).Error; err != nil {
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
