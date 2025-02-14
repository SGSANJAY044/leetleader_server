package models
import "time" 
// Student model
type Student struct {
	StudentID   uint   `gorm:"primaryKey"`
	Streak      int    `gorm:"not null"`
	SolvedEasy  int    `gorm:"default:0"`
	SolvedMedium int   `gorm:"default:0"`
	SolvedHard  int    `gorm:"default:0"`
	Ranking     int    `gorm:"not null"`
	Name        string `gorm:"size:100;not null"`
	ClassID     uint   `gorm:"not null"`
	Roll        string `gorm:"size:20;not null;"`
	DepartmentID uint  `gorm:"not null"`
	Phone       string `gorm:"size:15"`
	Mail        string `gorm:"size:100;unique"`
	Username    string `gorm:"size:50;"`
	Password    string `gorm:"not null"` // Secure password storage
}

// Staff model
type Staff struct {
	StaffID      uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100;not null"`
	ClassID      uint   `gorm:"not null"`
	DepartmentID uint   `gorm:"not null"`
	Phone        string `gorm:"size:15"`
	Mail         string `gorm:"size:100;unique"`
	IsHOD        bool   `gorm:"default:false"`
	Password     string `gorm:"not null"` // Secure password storage
}
// Class model
type Class struct {
	ClassID      uint   `gorm:"primaryKey"`
	DepartmentID uint   `gorm:"not null"`
	Name         string `gorm:"size:50;not null"`
	StudentCount int    `gorm:"not null"`
	Staffs       []Staff `gorm:"foreignKey:ClassID"`
}

// Department model
type Department struct {
	DepartmentID uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:50;not null"`
	StudentCount int    `gorm:"not null"`
	StaffCount   int    `gorm:"not null"`
}

// Question model
type Question struct {
	QuestionID    uint   `gorm:"primaryKey"`
	QuestionTitle string `gorm:"size:255;not null"`
	TitleSlug     string `gorm:"size:255;not null;unique"`
	Difficulty    string `gorm:"size:50;not null"`
	Question      string `gorm:"type:text;not null"`
}
// Assignment model
type Assignment struct {
    AssignmentID uint      `gorm:"primaryKey"`
    StudentID    uint      `gorm:"not null"`  
    TitleSlug    string    `gorm:"not null"`  
    AssignedAt   time.Time `gorm:"not null"`  
    Submitted    bool      `gorm:"default:false"`
    SubmittedAt  *time.Time 
}

// FriendsQuestions model
type FriendsQuestions struct {
	StudentID   uint      `gorm:"primaryKey"`
	TitleSlug   string    `gorm:"not null"`
}
