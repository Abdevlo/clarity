# Clarity - Personal Health Records App with AI Support

A comprehensive personal health records (PHR) application with AI-powered features including prescription scanning, health summarization, and doctor communication. Built with **Flutter** frontend and **Go backend** using **gRPC**.

## Features

### Core Features
- **OTP-Based Authentication**: Secure email-based OTP verification
- **Health Records Management**: Store and manage prescriptions, appointments, lab results, and symptoms
- **Prescription Scanning**: AI-powered prescription image scanning and extraction
- **Health Summarization**: AI-generated health summaries based on your records
- **Doctor Communication**: Chat with AI doctor or real healthcare professionals
- **Multi-Language Support**: Localization for English (en) and Spanish (es)
- **Responsive Design**: Works seamlessly on mobile and tablet devices
- **Dark Mode**: Built-in dark theme support

### Technical Highlights
- **Cloud Agnostic**: Database abstraction layer supporting AWS, GCP, Azure (currently using SQLite)
- **Simple & Clean**: Minimal complexity, focus on essential features
- **Easy to Run**: Docker Compose setup for quick local development
- **Globalization Ready**: Centralized translation system for easy language changes
- **gRPC**: Efficient binary protocol for frontend-backend communication

## Project Structure

```
clarity/
├── backend/                    # Go gRPC backend
│   ├── config/                # Configuration management
│   ├── database/              # Database abstraction layer
│   ├── handlers/              # gRPC service handlers
│   ├── models/                # Database models
│   ├── proto/                 # Protocol Buffer definitions
│   ├── services/              # Business logic
│   ├── main.go                # Entry point
│   ├── go.mod                 # Go dependencies
│   ├── Dockerfile             # Container configuration
│   ├── Makefile               # Build commands
│   └── .env.example           # Environment template
│
├── frontend/                   # Flutter mobile app
│   ├── lib/
│   │   ├── config/            # Theme, localization, constants
│   │   ├── providers/         # State management (Provider)
│   │   ├── screens/           # UI screens
│   │   │   ├── auth/          # Authentication screens
│   │   │   ├── health/        # Health records, scanner, summary
│   │   │   └── doctor/        # Doctor chat
│   │   └── main.dart          # App entry point
│   ├── assets/
│   │   └── translations/      # i18n JSON files
│   ├── pubspec.yaml           # Flutter dependencies
│   └── android/               # Android-specific config
│
├── docker-compose.yml         # Multi-container orchestration
├── .env.example               # Environment variables template
└── README.md                  # This file
```

## Quick Start

### Prerequisites
- **Backend**: Go 1.21+
- **Frontend**: Flutter 3.0+
- **Docker**: (Optional) For containerized setup

### Setup & Run

#### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd clarity

# Copy environment file
cp .env.example .env

# Start services
docker-compose up -d

# Backend will be available at localhost:50051
```

#### Option 2: Local Development

**Backend:**
```bash
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Generate gRPC code from protobuf
make proto

# Run server
make run

# Server runs on localhost:50051
```

**Frontend:**
```bash
cd frontend

# Install Flutter dependencies
flutter pub get

# Run app on emulator/device
flutter run
```

## Architecture

### Backend (Go)
- **gRPC Services**:
  - `AuthService`: OTP verification, token management
  - `HealthRecordsService`: CRUD operations for health records
  - `AIService`: Prescription scanning, summarization, doctor chat

- **Database Layer**: Abstraction supporting multiple cloud providers
  - Currently: SQLite (local)
  - Future: AWS RDS, Google Cloud SQL, Azure Database

- **Services**:
  - `AuthService`: User authentication and token management
  - `HealthRecordsService`: Health record operations
  - `AIService`: AI-powered features (mock implementations)

### Frontend (Flutter)
- **Providers** (State Management):
  - `AuthProvider`: Authentication state
  - `AppStateProvider`: App-wide settings (theme, language)
  - `HealthProvider`: Health records management

- **Screens**:
  - `OnboardingScreen`: App introduction
  - `LoginScreen`: Email entry
  - `OTPScreen`: OTP verification
  - `HomeScreen`: Dashboard with navigation
  - `RecordsScreen`: Health records list
  - `ScannerScreen`: Prescription scanning
  - `DoctorChatScreen`: AI doctor conversation
  - `SettingsScreen`: App preferences

## Configuration

### Backend Environment Variables

Create `backend/.env`:
```env
# Database
DB_TYPE=sqlite
DB_PATH=./clarity.db
CLOUD_PROVIDER=local

# Server
SERVER_PORT=50051
SERVER_HOST=localhost

