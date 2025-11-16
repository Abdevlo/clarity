package examples

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

// ExampleOpenAIIntegration shows how to integrate OpenAI API for prescription scanning
// This example uses ChatGPT 4 Vision model to analyze prescription images

// PrescriptionData represents extracted prescription information
type PrescriptionData struct {
	Medication  string `json:"medication"`
	Dosage      string `json:"dosage"`
	Frequency   string `json:"frequency"`
	Duration    string `json:"duration"`
	Indication  string `json:"indication"`
	Warnings    string `json:"warnings,omitempty"`
	Refills     string `json:"refills,omitempty"`
}

// Example 1: Scan Prescription Using Vision API
// This function takes a prescription image and extracts medication information
func ScanPrescriptionWithOpenAI(imageData []byte) (*PrescriptionData, error) {
	// Step 1: Get API key from environment
	// Set this before running: export OPENAI_API_KEY="sk-xxxxx"
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	// Step 2: Create OpenAI client
	client := openai.NewClient(apiKey)

	// Step 3: Convert image bytes to base64
	// OpenAI requires images to be base64 encoded
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// Step 4: Create the request with image
	// Note: This uses GPT-4 Vision which can analyze images
	request := openai.ChatCompletionRequest{
		Model: openai.GPT4VisionPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					// First, send the image
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ImageURL{
							// Use data URL format for inline image
							URL: fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage),
						},
					},
					// Then ask the question
					{
						Type: openai.ChatMessagePartTypeText,
						Text: `Please analyze this prescription image and extract the following information in JSON format:
{
  "medication": "the medication name",
  "dosage": "dose amount and unit",
  "frequency": "how often to take (e.g., twice daily)",
  "duration": "how long to take the medication",
  "indication": "reason for the prescription",
  "warnings": "any warnings or contraindications",
  "refills": "number of refills allowed"
}

Return ONLY the JSON object, no other text.`,
					},
				},
			},
		},
		MaxTokens: 500,
	}

	// Step 5: Call OpenAI API
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	// Step 6: Extract the response
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := response.Choices[0].Message.Content

	// Step 7: Parse JSON response
	var prescription PrescriptionData
	err = json.Unmarshal([]byte(responseText), &prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &prescription, nil
}

// Example 2: Health Summarization Using ChatGPT
// This function takes health records and generates a summary using ChatGPT
type HealthRecord struct {
	Type        string
	Title       string
	Description string
	Date        string
}

func SummarizeHealthWithOpenAI(records []HealthRecord) (string, []string, string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", nil, "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)

	// Build a formatted string of health records
	recordsText := "Health Records:\n"
	for _, record := range records {
		recordsText += fmt.Sprintf(
			"- [%s] %s: %s (Date: %s)\n",
			record.Type,
			record.Title,
			record.Description,
			record.Date,
		)
	}

	// Create prompt for summarization
	systemPrompt := `You are a medical assistant AI. Analyze the provided health records and:
1. Generate a brief health summary (2-3 sentences)
2. List 3 key findings
3. Provide health recommendations

Format your response as JSON:
{
  "summary": "...",
  "findings": ["...", "...", "..."],
  "recommendations": "..."
}`

	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: recordsText,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to call OpenAI: %w", err)
	}

	// Parse response
	type SummaryResponse struct {
		Summary         string   `json:"summary"`
		Findings        []string `json:"findings"`
		Recommendations string   `json:"recommendations"`
	}

	var result SummaryResponse
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Summary, result.Findings, result.Recommendations, nil
}

// Example 3: Doctor Chat - Streaming Conversation
// This demonstrates a streaming chat where the AI responds character by character
func DoctorChatWithStreamingOpenAI(messages []openai.ChatCompletionMessage) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)

	// Create system message for doctor AI
	systemMessage := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: `You are a helpful medical assistant AI.
- Be empathetic and professional
- Ask clarifying questions about symptoms
- Provide general health information (not medical diagnosis)
- Recommend seeing a doctor for serious concerns
- Keep responses concise and clear`,
	}

	// Combine system message with conversation history
	allMessages := []openai.ChatCompletionMessage{systemMessage}
	allMessages = append(allMessages, messages...)

	// Create streaming request
	request := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    allMessages,
		MaxTokens:   500,
		Temperature: 0.8,
		Stream:      true, // Enable streaming
	}

	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	// Collect streamed response
	fullResponse := ""
	for {
		response, err := stream.Recv()
		if err != nil {
			// Stream ended or error
			break
		}

		// Each delta contains a small piece of the response
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			fullResponse += chunk
			// In a real application, you'd send chunks to the client immediately
			fmt.Print(chunk)
		}
	}

	fmt.Println() // newline after streaming
	return fullResponse, nil
}

