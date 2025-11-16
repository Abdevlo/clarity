package models

import "time"

// User represents a user in the system
type User struct {
	ID            string    `gorm:"primaryKey"`
	Email         string    `gorm:"uniqueIndex"`
	Name          string
	DateOfBirth   string
	Gender        string
	BloodType     string
	PasswordHash  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OTPStore stores OTP data temporarily
type OTPStore struct {
	ID        string    `gorm:"primaryKey"`
	Email     string    `gorm:"index"`
	OTP       string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// HealthRecord stores health information
type HealthRecord struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"index"`
	RecordType  string    // prescription, appointment, lab_result, symptom
	Title       string
	Description string
	Metadata    string `gorm:"type:json"` // JSON string for flexibility
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DoctorConversation stores chat history
type DoctorConversation struct {
	ID             string    `gorm:"primaryKey"`
	UserID         string    `gorm:"index"`
	ConversationID string    `gorm:"index"`
	Message        string
	Response       string
	IsAI           bool
	CreatedAt      time.Time
}

// Token for JWT tokens
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
