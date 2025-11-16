# Clarity Architecture Documentation

## System Overview

Clarity is a modern Personal Health Records (PHR) application built with a clean separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Flutter Frontend                          │
│                   (Mobile App - iOS/Android)                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ gRPC (Protocol Buffers)
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Go Backend                                │
│              (gRPC Server on port 50051)                     │
├──────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Auth        │  │  Health      │  │  AI          │      │
│  │  Service     │  │  Service     │  │  Service     │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
├──────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐│
│  │         Database Abstraction Layer                       ││
│  │  (Supports SQLite, PostgreSQL, Cloud databases)          ││
│  └──────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        ▼              ▼              ▼
    SQLite (Local) AWS RDS       GCP Cloud SQL
                              (Azure Database)
```

## Frontend Architecture (Flutter)

### State Management: Provider Pattern

```
┌─────────────────────────────────────────┐
│          MultiProvider (Root)           │
├─────────────────────────────────────────┤
│  ├─ AuthProvider                        │
│  │   ├─ isAuthenticated                 │
│  │   ├─ accessToken, refreshToken       │
│  │   ├─ userId, userEmail               │
│  │   └─ Methods: sendOTP, verifyOTP     │
│  │                                       │
│  ├─ AppStateProvider                    │
│  │   ├─ locale (i18n)                   │
│  │   ├─ isDarkMode                      │
│  │   ├─ notificationsEnabled            │
│  │   └─ Methods: setLocale, toggleDarkMode│
│  │                                       │
│  └─ HealthProvider                      │
│      ├─ records list                    │
│      ├─ isLoading                       │
│      └─ Methods: fetchRecords, create, delete│
└─────────────────────────────────────────┘
```

### Screen Hierarchy

```
MaterialApp
├── OnboardingScreen
│   └── Carousel of intro pages
├── LoginScreen
│   ├── Email input
│   └── Send OTP button
├── OTPScreen
│   ├── 6-digit PIN entry
│   └── Resend countdown
└── HomeScreen (BottomNavigationBar)
    ├── Tab 0: DashboardScreen
    │   ├── User greeting
    │   ├── Quick actions grid
    │   └── Recent records
    ├── Tab 1: RecordsScreen
    │   ├── Filter chips (prescription, appointment, etc.)
    │   ├── Records list
    │   └── FAB: Add record
    ├── Tab 2: ScannerScreen
    │   ├── Image picker
    │   ├── Camera capture
    │   └── Results display
    ├── Tab 3: DoctorChatScreen
    │   ├── Chat history
    │   └── Input field + Send
    └── Tab 4: SettingsScreen
        ├── Language selector
        ├── Dark mode toggle
        ├── Notifications toggle
        └── Logout button
```

### Responsive Design Strategy

- **Mobile** (< 600px width): Single column, optimized touch
- **Tablet** (600px - 900px): 2-column layout where applicable
- **Desktop** (> 900px): Full width with constraints

```dart
bool isMobile = MediaQuery.of(context).size.width < 600;

// Adjust padding
EdgeInsets.all(isMobile ? 16 : 24)

// Adjust grid
GridView(
  gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
    crossAxisCount: isMobile ? 2 : 4,
  ),
)
```

### Localization (i18n)

Translation files in `assets/translations/`:
```json
{
  "key": "value",
  "parameterized_key": "Text with {param}"
}
```

Usage:
```dart
context.translate('key')
context.translate('key', {'param': 'value'})
```

Supported locales: English (en), Spanish (es)

## Backend Architecture (Go)

### Service Layer

```
┌─────────────────────────────────────────┐
│         gRPC Service Handlers           │
├─────────────────────────────────────────┤
│  ├─ AuthServer                          │
│  │   ├─ SendOTP                         │
│  │   ├─ VerifyOTP                       │
│  │   └─ RefreshToken                    │
│  │                                       │
│  ├─ HealthRecordsServer                 │
│  │   ├─ CreateRecord                    │
│  │   ├─ GetRecord                       │
│  │   ├─ ListRecords                     │
│  │   ├─ UpdateRecord                    │
│  │   └─ DeleteRecord                    │
│  │                                       │
│  └─ AIServer                            │
│      ├─ ScanPrescription                │
│      ├─ SummarizeHealth                 │
│      └─ DoctorChat (streaming)          │
└─────────────────────────────────────────┘
```

### Service Implementation

```go
// handlers/grpc_handlers.go
type AuthServer struct {
    authService *services.AuthService
}

// services/auth_service.go
type AuthService struct {
    db     *gorm.DB
    config *config.AuthConfig
}

// database/db.go
type Database interface {
    GetConnection() *gorm.DB
    Migrate() error
    Close() error
}
```

### Data Models

```go
// User: Core user information
User {
    ID, Email, Name, DateOfBirth, Gender, BloodType
}

// OTPStore: Temporary OTP storage
OTPStore {
    ID, Email, OTP, ExpiresAt
}

// HealthRecord: Any health data
HealthRecord {
    ID, UserID, RecordType, Title, Description, Metadata
}

// DoctorConversation: Chat history
DoctorConversation {
    ID, UserID, ConversationID, Message, Response, IsAI
}
```

### gRPC Communication Flow

```
1. Client (Flutter) → gRPC Request
                ↓
2. gRPC Handler receives request
                ↓
3. Handler calls Service layer
                ↓
