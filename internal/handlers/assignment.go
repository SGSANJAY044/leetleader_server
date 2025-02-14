package handlers

import (
	"net/http"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"time"
	"github.com/gin-gonic/gin"
	"strconv"
	"encoding/json"
    "fmt"

)

func AssignQuestion(c *gin.Context) {
	println("IN")
	var input struct {
		StudentID uint   `json:"student_id" binding:"required"`
		TitleSlug string `json:"title_slug" binding:"required"`
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
	if err := database.DB.Where("student_id = ? AND title_slug = ?", input.StudentID, input.TitleSlug).First(&existingAssignment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This question is already assigned to the student"})
		return
	}

	assignment := models.Assignment{
		StudentID:  input.StudentID,
		TitleSlug:  input.TitleSlug,
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
		StudentID  uint     `json:"student_id" binding:"required"`
		TitleSlugs []string `json:"title_slugs" binding:"required"`
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
	failedAssignments := []string{}

	for _, titleSlug := range input.TitleSlugs {
		// Check if the question is already assigned to the student
		var existingAssignment models.Assignment
		if err := database.DB.Where("student_id = ? AND title_slug = ?", input.StudentID, titleSlug).First(&existingAssignment).Error; err == nil {
			failedAssignments = append(failedAssignments, titleSlug)
			continue
		}

		assignment := models.Assignment{
			StudentID:  input.StudentID,
			TitleSlug:  titleSlug,
			AssignedAt: time.Now(),
			Submitted:  false,
		}

		if err := database.DB.Create(&assignment).Error; err != nil {
			failedAssignments = append(failedAssignments, titleSlug)
			continue
		}

		successfulAssignments = append(successfulAssignments, assignment)
	}

	response := gin.H{
		"message":               "Assignments processed",
		"successful_assignments": successfulAssignments,
		"failed_title_slugs":    failedAssignments,
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


// GetTodayAssignmentQuestions retrieves the full question details for today's assignments
func GetTodayAssignmentQuestions(c *gin.Context) {
	studentID := c.Param("student_id")
	id, err := strconv.ParseUint(studentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	var assignments []models.Assignment
	if err := database.DB.Where("student_id = ? AND assigned_at >= ?", id, twentyFourHoursAgo).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	if len(assignments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No assignments found for the last 24 hours", "questions": []models.Question{}})
		return
	}

	var titleSlugs []string
	for _, assignment := range assignments {
		titleSlugs = append(titleSlugs, assignment.TitleSlug)
	}

	var questions []models.Question
	if err := database.DB.Where("title_slug IN ?", titleSlugs).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	// Fetch missing questions from API
	existingTitleSlugs := make(map[string]bool)
	for _, q := range questions {
		existingTitleSlugs[q.TitleSlug] = true
	}

	var missingQuestions []models.Question
	for _, titleSlug := range titleSlugs {
		if !existingTitleSlugs[titleSlug] {
			apiURL := fmt.Sprintf("http://localhost:3000/select?titleSlug=%s", titleSlug)
			resp, err := http.Get(apiURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from API"})
				return
			}
			defer resp.Body.Close()

			var apiResponse struct {
				Link                string   `json:"link"`
				QuestionId         string   `json:"questionId"`
				QuestionFrontendId string   `json:"questionFrontendId"`
				QuestionTitle      string   `json:"questionTitle"`
				TitleSlug          string   `json:"titleSlug"`
				Difficulty         string   `json:"difficulty"`
				Question           string   `json:"question"`
				TopicTags         []struct {
					Name            string `json:"name"`
					Slug            string `json:"slug"`
					TranslatedName  string `json:"translatedName"`
				} `json:"topicTags"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
				return
			}

			questionID, _ := strconv.ParseUint(apiResponse.QuestionId, 10, 32)
			newQuestion := models.Question{
				QuestionID:    uint(questionID),
				QuestionTitle: apiResponse.QuestionTitle,
				TitleSlug:     apiResponse.TitleSlug,
				Difficulty:    apiResponse.Difficulty,
				Question:      apiResponse.Question,
			}

			missingQuestions = append(missingQuestions, newQuestion)
		}
	}

	// Store new questions in DB
	if len(missingQuestions) > 0 {
		if err := database.DB.Create(&missingQuestions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save questions to DB"})
			return
		}
		questions = append(questions, missingQuestions...)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Assignment questions retrieved successfully",
		"questions": questions,
	})
}