// Example 4: Using Go patterns to structure API calls
// This shows how to properly handle errors and responses in Go

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SafeAPICall wraps API calls with error handling
func SafeAPICall(fn func() (string, error)) (string, error) {
	if fn == nil {
		return "", fmt.Errorf("function is nil")
	}

	result, err := fn()
	if err != nil {
		// Log the error for debugging
		log.Printf("API call failed: %v", err)

		// Return a user-friendly error
		return "", fmt.Errorf("failed to process request: %w", err)
	}

	return result, nil
}

// Example 5: Retry logic for API calls (often needed due to rate limiting)
import "time"

// RetryConfig defines retry parameters
type RetryConfig struct {
	MaxAttempts int
	DelayMs     int
	BackoffMult float64
}

// CallWithRetry retries an API call with exponential backoff
func CallWithRetry(fn func() (string, error), config RetryConfig) (string, error) {
	var lastErr error
	delay := time.Duration(config.DelayMs) * time.Millisecond

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil // Success!
		}

		lastErr = err
		log.Printf("Attempt %d failed: %v. Retrying in %v...", attempt+1, err, delay)

		// Wait before retrying
		time.Sleep(delay)

		// Exponential backoff: increase delay for next attempt
		delay = time.Duration(float64(delay) * config.BackoffMult)
	}

	return "", fmt.Errorf("all %d attempts failed. Last error: %w", config.MaxAttempts, lastErr)
}

// Usage example:
// result, err := CallWithRetry(
//     func() (string, error) {
//         return ScanPrescriptionWithOpenAI(imageData)
//     },
//     RetryConfig{MaxAttempts: 3, DelayMs: 1000, BackoffMult: 2.0},
// )

// Example 6: Caching API responses to reduce costs
import "sync"

// CacheEntry stores cached data with expiry
type CacheEntry struct {
	Data      string
	ExpiresAt time.Time
}

// APICache provides simple caching for API responses
type APICache struct {
	mu    sync.RWMutex
	cache map[string]CacheEntry
	ttl   time.Duration
}

// NewAPICache creates a new cache with time-to-live
func NewAPICache(ttl time.Duration) *APICache {
	return &APICache{
		cache: make(map[string]CacheEntry),
		ttl:   ttl,
	}
}

// Get retrieves cached data
func (c *APICache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return "", false
	}

	// Check if cache entry has expired
	if time.Now().After(entry.ExpiresAt) {
		return "", false
	}

	return entry.Data, true
}

// Set stores data in cache
func (c *APICache) Set(key string, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Usage example:
// cache := NewAPICache(5 * time.Hour)
//
// if cached, ok := cache.Get("prescription_scan_123"); ok {
//     return cached, nil
// }
//
// result, _ := ScanPrescriptionWithOpenAI(imageData)
// cache.Set("prescription_scan_123", result)
// return result, nil

// Example 7: Integration with the backend service
// This shows how to integrate the API calls with the existing AI service

// EnhancedAIService adds real OpenAI integration to the mock service
type EnhancedAIService struct {
	openaiKey string
	cache     *APICache
	retryConfig RetryConfig
}

// NewEnhancedAIService creates an AI service with OpenAI integration
func NewEnhancedAIService(openaiKey string) *EnhancedAIService {
	return &EnhancedAIService{
		openaiKey: openaiKey,
		cache: NewAPICache(5 * time.Hour),
		retryConfig: RetryConfig{
			MaxAttempts: 3,
			DelayMs:     1000,
			BackoffMult: 2.0,
		},
	}
}

// ScanPrescription with real OpenAI integration
func (s *EnhancedAIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("scan_%s", userID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		log.Printf("Cache hit for scan_%s", userID)

		// Parse cached JSON
		var result map[string]string
		json.Unmarshal([]byte(cached), &result)
		return result, nil
	}

	// Call OpenAI API with retry logic
	result, err := CallWithRetry(
		func() (string, error) {
			prescription, err := ScanPrescriptionWithOpenAI(imageData)
			if err != nil {
				return "", err
			}

			// Convert to map
			data := map[string]string{
				"medication": prescription.Medication,
				"dosage":     prescription.Dosage,
				"frequency":  prescription.Frequency,
				"duration":   prescription.Duration,
				"indication": prescription.Indication,
			}

			jsonData, _ := json.Marshal(data)
			return string(jsonData), nil
		},
		s.retryConfig,
	)

	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.Set(cacheKey, result)

	// Parse and return
	var output map[string]string
	json.Unmarshal([]byte(result), &output)
	return output, nil
}

// Example 8: Main function showing how to use everything
func MainExample() {
	// This is how you would use the integration in your code

	// 1. Set API key
	os.Setenv("OPENAI_API_KEY", "sk-xxxxx")

	// 2. Create service
	service := NewEnhancedAIService(os.Getenv("OPENAI_API_KEY"))

	// 3. Use it
	prescription, err := service.ScanPrescription("user123", []byte{/* image data */})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Printf("Scanned prescription: %+v", prescription)
}
