# Quick Start Guide

Get Clarity running in minutes!

## Fastest Setup: Docker Compose

```bash
# 1. Clone/navigate to project
cd clarity

# 2. Copy environment file
cp .env.example .env

# 3. Start the backend
docker-compose up -d

# ✅ Backend running at localhost:50051
```

## Manual Setup: Backend

### Requirements
- Go 1.21+
- Protoc (for gRPC)

### Steps

```bash
# 1. Navigate to backend
cd backend

# 2. Install Go dependencies
go mod download

# 3. Generate gRPC code
make proto

# 4. Run the server
make run

# ✅ Server running at localhost:50051
```

## Manual Setup: Frontend

### Requirements
- Flutter 3.0+
- Android SDK or Xcode (for iOS)

### Steps

```bash
# 1. Navigate to frontend
cd frontend

# 2. Install Flutter dependencies
flutter pub get

# 3. Run on emulator/device
flutter run

# ✅ App running!
```

## Testing the Backend

### Using gRPC Client

```bash
# Install grpcurl
brew install grpcurl  # macOS
apt install grpcurl    # Linux

# List services
grpcurl -plaintext localhost:50051 list

# Call SendOTP
grpcurl -plaintext -d '{"email":"test@example.com"}' \
  localhost:50051 clarity.auth.AuthService/SendOTP
```

### Using Postman/Insomnia

1. Import proto files: `backend/proto/*.proto`
2. Connect to `localhost:50051`
3. Call services directly

## Default Credentials

The app uses **OTP authentication**. No default credentials needed:

1. Start app
2. Complete onboarding
3. Enter any email
4. Check backend logs for OTP (mock email)
5. Enter OTP and you're in!

OTP format: 6 digits (logged to console in development)

## Project Structure Overview

```
clarity/
├── backend/          → Go gRPC backend
├── frontend/         → Flutter mobile app
├── docker-compose.yml
├── README.md        → Full documentation
└── QUICK_START.md   → This file
```

## Next Steps

1. **Backend**: Customize AI providers in `config/config.go`
2. **Frontend**: Connect to gRPC backend in `providers/`
3. **Database**: Set up cloud provider in `.env`
4. **Deploy**: Use Docker for cloud deployment

## Common Commands

### Backend
```bash
cd backend
make build      # Compile binary
make run        # Run server
make test       # Run tests
make clean      # Clean build artifacts
make proto      # Generate gRPC code
```

### Frontend
```bash
cd frontend
flutter clean   # Clean build
flutter pub get # Install dependencies
flutter run     # Run app
flutter build apk  # Build Android
flutter build ios  # Build iOS
```

## Environment Variables

Copy `.env.example` to `.env` and customize:

```env
# Backend settings
JWT_SECRET=your-secret-key
AI_PROVIDER=openai
AI_API_KEY=your-api-key

# Cloud provider
CLOUD_PROVIDER=local  # or: aws, gcp, azure
```

## Troubleshooting

### Port 50051 already in use?
```bash
# Find process using port
lsof -i :50051

# Kill process
kill -9 <PID>
```

### Can't build backend?
```bash
# Update Go
go version

# Clean cache
go clean -cache

# Reinstall
go mod download
```

### Flutter not finding emulator?
```bash
flutter emulators  # List emulators
flutter emulators --launch <name>  # Start emulator
flutter run        # Run app
```

## API Testing

### Send OTP
```bash
curl -X POST http://localhost:50051/api/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

### Verify OTP
```bash
curl -X POST http://localhost:50051/api/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","otp":"123456"}'
```

## Support

- 📖 Full docs: See `README.md`
- 🐛 Issues: Check GitHub issues
- 💬 Questions: Create a discussion

---

**You're all set! Start hacking! 🚀**
