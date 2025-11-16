package examples

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	bedrockruntime "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

// AWS AI Integration Examples
// This shows how to use AWS services for medical image analysis and text generation

// Example 1: Using AWS Rekognition for Prescription Scanning
// Rekognition is AWS's computer vision service that can extract text from images

func ScanPrescriptionWithAWSRekognition(imageData []byte) (*PrescriptionData, error) {
	// Step 1: Load AWS configuration
	// This uses AWS credentials from environment or ~/.aws/credentials
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Step 2: Create Rekognition client
	client := rekognition.NewFromConfig(cfg)

	// Step 3: Call DetectText to extract text from the image
	ctx := context.Background()
	result, err := client.DetectText(ctx, &rekognition.DetectTextInput{
		Image: &types.Image{
			Bytes: imageData,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to detect text: %w", err)
	}

	// Step 4: Extract text from response
	// DetectText returns an array of detected text blocks
	fullText := ""
	for _, textDetection := range result.TextDetections {
		// Only use detected text with high confidence
		if textDetection.Confidence != nil && *textDetection.Confidence > 80.0 {
			if textDetection.DetectedText != nil {
				fullText += *textDetection.DetectedText + " "
			}
		}
	}

	// Step 5: Parse extracted text into structured data
	// In production, you'd use more sophisticated parsing
	prescription := &PrescriptionData{
		Medication: "Extracted medication",
		Dosage:     "Extracted dosage",
		Frequency:  "Extracted frequency",
	}

	return prescription, nil
}

// Example 2: Using AWS Bedrock for Health Summarization
// Bedrock gives you access to powerful foundation models like Claude, Llama, etc.

type BedrockRequest struct {
	Anthropic struct {
		Version string `json:"version"`
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"anthropic_version"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

func SummarizeHealthWithBedrock(records []HealthRecord) (string, []string, string, error) {
	// Step 1: Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Step 2: Create Bedrock Runtime client
	// This client invokes foundation models
	client := bedrockruntime.NewFromConfig(cfg)

	// Step 3: Build the prompt
	// Structure your prompt clearly for better results
	prompt := `You are a medical assistant. Analyze these health records:
`
	for _, record := range records {
		prompt += fmt.Sprintf("- %s: %s\n", record.Title, record.Description)
	}

	prompt += `

Provide response in JSON format:
{
  "summary": "2-3 sentence overview",
  "findings": ["finding 1", "finding 2", "finding 3"],
  "recommendations": "health recommendations"
}`

	// Step 4: Create the request payload
	// Different models have different request formats
	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-06-01",
		"max_tokens":        1000,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	requestJSON, _ := json.Marshal(requestBody)

	// Step 5: Invoke the model
	// "claude-3-sonnet-20240229-v1:0" is the model ID
	ctx := context.Background()
	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     "claude-3-sonnet-20240229-v1:0",
		ContentType: "application/json",
		Body:        requestJSON,
	})
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Step 6: Parse the response
	type BedrockResponse struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	var response BedrockResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Step 7: Extract and parse the model's response
	if len(response.Content) == 0 {
		return "", nil, "", fmt.Errorf("empty response from model")
	}

	modelResponse := response.Content[0].Text

	type SummaryResponse struct {
		Summary         string   `json:"summary"`
		Findings        []string `json:"findings"`
		Recommendations string   `json:"recommendations"`
	}

	var summary SummaryResponse
	err = json.Unmarshal([]byte(modelResponse), &summary)
	if err != nil {
		return "", nil, "", fmt.Errorf("failed to parse model response: %w", err)
	}

	return summary.Summary, summary.Findings, summary.Recommendations, nil
}

// Example 3: Doctor Chat using Claude from Bedrock

func DoctorChatWithBedrock(conversationHistory []map[string]string, newMessage string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	client := bedrockruntime.NewFromConfig(cfg)

	// Build messages array with conversation history
	messages := []map[string]interface{}{}

	// Add system message
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": "You are a helpful medical assistant AI. Be empathetic and professional.",
	})

	// Add conversation history
	for _, msg := range conversationHistory {
		messages = append(messages, map[string]interface{}{
			"role":    msg["role"],
			"content": msg["content"],
		})
	}

	// Add new message
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": newMessage,
	})

	// Create request
	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-06-01",
		"max_tokens":        500,
		"messages":          messages,
	}

	requestJSON, _ := json.Marshal(requestBody)

	ctx := context.Background()
	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     "claude-3-sonnet-20240229-v1:0",
		ContentType: "application/json",
		Body:        requestJSON,
	})
	if err != nil {
		return "", err
	}

	// Parse response
	var response map[string]interface{}
	json.Unmarshal(result.Body, &response)

	// Extract text from response
	if content, ok := response["content"].([]interface{}); ok && len(content) > 0 {
		if textBlock, ok := content[0].(map[string]interface{}); ok {
			return textBlock["text"].(string), nil
		}
	}

	return "", fmt.Errorf("unexpected response format")
}

// Example 4: AWS Configuration and Setup

type AWSConfig struct {
	Region   string
	RoleARN  string // For cross-account access if needed
}

func LoadAWSConfig() *AWSConfig {
	return &AWSConfig{
		Region:  os.Getenv("AWS_REGION"),
		RoleARN: os.Getenv("AWS_ROLE_ARN"),
	}
}

func (c *AWSConfig) Validate() error {
	if c.Region == "" {
		return fmt.Errorf("AWS_REGION not set")
	}

	return nil
}

// Example 5: Error Handling for AWS Services

