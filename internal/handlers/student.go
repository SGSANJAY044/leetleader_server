package handlers

import (
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"net/http"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// UpdateStudentDetails updates a student's details
func UpdateStudentDetails(c *gin.Context) {
	// Extract the student ID from the URL parameter
	Mail := c.Param("mail")

	var input struct {
		Name         string `json:"name"`
		ClassID      uint   `json:"class_id"`
		Roll         string `json:"roll"`
		DepartmentID uint   `json:"department_id"`
		Phone        string `json:"phone"`
		Username     string `json:"username"`
	}

	// Bind JSON input to the `input` struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the student from the database
	var student models.Student
	if err := database.DB.Where("mail = ?", Mail).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Update student fields if provided
	if input.Name != "" {
		student.Name = input.Name
	}
	if input.ClassID != 0 {
		student.ClassID = input.ClassID
	}
	if input.Roll != "" {
		student.Roll = input.Roll
	}
	if input.DepartmentID != 0 {
		student.DepartmentID = input.DepartmentID
	}
	if input.Phone != "" {
		student.Phone = input.Phone
	}
	if input.Username != "" {
		student.Username = input.Username
	}

	// Save the updated student details
	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student details updated successfully", "student": student})
}

func GetStudentDetails(c *gin.Context) {
	roll := c.Param("roll") // Get roll from route parameter

	var student models.Student
	err := database.DB.Where("roll = ?", roll).First(&student).Error

	// Handle the error based on the query result
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Respond with a "not found" message if the student is not found
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Student not found",
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

	// Respond with the student details on success
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   student,
	})
}
