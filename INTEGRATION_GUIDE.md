# Integration Guide - Connecting to Real Services

This guide helps you integrate real AI, email, and cloud services into Clarity.

## Table of Contents
1. [AI Service Integration](#ai-service-integration)
2. [Email Service Integration](#email-service-integration)
3. [Cloud Database Setup](#cloud-database-setup)
4. [Frontend-Backend Connection](#frontend-backend-connection)
5. [Testing & Debugging](#testing--debugging)

---

## AI Service Integration

### Option 1: OpenAI API (ChatGPT + Vision)

#### Setup
```bash
# 1. Get API key from https://platform.openai.com/api-keys
# 2. Add to .env
echo "AI_API_KEY=sk-your-api-key" >> .env
```

#### Update Backend Code
File: `backend/services/ai_service.go`

```go
import "github.com/sashabaranov/go-openai"

type AIService struct {
    db        *gorm.DB
    openaiKey string
}

// Replace ScanPrescription implementation
func (as *AIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
    client := openai.NewClient(as.openaiKey)

    // Encode image to base64
    encodedImage := base64.StdEncoding.EncodeToString(imageData)

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: openai.GPT4VisionPreview,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleUser,
                MultiContent: []openai.ChatMessagePart{
                    {
                        Type: openai.ChatMessagePartTypeImage,
                        ImageURL: &openai.ImageURL{
                            URL: fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage),
                        },
                    },
                    {
                        Type: openai.ChatMessagePartTypeText,
                        Text: "Extract medication information from this prescription: medication, dosage, frequency, duration, indication",
                    },
                },
            },
        },
    })

    // Parse response and return extracted data
    // Parse JSON response to extract: medication, dosage, frequency, etc.

    return extractedData, nil
}

// Replace DoctorChat implementation
func (as *AIService) DoctorChat(userID, conversationID, message string) (string, error) {
    client := openai.NewClient(as.openaiKey)

    // Get conversation history
    history, _ := as.GetConversationHistory(conversationID)

    messages := []openai.ChatCompletionMessage{
        {
            Role:    openai.ChatMessageRoleSystem,
            Content: "You are a helpful medical assistant. Provide health advice based on symptoms described.",
        },
    }

    // Add history
    for _, h := range history {
        messages = append(messages, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleUser,
            Content: h.Message,
        })
        messages = append(messages, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleAssistant,
            Content: h.Response,
        })
    }

    // Add new message
    messages = append(messages, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: message,
    })

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model:    openai.GPT3Dot5Turbo,
        Messages: messages,
    })

    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}
```

#### Update go.mod
```bash
go get github.com/sashabaranov/go-openai
```

### Option 2: Google Cloud Vision + PaLM

#### Setup
```bash
# 1. Create GCP project and enable Vision API
# 2. Download service account key (JSON)
# 3. Set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/key.json"
```

#### Update Backend Code
File: `backend/services/ai_service.go`

```go
import (
    vision "cloud.google.com/go/vision/v2"
    "cloud.google.com/go/aiplatform/apiv1"
)

func (as *AIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
    ctx := context.Background()
    client, _ := vision.NewImageAnnotatorClient(ctx)
    defer client.Close()

    image := vision.NewImageFromBytes(imageData)
    annotations, _ := client.DetectTexts(ctx, image, nil)

    // Parse OCR results to extract prescription data
    return parsePrescriptionText(annotations), nil
}
```

### Option 3: AWS Rekognition + Bedrock

#### Setup
```bash
# 1. Create AWS account and get credentials
# 2. Enable Rekognition and Bedrock
# 3. Add to .env
echo "AWS_ACCESS_KEY_ID=xxx" >> .env
echo "AWS_SECRET_ACCESS_KEY=xxx" >> .env
```

#### Update Backend Code
```go
import "github.com/aws/aws-sdk-go-v2/service/rekognition"

func (as *AIService) ScanPrescription(userID string, imageData []byte) (map[string]string, error) {
    svc := rekognition.NewFromConfig(awsCfg)

    result, _ := svc.DetectText(context.Background(), &rekognition.DetectTextInput{
        Image: &types.Image{
            Bytes: imageData,
        },
    })

    // Parse results
    return parseRekognitionResults(result), nil
}
```

---

## Email Service Integration

### Option 1: SendGrid

#### Setup
```bash
# 1. Get API key from https://sendgrid.com/
# 2. Add to .env
echo "SENDGRID_API_KEY=SG.xxx" >> .env
```

#### Update Backend Code
File: `backend/services/auth_service.go`

```go
import "github.com/sendgrid/sendgrid-go"
import "github.com/sendgrid/sendgrid-go/helpers/mail"

func (as *AuthService) SendOTP(email string) (string, error) {
    otp := generateOTP(as.config.OTPLength)

    from := mail.NewEmail("Clarity", "noreply@clarity.health")
    subject := "Your Clarity Health Verification Code"
    to := mail.NewEmail("User", email)

    htmlContent := fmt.Sprintf(`
        <h2>Your Verification Code</h2>
        <p>Enter this code to verify your email:</p>
        <h1>%s</h1>
        <p>This code expires in 10 minutes.</p>
    `, otp)

    message := mail.NewSingleEmail(from, subject, to, "OTP: "+otp, htmlContent)

    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    response, err := client.Send(message)

    if err != nil {
        return "", err
    }

    log.Printf("Email sent to %s, status: %d", email, response.StatusCode)

    // Store OTP as before
    return otp, nil
}
```

#### Update go.mod
```bash
go get github.com/sendgrid/sendgrid-go
```

### Option 2: AWS SES

```go
import "github.com/aws/aws-sdk-go-v2/service/ses"

func (as *AuthService) SendOTP(email string) (string, error) {
    otp := generateOTP(as.config.OTPLength)

    svc := ses.NewFromConfig(awsCfg)

    svc.SendEmail(context.Background(), &ses.SendEmailInput{
        Source: aws.String("noreply@clarity.health"),
        Destination: &types.Destination{
            ToAddresses: []string{email},
        },
        Message: &types.Message{
            Subject: &types.Content{
                Data: aws.String("Your Clarity Verification Code"),
            },
            Body: &types.Body{
                Html: &types.Content{
                    Data: aws.String(fmt.Sprintf(`<h1>%s</h1>`, otp)),
                },
            },
        },
    })

    return otp, nil
}
```

---

## Cloud Database Setup

### AWS RDS (PostgreSQL)

#### 1. Create RDS Instance
```bash
# Via AWS Console or CLI
aws rds create-db-instance \
  --db-instance-identifier clarity-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --master-username admin \
  --master-user-password YourSecurePassword
```

#### 2. Update Backend Configuration
File: `backend/.env`
```env
DB_TYPE=postgres
CLOUD_PROVIDER=aws
DB_HOST=clarity-db.xxxxxx.us-east-1.rds.amazonaws.com
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=YourSecurePassword
DB_NAME=clarity
```

#### 3. Update database/db.go
```go
func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
    switch cfg.Type {
    case "postgres":
        return newPostgresDB(cfg)
    // ...
    }
}

func newPostgresDB(cfg *config.DatabaseConfig) (Database, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return &PostgresDB{conn: db}, nil
}
```

#### 4. Update go.mod
```bash
go get gorm.io/driver/postgres
```

### Google Cloud SQL

#### 1. Create Instance
```bash
gcloud sql instances create clarity-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1
```

#### 2. Create Database
```bash
gcloud sql databases create clarity --instance=clarity-db
```

#### 3. Create User
```bash
gcloud sql users create clarity-user --instance=clarity-db --password
```

#### 4. Get Connection String
```bash
gcloud sql instances describe clarity-db --format="value(connectionName)"
# Output: project:region:clarity-db
```

#### 5. Update Backend
File: `backend/.env`
```env
DB_TYPE=postgres
CLOUD_PROVIDER=gcp
GOOGLE_CLOUD_SQL_CONNECTION=project:region:clarity-db
DB_USER=clarity-user
DB_PASSWORD=YourPassword
DB_NAME=clarity
```

### Azure Database for PostgreSQL

#### 1. Create Instance
```bash
az postgres server create \
  --resource-group clarity-rg \
  --name clarity-db \
  --location eastus \
  --admin-user azureuser \
  --admin-password YourPassword \
  --sku-name B_Gen5_1
```

#### 2. Get Connection String
```bash
# From Azure Portal: Connection Strings
postgresql://azureuser@clarity-db:password@clarity-db.postgres.database.azure.com:5432/clarity
```

#### 3. Update Backend
File: `backend/.env`
```env
DB_TYPE=postgres
CLOUD_PROVIDER=azure
DB_HOST=clarity-db.postgres.database.azure.com
DB_USER=azureuser@clarity-db
DB_PASSWORD=YourPassword
DB_NAME=clarity
```

---

## Frontend-Backend Connection

### Update Flutter App to Connect to Real Backend

File: `frontend/lib/services/grpc_client.dart` (create new file)

```dart
import 'package:grpc/grpc.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

class GrpcClient {
  static final GrpcClient _instance = GrpcClient._internal();

  late ClientChannel _channel;

  factory GrpcClient() {
    return _instance;
  }

  GrpcClient._internal();

  Future<void> initialize() async {
    final host = dotenv.env['GRPC_HOST'] ?? 'localhost';
    final port = int.parse(dotenv.env['GRPC_PORT'] ?? '50051');
    final useSSL = dotenv.env['GRPC_USE_SSL'] == 'true';

    _channel = ClientChannel(
      host,
      port: port,
      options: ChannelOptions(
        credentials: useSSL
            ? const ChannelCredentials.secure()
            : const ChannelCredentials.insecure(),
      ),
    );
  }

  ClientChannel get channel => _channel;

  Future<void> close() async {
    await _channel.shutdown();
  }
}
```

### Update AuthProvider to Use gRPC

File: `frontend/lib/providers/auth_provider.dart`

```dart
import 'package:clarity/services/grpc_client.dart';
// Import generated gRPC code
// import 'package:clarity/gen/auth.pb.dart';
// import 'package:clarity/gen/auth.pbgrpc.dart';

Future<bool> sendOTP(String email) async {
  try {
    _isLoading = true;
    notifyListeners();

    final grpcClient = GrpcClient();
    // final authService = AuthServiceClient(grpcClient.channel);
    //
    // final response = await authService.sendOTP(
    //   SendOTPRequest(email: email),
    // );

    // For now, mock implementation
    await Future.delayed(Duration(seconds: 2));

    _error = null;
    _isLoading = false;
    notifyListeners();
    return true;
  } catch (e) {
    _error = e.toString();
    _isLoading = false;
    notifyListeners();
    return false;
  }
}
```

### Generate Flutter gRPC Code

```bash
# Install protoc compiler and plugins
flutter pub global activate protoc_plugin

# Generate Dart code from proto files
protoc --dart_out=grpc:lib/gen \
  --plugin=protoc-gen-dart=$(which protoc-gen-dart) \
  -I../backend/proto \
  ../backend/proto/*.proto
```

---

## Testing & Debugging

### Test OTP Service
```bash
grpcurl -plaintext -d '{"email":"test@example.com"}' \
  localhost:50051 clarity.auth.AuthService/SendOTP
```

### Test Health Records Service
```bash
grpcurl -plaintext -d '{
  "user_id":"user123",
  "record_type":"prescription",
  "title":"Aspirin",
  "description":"500mg tablets",
  "metadata":{"dosage":"500mg","frequency":"2x daily"}
}' localhost:50051 clarity.health.HealthRecordsService/CreateRecord
```

### Test AI Service
```bash
# Note: Prescription scanning requires binary image data
grpcurl -plaintext -d '{
  "user_id":"user123",
  "days":7
}' localhost:50051 clarity.ai.AIService/SummarizeHealth
```

### Debug with grpcui
```bash
# Install gRPC UI
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# Run web interface
grpcui -plaintext localhost:50051
# Opens at http://localhost:7469
```

### Test Frontend Connection
File: `frontend/lib/main.dart`

```dart
void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize gRPC client
  await GrpcClient().initialize();

  runApp(MyApp());
}
```

---

## Deployment Checklist

- [ ] Set up cloud database
- [ ] Configure AI service API keys
- [ ] Set up email service
- [ ] Update environment variables
- [ ] Generate gRPC code for frontend
- [ ] Test gRPC connections
- [ ] Set JWT_SECRET to secure value
- [ ] Enable HTTPS/SSL
- [ ] Configure CORS if needed
- [ ] Set up monitoring/logging
- [ ] Deploy backend to container service
- [ ] Build and deploy mobile apps
- [ ] Configure CI/CD pipeline

---

## Troubleshooting

### gRPC Connection Error
```
Error: failed to dial server at localhost:50051
```
**Solution**: Ensure backend is running and port is open

### Database Connection Error
```
Error: connection refused
```
**Solution**: Check DB credentials, firewall rules, and network connectivity

### AI API Rate Limit
```
Error: rate limit exceeded
```
**Solution**: Implement request queuing and exponential backoff

### Email Not Sending
```
Error: authentication failed
```
**Solution**: Verify API key, sender email, and service quotas

---

**Happy integrating! 🚀**
