# AI Service Integration Examples

This directory contains working Go examples showing how to integrate real AI services into the Clarity backend.

## Overview

Three complete integration examples are provided:

1. **OpenAI Integration** (`openai_integration.go`)
   - ChatGPT 4 Vision for prescription scanning
   - ChatGPT for health summarization and doctor chat
   - Streaming responses
   - Error handling and retries

2. **Google Cloud Integration** (`google_cloud_integration.go`)
   - Google Cloud Vision for OCR
   - Gemini for health analysis
   - Streaming with Gemini
   - Batch processing

3. **AWS Integration** (`aws_integration.go`)
   - AWS Rekognition for prescription scanning
   - AWS Bedrock for health analysis
   - Parallel processing with goroutines
   - S3 caching

## Key Go Programming Concepts Used

### 1. **Context** - Request lifecycle management
```go
ctx := context.Background()
// Use ctx for cancellation, deadlines
result, err := client.DoSomething(ctx, request)
```

### 2. **Defer** - Cleanup resources
```go
client, err := createClient()
defer client.Close() // Runs at end of function
```

### 3. **Error Handling** - Go's error pattern
```go
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}
```

### 4. **Interfaces** - Abstraction for swapping implementations
```go
type AIProvider interface {
    ScanPrescription(imageData []byte) (map[string]string, error)
    SummarizeHealth(records []HealthRecord) (string, error)
}
```

### 5. **Goroutines** - Concurrent processing
```go
go func() {
    result := processImage(imageData)
}()
```

### 6. **WaitGroup** - Synchronize goroutines
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // do work
}()
wg.Wait()
```

### 7. **Channels** - Communicate between goroutines
```go
resultChan := make(chan Result)
go func() {
    resultChan <- computeResult()
}()
result := <-resultChan
```

### 8. **Mutex** - Protect shared data
```go
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}
```

## Installation

### 1. Install Dependencies

**OpenAI Integration:**
```bash
go get github.com/sashabaranov/go-openai
```

**Google Cloud Integration:**
```bash
go get cloud.google.com/go/vision/v2
go get google.cloud.org/genai
```

**AWS Integration:**
```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/service/rekognition
go get github.com/aws/aws-sdk-go-v2/service/bedrockruntime
```

### 2. Set Environment Variables

**OpenAI:**
```bash
export OPENAI_API_KEY="sk-xxxxx"
```

**Google Cloud:**
```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials.json"
export GOOGLE_API_KEY="your-api-key"
```

**AWS:**
```bash
export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="xxxxx"
export AWS_SECRET_ACCESS_KEY="xxxxx"
```

## How to Use These Examples

### Basic Pattern for Each Service

#### OpenAI
```go
// 1. Create service
service := NewEnhancedAIService(os.Getenv("OPENAI_API_KEY"))

// 2. Use it
prescription, err := service.ScanPrescription("user123", imageData)

// 3. Handle errors
if err != nil {
    log.Printf("Error: %v", err)
}
```

#### Google Cloud
```go
// 1. Load config
config := LoadGoogleCloudConfig()
config.Validate()

// 2. Create service
service, _ := NewGoogleCloudAIService(config)

// 3. Use it
result, _ := service.ScanPrescription("user123", imageData)
```

#### AWS
```go
// 1. Load config
config := LoadAWSConfig()
config.Validate()

// 2. Create service
service, _ := NewAWSHealthAIService(config)

// 3. Use it
result, _ := service.ScanPrescription("user123", imageData)
```

## Code Explanations

### Example: Prescription Scanning

#### OpenAI Version
```go
func ScanPrescriptionWithOpenAI(imageData []byte) (*PrescriptionData, error) {
    // Step 1: Create API client
    client := openai.NewClient(apiKey)

    // Step 2: Convert image to base64 (required by OpenAI)
    encodedImage := base64.StdEncoding.EncodeToString(imageData)

    // Step 3: Create request with vision model
    request := openai.ChatCompletionRequest{
        Model: openai.GPT4VisionPreview,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleUser,
                MultiContent: []openai.ChatMessagePart{
                    {
                        Type: openai.ChatMessagePartTypeImageURL,
                        ImageURL: &openai.ImageURL{
                            URL: fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage),
                        },
                    },
                    {
                        Type: openai.ChatMessagePartTypeText,
                        Text: "Extract medication information...",
                    },
                },
            },
        },
    }

    // Step 4: Call API
    ctx := context.Background()
    response, err := client.CreateChatCompletion(ctx, request)

    // Step 5: Parse response
    return parseResponse(response), nil
}
```

**What's happening:**
1. We create an OpenAI client with our API key
2. Convert the image to base64 (OpenAI requires this format)
3. Create a request with both image and text
4. Call the API with context
5. Parse the returned JSON

#### Key Patterns:
- **Context**: Used for timeout and cancellation management
- **Error propagation**: Using `fmt.Errorf` with `%w` for error wrapping
- **Type safety**: Using structs for requests/responses

### Example: Health Summarization

```go
func SummarizeHealthWithGemini(records []HealthRecord) (string, []string, string, error) {
    // Get model
    model := client.GenerativeModel("gemini-1.5-pro")

    // Configure parameters
    model.SetTemperature(0.7)        // Creativity level
    model.SetMaxOutputTokens(1000)   // Response length limit

    // Build prompt from records
    prompt := buildHealthPrompt(records)

    // Generate response
    response, err := model.GenerateContent(ctx, genai.Text(prompt))

    // Parse response
    var result SummaryResponse
    json.Unmarshal([]byte(responseText), &result)

    return result.Summary, result.Findings, result.Recommendations, nil
}
```

**What's happening:**
1. Get the model instance from Gemini client
2. Configure generation parameters (temperature, max tokens)
3. Build the prompt from input data
4. Call the API
5. Parse JSON response into structured data

### Example: Concurrent Processing

```go
func BatchScanPrescriptionsWithAWS(prescriptions map[string][]byte) (map[string]*PrescriptionData, error) {
    var wg sync.WaitGroup  // WaitGroup tracks goroutines
    results := make(map[string]*PrescriptionData)

    // Process each image concurrently
    for userID, imageData := range prescriptions {
        wg.Add(1)  // Increment counter

        go func(uid string, img []byte) {
            defer wg.Done()  // Decrement when done

            // Process image
            result, err := client.DetectText(ctx, &rekognition.DetectTextInput{
                Image: &types.Image{Bytes: img},
            })

            // Store result
            resultsChan <- struct{userID string; data *PrescriptionData}{uid, parsed}
        }(userID, imageData)
    }

    wg.Wait()  // Wait for all goroutines
    return results, nil
}
```

**What's happening:**
1. `WaitGroup` tracks the number of goroutines
2. For each image, we launch a goroutine (concurrent execution)
3. Each goroutine processes an image independently
4. We wait for all to complete before returning
5. Results are sent through a channel

**Performance benefit**: All images processed in parallel instead of one-by-one

### Example: Caching with TTL

```go
type APICache struct {
    mu    sync.RWMutex              // Protects the cache
    cache map[string]CacheEntry
    ttl   time.Duration              // Time to live
}

