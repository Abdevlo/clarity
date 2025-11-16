package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/clarity/backend/config"
	"github.com/clarity/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.AuthConfig
}

func NewAuthService(db *gorm.DB, cfg *config.AuthConfig) *AuthService {
	return &AuthService{
		db:     db,
		config: cfg,
	}
}

// SendOTP generates and stores an OTP
func (as *AuthService) SendOTP(email string) (string, error) {
	otp := generateOTP(as.config.OTPLength)

	otpStore := models.OTPStore{
		ID:        uuid.New().String(),
		Email:     email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(time.Duration(as.config.OTPExpiry) * time.Second),
		CreatedAt: time.Now(),
	}

	// In production, send via email service
	log.Printf("OTP for %s: %s (expires in %d seconds)", email, otp, as.config.OTPExpiry)

	if err := as.db.Create(&otpStore).Error; err != nil {
		return "", fmt.Errorf("failed to store OTP: %w", err)
	}

	return otp, nil // In production, don't return OTP
}

// VerifyOTP validates the OTP and returns tokens
func (as *AuthService) VerifyOTP(email, otp string) (*models.User, string, string, error) {
	var otpStore models.OTPStore

	if err := as.db.Where("email = ? AND otp = ?", email, otp).First(&otpStore).Error; err != nil {
		return nil, "", "", fmt.Errorf("invalid OTP")
	}

	if time.Now().After(otpStore.ExpiresAt) {
		as.db.Delete(&otpStore)
		return nil, "", "", fmt.Errorf("OTP expired")
	}

	// Get or create user
	var user models.User
	if err := as.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				ID:        uuid.New().String(),
				Email:     email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := as.db.Create(&user).Error; err != nil {
				return nil, "", "", fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, "", "", fmt.Errorf("failed to fetch user: %w", err)
		}
	}

	// Generate tokens
	accessToken := generateToken(user.ID, 24*time.Hour)
	refreshToken := generateToken(user.ID, 7*24*time.Hour)

	// Delete used OTP
	as.db.Delete(&otpStore)

	return &user, accessToken, refreshToken, nil
}

// Helper functions
func generateOTP(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return fmt.Sprintf("%0*d", length, int64(bytes[0])%1000000)
}

func generateToken(userID string, duration time.Duration) string {
	token := hex.EncodeToString([]byte(userID + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)))
	log.Printf("Generated token for user %s", userID)
	return token
}

// RefreshToken validates refresh token and returns new access token
func (as *AuthService) RefreshToken(refreshToken string) (string, error) {
	// In production, implement proper JWT validation
	if refreshToken == "" {
		return "", fmt.Errorf("invalid refresh token")
	}

	accessToken := generateToken("user-id", 24*time.Hour)
	return accessToken, nil
}
