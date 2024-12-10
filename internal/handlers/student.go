package handlers

import (
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateStudentDetails updates a student's details
func UpdateStudentDetails(c *gin.Context) {
	// Extract the student ID from the URL parameter
	studentID := c.Param("id")

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
	if err := database.DB.Where("student_id = ?", studentID).First(&student).Error; err != nil {
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
