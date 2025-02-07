package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"net/http"
	"time"
	"strconv"
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

func UpdateDailyStreak(c *gin.Context) {
	type Submission struct {
	Title         string `json:"title"`         
	TitleSlug     string `json:"titleSlug"`     
	Timestamp     string `json:"timestamp"`     
	StatusDisplay string `json:"statusDisplay"` 
	Lang          string `json:"lang"`          
   }
   
	type SubmissionsResponse struct {
		Count       int          `json:"count"`
		Submissions []Submission `json:"submission"`
	}
	
	type Student struct {
		StudentID uint   `gorm:"primaryKey"`
		Username  string `gorm:"size:50;unique"`
		Streak    int    `gorm:"not null"`
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username is required",
		})
		return
	}

	apiURL := fmt.Sprintf("http://localhost:3000/%s/acSubmission", username) // Replace with actual API URL
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch data from LeetCode API",
		})
		return
	}
	defer resp.Body.Close()

	// Decode the API response
	var apiResponse SubmissionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Failed to decode API response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to process API response",
		})
		return
	}

	// Check if submissions exist
	if apiResponse.Count == 0 || len(apiResponse.Submissions) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No submissions found for today",
		})
		return
	}

	// Check if the first submission timestamp matches today's date
	firstSubmission := apiResponse.Submissions[0]
	timestampInt, err := strconv.ParseInt(firstSubmission.Timestamp, 10, 64)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": "Failed to parse submission timestamp",
            })
            return
        }
	submissionDate := time.Unix(timestampInt, 0).UTC()
	currentDate := time.Now().UTC()

	if submissionDate.Year() == currentDate.Year() && submissionDate.YearDay() == currentDate.YearDay() {
		// Increment streak if dates match
		result := database.DB.Model(&Student{}).Where("username = ?", username).Update("streak", gorm.Expr("streak + 1"))
		if result.Error != nil {
			fmt.Printf("Failed to update streak: %v\n", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to update streak",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Streak updated successfully",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No submissions for today; streak not updated",
		})
	}
}

func GetStudentsByClass(c *gin.Context) {
    classID := c.Param("classID") // Get ClassID from route parameter

    var students []models.Student
    if err := database.DB.Where("class_id = ?", classID).Find(&students).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to fetch students",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "students": students,
    })
}

func GetStudentsByDept(c *gin.Context) {
    departmentID := c.Param("departmentID")

    var students []models.Student
    if err := database.DB.Where("department_id = ?", departmentID).Find(&students).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to fetch students",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "students": students,
    })
}

func GetStudentsSubmissions(c *gin.Context) {
	username := c.Param("username")
	apiURL := fmt.Sprintf("http://localhost:3000/%s/submission", username)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch data from LeetCode API",
		})
		return
	}
	defer resp.Body.Close()

	type submission struct {
		Title         string `json:"title"`
		TitleSlug     string `json:"titleSlug"`
		Timestamp     string `json:"timestamp"`
		StatusDisplay string `json:"statusDisplay"`
		Lang          string `json:"lang"`
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error", 
			"message": "Failed to read response body",
		})
		return
	}

	type SubmissionsResponse struct {
        Count       int          `json:"count"`
        Submissions []submission `json:"submission"`
    }

    // Unmarshal into the wrapper structure
    var response SubmissionsResponse
    if err := json.Unmarshal(body, &response); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to unmarshal response",
        })
        return
    }

	

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"submission": response.Submissions,
	})
}