4. Service queries Database
                ↓
5. Database returns result
                ↓
6. Service processes result
                ↓
7. Handler returns gRPC Response
                ↓
8. Client (Flutter) ← gRPC Response
```

## Database Layer Architecture

### Multi-Cloud Support

```
┌──────────────────────────────────────┐
│   Database Interface (Abstract)      │
├──────────────────────────────────────┤
│  interface Database {                │
│    GetConnection() *gorm.DB          │
│    Migrate() error                   │
│    Close() error                     │
│  }                                    │
└──────────────────────────────────────┘
        │         │         │
        ▼         ▼         ▼
    SQLite    RDS       Cloud SQL
   (Local)   (AWS)      (GCP/Azure)
```

### Configuration-Driven Selection

```go
// config/config.go
type DatabaseConfig struct {
    Type            string // "sqlite", "postgres", "mysql"
    CloudProvider   string // "local", "aws", "gcp", "azure"
    Path            string // For SQLite
    Host, Port, DB  string // For remote databases
}

// database/db.go
func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
    switch cfg.Type {
    case "sqlite":
        return newSQLiteDB(cfg)
    // Add more cases for postgres, mysql, etc.
    }
}
```

## Authentication Flow

```
┌─────────────────────────────────────────────┐
│  User enters email & taps "Send OTP"        │
└────────────┬────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│  SendOTP() gRPC call                        │
├─────────────────────────────────────────────┤
│  1. Generate random 6-digit OTP             │
│  2. Set expiry (10 minutes)                 │
│  3. Store in OTPStore table                 │
│  4. Send to email (mock: log to console)    │
└────────────┬────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│  User receives OTP & enters code            │
└────────────┬────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│  VerifyOTP() gRPC call                      │
├─────────────────────────────────────────────┤
│  1. Lookup OTPStore record                  │
│  2. Verify OTP matches                      │
│  3. Check expiry                            │
│  4. Get or create User                      │
│  5. Generate JWT tokens                     │
│  6. Delete used OTP                         │
└────────────┬────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│  Return access_token + refresh_token        │
│  Store in SharedPreferences                 │
│  Navigate to HomeScreen                     │
└─────────────────────────────────────────────┘
```

## API Design Principles

1. **Request/Response**: Protocol buffers for type safety
2. **Stateless**: Backend doesn't maintain session state
3. **Versioning**: Built into proto package names
4. **Error Handling**: Structured error messages
5. **Pagination**: limit/offset pattern for large datasets

## Performance Considerations

### Frontend
- **Provider**: Efficient rebuilds with `Consumer` scoping
- **Image Caching**: Automatic with `Image.file()`
- **List Virtualization**: `ListView.builder()` for long lists
- **Code Splitting**: Lazy load screens as needed

### Backend
- **gRPC**: Binary protocol, ~7x faster than REST JSON
- **Connection Pooling**: gorm manages DB connections
- **Caching**: Potential for Redis in future
- **Indexing**: Indexes on frequently queried columns

## Security Considerations

### Frontend
- **Token Storage**: SharedPreferences (can be upgraded to secure storage)
- **No Hardcoded Secrets**: Use .env file
- **SSL Pinning**: Can be added for production
- **Secure Storage**: Upgrade to `flutter_secure_storage`

### Backend
- **OTP Validation**: Email verification before token issue
- **Token Expiry**: Short-lived access tokens
- **Rate Limiting**: Can be added to prevent brute force
- **Encryption**: Database encryption at rest (cloud provider feature)

## Testing Strategy

### Unit Tests
- Service layer logic
- Helper functions
- Data validation

### Integration Tests
- Database operations
- gRPC service calls
- End-to-end flows

### Widget Tests (Flutter)
- Screen rendering
- User interactions
- Navigation flows

## Deployment Architecture

### Docker Containerization

```dockerfile
# Backend
FROM golang:1.21-alpine
  ├─ Build stage: compile Go binary
  └─ Runtime stage: minimal Alpine image

# Frontend (future)
FROM cirrusci/flutter
  └─ Build APK/IPA, publish to stores
```

### Scaling Strategy

```
┌────────────────┐
│  Load Balancer │
├────────────────┤
│   (Nginx/HAProxy)
└────────┬───────┘
    ┌───┴────┬─────────┐
    ▼        ▼         ▼
┌────────┐┌────────┐┌────────┐
│Backend │ Backend │ Backend │  (Multiple gRPC instances)
└────────┘└────────┘└────────┘
         │
    ┌────▼─────┐
    │ Database  │  (Cloud SQL / RDS cluster)
    └──────────┘
```

## Future Enhancements

### Architecture Level
1. **API Gateway**: Add Kong/Ambassador
2. **Message Queue**: Add RabbitMQ/Kafka for async processing
3. **Caching Layer**: Redis for frequently accessed data
4. **Service Mesh**: Istio for production traffic management
5. **Event Sourcing**: Store health events for audit trail

### AI Integration
1. **LLM Service**: Separate service for AI operations
2. **Vision API**: External prescription scanning
3. **ML Models**: Local models for privacy

### Security Enhancements
1. **OAuth2/OIDC**: Third-party authentication
2. **2FA**: Additional security layer
3. **End-to-End Encryption**: Data encryption at rest
4. **Audit Logging**: Comprehensive logging for compliance

---

**Architecture is designed for simplicity, scalability, and security.**