// AWSServiceError wraps AWS errors with additional context
type AWSServiceError struct {
	Service   string // e.g., "Rekognition", "Bedrock"
	Operation string // e.g., "DetectText", "InvokeModel"
	Code      string // AWS error code
	Message   string
	Original  error
}

func (e *AWSServiceError) Error() string {
	return fmt.Sprintf("AWS %s %s failed [%s]: %s - %v",
		e.Service, e.Operation, e.Code, e.Message, e.Original)
}

// Example 6: Batch Processing with AWS Rekognition
// Process multiple images in parallel for better performance

import (
	"sync"
)

func BatchScanPrescriptionsWithAWS(prescriptions map[string][]byte) (map[string]*PrescriptionData, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := rekognition.NewFromConfig(cfg)
	results := make(map[string]*PrescriptionData)

	// Use WaitGroup to process images concurrently
	// This allows multiple images to be processed at the same time
	var wg sync.WaitGroup
	resultsChan := make(chan struct {
		userID string
		data   *PrescriptionData
		err    error
	}, len(prescriptions))

	ctx := context.Background()

	// Create goroutine for each image
	for userID, imageData := range prescriptions {
		wg.Add(1)

		// Launch concurrent processing
		go func(uid string, img []byte) {
			defer wg.Done()

			result, err := client.DetectText(ctx, &rekognition.DetectTextInput{
				Image: &types.Image{
					Bytes: img,
				},
			})

			if err != nil {
				resultsChan <- struct {
					userID string
					data   *PrescriptionData
					err    error
				}{uid, nil, err}
				return
			}

			// Parse text detections
			prescription := parseRekognitionResults(result)
			resultsChan <- struct {
				userID string
				data   *PrescriptionData
				err    error
			}{uid, prescription, nil}
		}(userID, imageData)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for res := range resultsChan {
		if res.err != nil {
			log.Printf("Error processing %s: %v", res.userID, res.err)
		} else {
			results[res.userID] = res.data
		}
	}

	return results, nil
}

// Helper function to parse Rekognition results
func parseRekognitionResults(result *rekognition.DetectTextOutput) *PrescriptionData {
	// Extract text blocks from Rekognition response
	prescription := &PrescriptionData{}

	for _, detection := range result.TextDetections {
		if detection.DetectedText != nil {
			// Simple parsing - in production, use more sophisticated logic
			log.Printf("Detected: %s (confidence: %.2f%%)",
				*detection.DetectedText,
				*detection.Confidence)
		}
	}

	return prescription
}

// Example 7: Caching API Responses

// S3Cache stores responses in AWS S3 for distributed caching
// This is useful for expensive API calls

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3Cache struct {
	s3Client   *s3.Client
	bucketName string
}

func NewS3Cache(cfg config.Config, bucketName string) *S3Cache {
	return &S3Cache{
		s3Client:   s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}
}

// Get retrieves cached data from S3
func (c *S3Cache) Get(key string) (string, error) {
	ctx := context.Background()

	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucketName,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(result.Body)
	return buffer.String(), nil
}

// Set stores data in S3 cache
func (c *S3Cache) Set(key string, data string) error {
	ctx := context.Background()

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.bucketName,
		Key:    &key,
		Body:   bytes.NewReader([]byte(data)),
	})
	return err
}

// Example 8: Streaming Responses from Bedrock
// Some models support streaming for real-time responses

import "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

func DoctorChatWithStreamingBedrock(message string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	client := bedrockruntime.NewFromConfig(cfg)

	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-06-01",
		"max_tokens":        500,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": message,
			},
		},
	}

	requestJSON, _ := json.Marshal(requestBody)

	ctx := context.Background()

	// Use InvokeModelWithResponseStream for streaming
	result, err := client.InvokeModelWithResponseStream(ctx,
		&bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     "claude-3-sonnet-20240229-v1:0",
			ContentType: "application/json",
			Body:        requestJSON,
		})
	if err != nil {
		return "", err
	}

	// Collect streamed chunks
	fullResponse := ""
	for event := range result.EventStream.Events() {
		switch e := event.(type) {
		case *types.ContentBlockDeltaEvent:
			if delta := e.Delta.(*types.ContentBlockDelta); delta != nil {
				// Process streaming text delta
				fmt.Print(".")
			}
		}
	}

	return fullResponse, nil
}

// Example 9: AWS Integration Best Practices

// AWSHealthAIService integrates AWS services for healthcare
type AWSHealthAIService struct {
	rekognitionClient *rekognition.Client
	bedrockClient     *bedrockruntime.Client
	cache             *APICache
	config            *AWSConfig
}

func NewAWSHealthAIService(cfg *AWSConfig) (*AWSHealthAIService, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &AWSHealthAIService{
		rekognitionClient: rekognition.NewFromConfig(awsCfg),
		bedrockClient:     bedrockruntime.NewFromConfig(awsCfg),
		cache:             NewAPICache(5 * time.Hour),
		config:            cfg,
	}, nil
}

// ScanPrescription with AWS
func (s *AWSHealthAIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("aws_scan_%s", userID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		var result map[string]string
		json.Unmarshal([]byte(cached), &result)
		return result, nil
	}

	// Use Rekognition
	prescription, err := ScanPrescriptionWithAWSRekognition(imageData)
	if err != nil {
		return nil, err
	}

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

// Key AWS Integration Points:
// 1. Always load config at the start
// 2. Use context for request lifecycle
// 3. Defer client closing to prevent resource leaks
// 4. Implement proper error handling
// 5. Cache expensive API calls
// 6. Use goroutines for parallel processing
// 7. Set appropriate confidence thresholds
// 8. Use appropriate model IDs for Bedrock
// 9. Monitor API costs and usage
// 10. Implement circuit breakers for fault tolerance
