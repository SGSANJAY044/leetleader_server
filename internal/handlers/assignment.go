package handlers

import (
	"net/http"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"time"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AssignQuestion(c *gin.Context) {
	println("IN")
	var input struct {
		StudentID   uint      `json:"student_id" binding:"required"`
		QuestionID  uint      `json:"question_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify student exists
	var student models.Student
	if err := database.DB.First(&student, input.StudentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Check if the question is already assigned to the student
	var existingAssignment models.Assignment
	if err := database.DB.Where("student_id = ? AND question_id = ?", input.StudentID, input.QuestionID).First(&existingAssignment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This question is already assigned to the student"})
		return
	}

	assignment := models.Assignment{
		StudentID:  input.StudentID,
		QuestionID: input.QuestionID,
		AssignedAt: time.Now(),
		Submitted:  false,
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Assignment created successfully", "assignment": assignment})
}

func AssignQuestions(c *gin.Context) {
	var input struct {
		StudentID   uint   `json:"student_id" binding:"required"`
		QuestionIDs []uint `json:"question_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify student exists
	var student models.Student
	if err := database.DB.First(&student, input.StudentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	successfulAssignments := []models.Assignment{}
	failedAssignments := []uint{}

	for _, questionID := range input.QuestionIDs {
		// Check if the question is already assigned to the student
		var existingAssignment models.Assignment
		if err := database.DB.Where("student_id = ? AND question_id = ?", input.StudentID, questionID).First(&existingAssignment).Error; err == nil {
			failedAssignments = append(failedAssignments, questionID)
			continue
		}

		assignment := models.Assignment{
			StudentID:  input.StudentID,
			QuestionID: questionID,
			AssignedAt: time.Now(),
			Submitted:  false,
		}

		if err := database.DB.Create(&assignment).Error; err != nil {
			failedAssignments = append(failedAssignments, questionID)
			continue
		}

		successfulAssignments = append(successfulAssignments, assignment)
	}

	response := gin.H{
		"message":               "Assignments processed",
		"successful_assignments": successfulAssignments,
		"failed_question_ids":   failedAssignments,
	}

	if len(successfulAssignments) == 0 {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(failedAssignments) == 0 {
		c.JSON(http.StatusCreated, response)
		return
	}

	c.JSON(http.StatusPartialContent, response)
}

// GetTodaysAssignments retrieves assignments for a student within the last 24 hours
func GetTodaysAssignments(c *gin.Context) {
	// Get student ID from URL parameter
	studentID := c.Param("student_id")

	// Convert string ID to uint
	id, err := strconv.ParseUint(studentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Calculate time 24 hours ago
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	// Query assignments within last 24 hours
	var assignments []models.Assignment
	if err := database.DB.Where("student_id = ? AND assigned_at >= ?", id, twentyFourHoursAgo).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	// If no assignments found
	if len(assignments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No assignments found for the last 24 hours",
			"assignments": []models.Assignment{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignments retrieved successfully",
		"assignments": assignments,
	})
}

