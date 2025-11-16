package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	vision "cloud.google.com/go/vision/v2"
	"google.cloud.org/genai"
)

// Google Cloud Vision API Example
// This shows how to use Google's Vision API for prescription scanning
// and Gemini API for health summarization

// Example 1: Using Google Cloud Vision for Prescription OCR
// The Vision API uses Optical Character Recognition to extract text from images

func ScanPrescriptionWithGoogleVision(imageData []byte) (*PrescriptionData, error) {
	// Step 1: Create a context
	// Context is a Go pattern for managing request deadlines and cancellation
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

// Helper function to parse prescription text
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

// Example 2: Using Google Gemini API for Health Summarization
// Gemini is Google's most advanced language model for content generation

type GeminiClient struct {
	apiKey string
	client *genai.Client
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()

	// Create the Gemini client with API key
	client, err := genai.NewClient(ctx, genai.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		apiKey: apiKey,
		client: client,
	}, nil
}

// SummarizeHealthWithGemini generates a health summary using Gemini
func (g *GeminiClient) SummarizeHealthWithGemini(records []HealthRecord) (string, []string, string, error) {
	ctx := context.Background()

	// Get the Gemini 1.5 Pro model
	model := g.client.GenerativeModel("gemini-1.5-pro")

	// Configure model parameters
	// Temperature: How creative (0=factual, 1=creative)
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(1000)

	// Build prompt from health records
	prompt := buildHealthPrompt(records)

	// Generate response
	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract response text
	if len(response.Candidates) == 0 {
		return "", nil, "", fmt.Errorf("no response from Gemini")
	}

	responseText := extractTextFromCandidate(response.Candidates[0])

	// Parse JSON response
	type SummaryResponse struct {
		Summary         string   `json:"summary"`
		Findings        []string `json:"findings"`
		Recommendations string   `json:"recommendations"`
	}

	var result SummaryResponse
	err = json.Unmarshal([]byte(responseText), &result)
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Summary, result.Findings, result.Recommendations, nil
}

// Helper function to build prompt
func buildHealthPrompt(records []HealthRecord) string {
	prompt := `Analyze these health records and provide a JSON response with:
{
  "summary": "2-3 sentence overview",
  "findings": ["key finding 1", "key finding 2", "key finding 3"],
  "recommendations": "health recommendations"
}

Health Records:`

	for _, record := range records {
		prompt += fmt.Sprintf("\n- %s: %s", record.Title, record.Description)
	}

	return prompt
}

// Helper to extract text from response
func extractTextFromCandidate(candidate interface{}) string {
	// This is simplified - actual implementation depends on response structure
	return "parsed response"
}

// Example 3: Streaming Response with Gemini
// Gemini supports streaming which sends response chunks as they're generated

func (g *GeminiClient) DoctorChatWithStreaming(messages []string) (string, error) {
	ctx := context.Background()

	model := g.client.GenerativeModel("gemini-1.5-pro")

	// Build chat history
	cs := model.StartChat()

	// Add system context
	systemPrompt := `You are a helpful medical assistant AI.
- Be empathetic and professional
- Ask clarifying questions about symptoms
- Provide general health information
- Recommend seeing a doctor for serious concerns`

	// Add system message
	cs.History = append(cs.History, &genai.Content{
		Role: "user",
		Parts: []genai.Part{
			genai.Text(systemPrompt),
		},
	})

	// Add previous messages to chat
	for _, msg := range messages {
		cs.History = append(cs.History, &genai.Content{
			Role: "user",
			Parts: []genai.Part{
				genai.Text(msg),
			},
		})
	}

	// Send message and get streaming response
	iter := cs.SendMessageStream(ctx, genai.Text("What health advice do you have?"))

	fullResponse := ""
	for {
		chunk, err := iter.Next()
		if err != nil {
			break // Stream ended
		}

		// Each chunk contains a part of the response
		if len(chunk.Content.Parts) > 0 {
			text := fmt.Sprintf("%v", chunk.Content.Parts[0])
			fullResponse += text
			// In a real app, send this to the client immediately
			fmt.Print(text)
		}
	}

	fmt.Println() // newline
	return fullResponse, nil
}

// Example 4: Best Practices for Google Cloud Integration

// GoogleCloudConfig holds configuration for Google Cloud services
type GoogleCloudConfig struct {
	ProjectID   string
	APIKey      string
	CredentialsPath string // Path to service account JSON
}

// LoadGoogleCloudConfig loads configuration from environment
func LoadGoogleCloudConfig() *GoogleCloudConfig {
	return &GoogleCloudConfig{
		ProjectID:   os.Getenv("GOOGLE_CLOUD_PROJECT"),
		APIKey:      os.Getenv("GOOGLE_API_KEY"),
		CredentialsPath: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}
}

// Validate checks if config is valid
func (c *GoogleCloudConfig) Validate() error {
	if c.APIKey == "" && c.CredentialsPath == "" {
		return fmt.Errorf("either GOOGLE_API_KEY or GOOGLE_APPLICATION_CREDENTIALS must be set")
	}

	if c.CredentialsPath != "" {
		// Check if credentials file exists
		if _, err := os.Stat(c.CredentialsPath); os.IsNotExist(err) {
			return fmt.Errorf("credentials file not found: %s", c.CredentialsPath)
		}
	}

	return nil
}

// Example 5: Error Handling Pattern for Google Cloud APIs

// GoogleAPIError wraps errors from Google Cloud APIs
type GoogleAPIError struct {
	Service   string // e.g., "Vision", "Gemini"
	Operation string // e.g., "DetectText", "GenerateContent"
	Original  error
}

func (e *GoogleAPIError) Error() string {
	return fmt.Sprintf("Google %s API error in %s: %v", e.Service, e.Operation, e.Original)
}

// WrapGoogleError wraps an error with service context
func WrapGoogleError(service, operation string, err error) error {
	return &GoogleAPIError{
		Service:   service,
		Operation: operation,
		Original:  err,
	}
}

// Example usage:
// _, err := client.DetectTexts(ctx, image, nil)
// if err != nil {
//     return nil, WrapGoogleError("Vision", "DetectTexts", err)
// }

// Example 6: Batch Processing Multiple Images
// Shows how to process multiple prescription images efficiently

func BatchScanPrescriptions(prescriptionImages map[string][]byte) (map[string]*PrescriptionData, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	results := make(map[string]*PrescriptionData)

	// Process each image
	for userID, imageData := range prescriptionImages {
		image := vision.NewImageFromBytes(imageData)

		annotations, err := client.DetectTexts(ctx, image, nil)
		if err != nil {
			log.Printf("Error processing image for user %s: %v", userID, err)
			continue
		}

		// Extract text from annotations
		fullText := ""
		for _, annotation := range annotations {
			fullText += annotation.GetDescription()
		}

		// Parse prescription
		prescription := parsePrescriptionText(fullText)
		results[userID] = prescription
	}

	return results, nil
}

// Example 7: Caching with TTL for Google Cloud Responses

// CloudCacheEntry stores cached response with metadata
type CloudCacheEntry struct {
	Data        string
	ExpiresAt   time.Time
	Source      string // e.g., "Vision", "Gemini"
	RequestHash string // Hash of the request to verify cache hit
}

// VerifyCache checks if cache entry matches request
func (e *CloudCacheEntry) VerifyCache(requestHash string) bool {
	if e.RequestHash != requestHash {
		return false // Different request
	}

	if time.Now().After(e.ExpiresAt) {
		return false // Expired
	}

	return true
}

// Example 8: Integration with existing backend service

// GoogleCloudAIService integrates Google Cloud APIs with the backend
type GoogleCloudAIService struct {
	config         *GoogleCloudConfig
	visionClient   *vision.ImageAnnotatorClient
	geminiClient   *GeminiClient
	cache          *APICache
}

// NewGoogleCloudAIService creates a new service
func NewGoogleCloudAIService(config *GoogleCloudConfig) (*GoogleCloudAIService, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Create Vision client
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vision client: %w", err)
	}

	// Create Gemini client
	geminiClient, err := NewGeminiClient(config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GoogleCloudAIService{
		config:       config,
		visionClient: visionClient,
		geminiClient: geminiClient,
		cache:        NewAPICache(5 * time.Hour),
	}, nil
}

// ScanPrescription with Google Vision
func (s *GoogleCloudAIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
	cacheKey := fmt.Sprintf("vision_%s", userID)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		var result map[string]string
		json.Unmarshal([]byte(cached), &result)
		return result, nil
	}

	// Call Vision API
	prescription, err := ScanPrescriptionWithGoogleVision(imageData)
	if err != nil {
		return nil, err
	}

	// Convert to map
	result := map[string]string{
		"medication": prescription.Medication,
		"dosage":     prescription.Dosage,
		"frequency":  prescription.Frequency,
	}

	// Cache result
	jsonData, _ := json.Marshal(result)
	s.cache.Set(cacheKey, string(jsonData))

	return result, nil
}

