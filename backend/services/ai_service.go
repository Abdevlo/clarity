package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	vision "cloud.google.com/go/vision/v2"
	"github.com/clarity/backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrescriptionData struct {
	Medication string `json:"medication"`
	Dosage     string `json:"dosage"`
	Frequency  string `json:"frequency"`
	Duration   string `json:"duration"`
	Indication string `json:"indication"`
	Warnings   string `json:"warnings,omitempty"`
	Refills    string `json:"refills,omitempty"`
}

type AIService struct {
	db *gorm.DB
}

func NewAIService(db *gorm.DB) *AIService {
	return &AIService{db: db}
}

// ScanPrescription extracts data from prescription image
func (as *AIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
	// Placeholder for AI prescription scanning
	// In production, integrate with OpenAI Vision API or similar

	log.Printf("Scanning prescription for user %s", userID)

	// Mock extracted data
	extractedData := map[string]string{
		"medication": "Aspirin",
		"dosage":     "500mg",
		"frequency":  "Twice daily",
		"duration":   "7 days",
		"indication": "Headache/Pain relief",
	}

	return extractedData, nil
}

// SummarizeHealth generates a health summary
func (as *AIService) SummarizeHealth(userID string, days int) (string, []string, string, error) {
	// Fetch user's recent health records
	var records []models.HealthRecord
	startDate := time.Now().AddDate(0, 0, -days)

	if err := as.db.Where("user_id = ? AND created_at > ?", userID, startDate).
		Find(&records).Error; err != nil {
		return "", nil, "", fmt.Errorf("failed to fetch records: %w", err)
	}

	log.Printf("Summarizing %d health records for user %s", len(records), userID)

	// Mock summarization (in production, use AI model)
	summary := fmt.Sprintf("Health Summary for last %d days: %d records found.", days, len(records))

	keyFindings := []string{
		"Overall health status: Good",
		"Recent medications: None critical",
		"Recommended actions: Regular check-up",
	}

	recommendations := "Stay hydrated, maintain regular exercise, and schedule a check-up next month."

	return summary, keyFindings, recommendations, nil
}

// DoctorChat handles conversation with AI doctor
func (as *AIService) DoctorChat(userID, conversationID, message string) (string, error) {
	// Placeholder for AI-powered doctor chat
	// In production, integrate with LLM API

	log.Printf("Doctor chat for user %s: %s", userID, message)

	// Mock AI response
	response := fmt.Sprintf("AI Doctor: I've noted your concern about '%s'. Please provide more details about your symptoms.", message)

	// Store conversation
	conversation := models.DoctorConversation{
		ID:             uuid.New().String(),
		UserID:         userID,
		ConversationID: conversationID,
		Message:        message,
		Response:       response,
		IsAI:           true,
		CreatedAt:      time.Now(),
	}

	if err := as.db.Create(&conversation).Error; err != nil {
		return "", fmt.Errorf("failed to store conversation: %w", err)
	}

	return response, nil
}

// GetConversationHistory retrieves chat history
func (as *AIService) GetConversationHistory(conversationID string) ([]models.DoctorConversation, error) {
	var conversations []models.DoctorConversation
	if err := as.db.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch conversations: %w", err)
	}
	return conversations, nil
}

// Helper to parse prescription metadata
func parsePrescriptionMetadata(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	err := json.Unmarshal(data, &result)
	return result, err
}

func extractDataFromScanWithVisionAPI(imageData []byte) (*PrescriptionData, error) {
	ctx := context.Background()

	// Step 2: Create Vision client
	// This requires Google Cloud credentials (JSON key file)
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vision client: %w", err)
	}
	defer client.Close()

	// Step 3: Create image from bytes
	// Vision can work with image bytes, URLs, or cloud storage paths
	image := vision.NewImageFromBytes(imageData)

	// Step 4: Detect text in image
	// This performs OCR on the prescription image
	annotations, err := client.DetectTexts(ctx, image, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detect text: %w", err)
	}

	// Step 5: Extract and parse text
	fullText := ""
	for _, annotation := range annotations {
		fullText += annotation.GetDescription() + "\n"
	}

	// Step 6: Parse extracted text into structured data
	// This is simplified - in real use, you'd use more sophisticated parsing
	prescription := parsePrescriptionText(fullText)

	return prescription, nil
}

func parsePrescriptionText(text string) *PrescriptionData {
	// This is a simple example - real parsing would be more complex
	// You could use regex patterns or a more advanced NLP model

	prescription := &PrescriptionData{
		Medication: "Extract from text",
		Dosage:     "Extract from text",
		Frequency:  "Extract from text",
	}

	// In production, use pattern matching or AI to parse structured data
	return prescription
}