# Authentication
JWT_SECRET=your-secret-key
OTP_EXPIRY=600

# AI (e.g., OpenAI)
AI_PROVIDER=openai
AI_API_KEY=sk-xxxxx
```

### Cloud Provider Setup

**AWS**: Set environment variables
```env
CLOUD_PROVIDER=aws
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
AWS_REGION=us-east-1
```

**GCP**: Set credentials path
```env
CLOUD_PROVIDER=gcp
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

**Azure**: Set subscription details
```env
CLOUD_PROVIDER=azure
AZURE_SUBSCRIPTION_ID=xxx
AZURE_RESOURCE_GROUP=xxx
AZURE_TENANT_ID=xxx
```

## Internationalization (i18n)

Add translations in `frontend/assets/translations/`:

**English** (`en.json`): Already included
**Spanish** (`es.json`): Already included

To add a new language:
1. Create `assets/translations/[language-code].json`
2. Add to `pubspec.yaml`:
```yaml
assets:
  - assets/translations/
```
3. Update `lib/config/localization.dart` supported locales

## API Documentation

### gRPC Services

#### AuthService
- `SendOTP(email)`: Send OTP to email
- `VerifyOTP(email, otp)`: Verify OTP and get tokens
- `RefreshToken(refresh_token)`: Refresh access token

#### HealthRecordsService
- `CreateRecord(userId, recordType, title, description, metadata)`: Create record
- `GetRecord(recordId)`: Get single record
- `ListRecords(userId, limit, offset)`: List records with pagination
- `UpdateRecord(recordId, title, description, metadata)`: Update record
- `DeleteRecord(recordId)`: Delete record

#### AIService
- `ScanPrescription(userId, imageData)`: Scan prescription image
- `SummarizeHealth(userId, days)`: Generate health summary
- `DoctorChat(userId, message)`: Chat with AI doctor

## Development

### Adding Backend Services

1. Define protobuf in `backend/proto/yourservice.proto`
2. Generate code: `make proto`
3. Create service handler in `backend/handlers/`
4. Register in `main.go`

### Adding Frontend Screens

1. Create screen in `lib/screens/`
2. Use providers for state management
3. Apply theme using `AppTheme`
4. Add translations to `assets/translations/`
5. Add to navigation in `HomeScreen`

## Testing

### Backend
```bash
cd backend
go test ./... -v
```

### Frontend
```bash
cd frontend
flutter test
```

## Deployment

### Backend
```bash
cd backend
docker build -t clarity-backend:latest .
docker run -p 50051:50051 clarity-backend:latest
```

### Frontend
```bash
cd frontend
flutter build apk     # Android
flutter build ios     # iOS
flutter build web     # Web
```

## Key Design Decisions

1. **OTP Authentication**: Simple, secure, no password management
2. **SQLite by Default**: No external dependencies for quick start
3. **Cloud Abstraction**: Future-proof for enterprise cloud deployments
4. **gRPC**: Efficient, strongly-typed communication
5. **Provider Pattern**: Simple state management without boilerplate
6. **JSON-based i18n**: Easy to maintain and extend translations
7. **Responsive Design**: Single codebase for multiple screen sizes

## Future Enhancements

- [ ] Real AI integration (OpenAI, Google, Hugging Face)
- [ ] Doctor authentication and approval system
- [ ] End-to-end encryption for sensitive data
- [ ] Real-time health notifications
- [ ] Wearable device integration
- [ ] Web dashboard for doctors
- [ ] Production database setup (PostgreSQL, Cloud SQL)
- [ ] Advanced analytics and trends
- [ ] Export health records (PDF, HL7)
- [ ] Integration with healthcare providers

## Contributing

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Commit changes: `git commit -am 'Add feature'`
3. Push to branch: `git push origin feature/my-feature`
4. Submit a pull request

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or suggestions:
1. Open an issue on GitHub
2. Check existing documentation
3. Contact the development team

## Security Considerations

- Always use HTTPS in production
- Rotate JWT_SECRET regularly
- Enable encryption at rest for databases
- Implement rate limiting on OTP endpoints
- Regular security audits recommended
- Use environment variables for secrets (never commit to git)

## Troubleshooting

### Backend won't start
- Check port 50051 is available
- Verify Go version: `go version`
- Check environment variables in `.env`

### Frontend build errors
- Clear Flutter cache: `flutter clean`
- Update Flutter: `flutter upgrade`
- Reinstall dependencies: `flutter pub get`

### gRPC connection issues
- Verify backend is running
- Check firewall settings
- Use `grpcurl` to debug: `grpcurl -plaintext localhost:50051 list`

---

**Made with ❤️ for healthcare accessibility**
