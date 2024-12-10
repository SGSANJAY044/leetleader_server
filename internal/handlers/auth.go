package handlers

import (
	"net/http"
	"leetleader_server/internal/database"
	"leetleader_server/internal/models"
	"leetleader_server/internal/utils"

	"github.com/gin-gonic/gin"
)

// SignupStudent handles student registration
func SignupStudent(c *gin.Context) {
	var input struct {
		Mail     string `json:"mail" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create student
	student := models.Student{
		Mail:     input.Mail,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Student registered successfully"})
}

// LoginStudent handles student login
func LoginStudent(c *gin.Context) {
	var input struct {
		Mail string `json:"mail" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find student by mail
	var student models.Student
	if err := database.DB.Where("mail = ?", input.Mail).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Mail"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(input.Password, student.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(student.StudentID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SignupStaff handles staff registration
func SignupStaff(c *gin.Context) {
	var input struct {
		ClassID  uint   `json:"class_id" binding:"required"`
		DepartmentID uint   `json:"department_id" binding:"required"`
		Mail     string `json:"mail" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create staff
	staff := models.Staff{
		Mail:     input.Mail,
		Password: hashedPassword,
		ClassID:  input.ClassID,
		DepartmentID: input.DepartmentID,
	}

	if err := database.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Staff registered successfully"})
}

// LoginStaff handles staff login
func LoginStaff(c *gin.Context) {
	var input struct {
		Mail string `json:"mail" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find staff by mail
	var staff models.Staff
	if err := database.DB.Where("mail = ?", input.Mail).First(&staff).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(input.Password, staff.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(staff.StaffID, "staff")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