// Example 9: How to set up Google Cloud credentials

func SetupGoogleCloudCredentials() {
	// Option 1: Using API Key (simpler, for APIs that support it)
	os.Setenv("GOOGLE_API_KEY", "your-api-key")

	// Option 2: Using Service Account (more secure, for production)
	// 1. Download JSON credentials from Google Cloud Console
	// 2. Set path to the file
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/service-account.json")

	// Option 3: Set in .env file (recommended for local development)
	// GOOGLE_API_KEY=your-api-key
	// GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
}

// Example 10: Testing the integration

func TestGoogleCloudIntegration() {
	config := LoadGoogleCloudConfig()

	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	service, err := NewGoogleCloudAIService(config)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// Read a test image
	imageData, err := ioutil.ReadFile("/path/to/test/prescription.jpg")
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	// Scan prescription
	result, err := service.ScanPrescription("test-user", imageData)
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	log.Printf("Scanned prescription: %+v", result)
}

// Key Points for Google Cloud Integration:
// 1. Use context.Background() for request context
// 2. Always defer client.Close() to clean up resources
// 3. Cache API responses to reduce costs
// 4. Implement retry logic with exponential backoff
// 5. Use structured error handling with custom error types
// 6. Validate configuration before using services
// 7. Use environment variables for sensitive data
// 8. Implement request deduplication to avoid duplicate API calls