func (c *APICache) Get(key string) (string, bool) {
    c.mu.RLock()                    // Lock for reading
    defer c.mu.RUnlock()

    entry, exists := c.cache[key]
    if !exists {
        return "", false
    }

    // Check if expired
    if time.Now().After(entry.ExpiresAt) {
        return "", false
    }

    return entry.Data, true
}

func (c *APICache) Set(key string, data string) {
    c.mu.Lock()                     // Lock for writing
    defer c.mu.Unlock()

    c.cache[key] = CacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(c.ttl),  // Set expiry time
    }
}
```

**What's happening:**
1. `RWMutex` allows multiple readers or one writer (thread-safe)
2. RLock for read operations, Lock for write operations
3. Automatic expiry checking on each Get
4. Defer ensures locks are always released

**Benefit**: Reduces API costs by not calling the same API twice within TTL

### Example: Error Handling Pattern

```go
// Custom error type
type GoogleAPIError struct {
    Service   string
    Operation string
    Original  error
}

func (e *GoogleAPIError) Error() string {
    return fmt.Sprintf("Google %s API error in %s: %v",
        e.Service, e.Operation, e.Original)
}

// Using it
if err != nil {
    return nil, &GoogleAPIError{
        Service:   "Vision",
        Operation: "DetectTexts",
        Original:  err,
    }
}
```

**What's happening:**
1. Define custom error type with context
2. Implement Error() interface
3. Wrap errors with additional information
4. Caller gets complete error context

## Integration Checklist

- [ ] Install required Go packages
- [ ] Set environment variables
- [ ] Choose which AI provider to use
- [ ] Copy relevant example file into your project
- [ ] Update `backend/services/ai_service.go` to use the new implementation
- [ ] Test with sample prescription images
- [ ] Monitor API costs
- [ ] Implement rate limiting if needed

## Common Mistakes to Avoid

1. **Forgetting to defer Close()**
   ```go
   client, _ := createClient()
   defer client.Close()  // Don't forget this!
   ```

2. **Not handling errors**
   ```go
   // Wrong
   response, _ := api.Call()  // Ignoring error!

   // Right
   response, err := api.Call()
   if err != nil {
       return nil, err
   }
   ```

3. **Race conditions with goroutines**
   ```go
   // Wrong - concurrent map access without mutex
   results := make(map[string]Result)
   go func() {
       results["key"] = value  // Race condition!
   }()

   // Right - use channel or mutex
   resultChan := make(chan Result)
   go func() {
       resultChan <- value  // Safe
   }()
   ```

4. **Not setting context timeout**
   ```go
   // Wrong - might wait forever
   result, _ := api.Call(context.Background())

   // Right - set timeout
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   result, _ := api.Call(ctx)
   ```

## Testing

### Unit Test Example
```go
func TestScanPrescription(t *testing.T) {
    // Load test image
    imageData, _ := ioutil.ReadFile("test_prescription.jpg")

    // Call function
    result, err := ScanPrescriptionWithOpenAI(imageData)

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if result == nil {
        t.Fatal("Expected result, got nil")
    }

    if result.Medication == "" {
        t.Error("Expected medication, got empty string")
    }
}
```

## Performance Tips

1. **Cache responses** - Store API results to avoid repeated calls
2. **Batch requests** - Process multiple items in parallel with goroutines
3. **Use streaming** - For long responses, use streaming APIs
4. **Set appropriate timeouts** - Prevent hanging requests
5. **Implement rate limiting** - Respect API quotas

## Cost Optimization

1. **Cache API calls** - Reduces number of API calls
2. **Batch processing** - More efficient than individual calls
3. **Use appropriate models** - Cheaper models for simpler tasks
4. **Monitor usage** - Track API costs regularly
5. **Implement circuit breakers** - Stop calling failing APIs

## Next Steps

1. Choose one AI provider to start with
2. Copy the example into your project
3. Update `ai_service.go` to use real implementation
4. Set up environment variables
5. Test with sample data
6. Monitor performance and costs
7. Consider implementing caching layer
8. Add monitoring and logging

## Support

Each example includes:
- Detailed comments explaining each step
- Error handling patterns
- Best practices for production
- Integration with existing backend structure

Study the examples carefully to understand:
- How Go clients work
- How to handle errors
- How to use goroutines safely
- How to implement caching
- How to structure API calls

Good luck integrating AI into Clarity! 🚀
