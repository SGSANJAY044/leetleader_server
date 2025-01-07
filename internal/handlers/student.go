package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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


func UpdateProblemCount(c *gin.Context) {
	// Extract the username from the request parameters
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username is required",
		})
		return
	}

	// Call the LeetCode API to fetch student details
	apiURL := fmt.Sprintf("http://localhost:3000/%s/solved", username) // Replace with actual API URL
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch data from LeetCode API",
		})
		return
	}
	defer resp.Body.Close()

	// Parse the response from the LeetCode API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to read response from LeetCode API",
		})
		return
	}

	// Define a structure to parse the API response
	type LeetCodeData struct {
		SolvedEasy  int    `json:"easySolved"`
		SolvedMedium int   `json:"mediumSolved"`
		SolvedHard  int    `json:"hardSolved"`
	}

	var leetCodeData LeetCodeData
	if err := json.Unmarshal(body, &leetCodeData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to parse response from LeetCode API",
		})
		return
	}

	var student models.Student
	if err := database.DB.Where("username = ?", username).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Student not found in database",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to query database",
			})
		}
		return
	}

	// Update student fields
	student.SolvedEasy = leetCodeData.SolvedEasy
	student.SolvedMedium = leetCodeData.SolvedMedium
	student.SolvedHard = leetCodeData.SolvedHard

	// Save the updated record
	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update student details",
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Student details updated successfully",
		"data":    student,
	})
}