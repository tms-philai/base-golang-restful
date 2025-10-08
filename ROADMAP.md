# Roadmap: Base Golang RESTful API - Enterprise Features

## Tổng Hợp Base Golang - Tương Đương NestJS Bottle Keeping API

| Mục | Chi Tiết | Logic/Mô Tả | Trạng Thái | Priority |
|-----|----------|-------------|------------|----------|
| 🏗️ **Architecture** | Repository Pattern + Service Layer | Interface → Repository → Service → Handler | ✅ Completed | - |
| 🗄️ **Database** | PostgreSQL + GORM | Entity Relations, Query Builder, Raw SQL, Cloud SQL support | 🔄 Phase 1 | HIGH |
| 🌐 **i18n** | Japanese/English Support | go-i18n, fallback to Japanese | 📋 Phase 2 | HIGH |
| 🔐 **Authentication** | JWT + bcrypt | Token-based auth, password hashing | ✅ Completed | - |
| 📱 **User Types** | App Users + Shop Staff + Admin | Role-based access control (RBAC) | 📋 Phase 2 | HIGH |
| 🏪 **Core Entities** | User, Shop, Tag, UserTag, Device | Many-to-many relationships | 📋 Phase 3 | MEDIUM |
| 📊 **Pagination** | Standard pattern | limit, offset, totalCount | 📋 Phase 1 | HIGH |
| 🔍 **Search** | PostgreSQL ILIKE | Case-insensitive search with %pattern% | 📋 Phase 1 | MEDIUM |
| 🇯🇵 **Sorting** | Japanese-English Mixed | Japanese first, English last with CASE WHEN | 📋 Phase 2 | LOW |
| ⚡ **Performance** | Query Optimization | Proper SELECT fields, efficient JOINs | 📋 Phase 3 | MEDIUM |
| 🛡️ **Error Handling** | Localized Exceptions | Custom errors with i18n messages | 📋 Phase 2 | HIGH |
| 📝 **Logging** | Request/Response Logging | Structured logging with zerolog/zap | 📋 Phase 1 | HIGH |
| 🔄 **Transactions** | GORM Transaction Support | Complex operations with rollback | 📋 Phase 1 | HIGH |
| 📦 **Response Format** | Consistent Structure | data, metadata, error fields | 📋 Phase 1 | HIGH |
| 🎯 **Validation** | DTO + Validator | Input validation with struct tags | ✅ Completed | - |
| 📧 **Email** | SMTP Integration | Verification, notifications | 📋 Phase 4 | MEDIUM |
| ⏰ **Scheduling** | Cron Jobs | Automated tasks with gocron | 📋 Phase 4 | LOW |
| 🔔 **Notifications** | Firebase + Email | Push notifications, email alerts | 📋 Phase 5 | LOW |
| 📊 **Analytics** | Tag Usage Tracking | User behavior, shop analytics | 📋 Phase 5 | LOW |
| 🏷️ **Tag Management** | Physical Tag System | BLE devices, RFID tracking | 📋 Phase 6 | LOW |
| 💾 **Storage** | File Upload Support | Images, documents with GCS/S3/MinIO | 📋 Phase 3 | MEDIUM |
| 🌍 **Environment** | Multi-environment Config | Development, Staging, Production | 📋 Phase 1 | HIGH |
| 🐳 **Deployment** | Docker Support | Containerized deployment with GKE/Cloud Run | 📋 Phase 4 | MEDIUM |
| 📚 **Documentation** | Swagger/OpenAPI | Auto-generated API docs | ✅ Completed | - |
| 🧪 **Testing** | Unit + Integration Tests | Testify framework | ✅ Completed | - |

---

## Phase 1: Core Infrastructure (Week 1-2) 🚀

### 1.1 Database Integration - PostgreSQL + GORM
**Files to create:**
```
app/
├── database/
│   ├── connection.go          # Database connection setup
│   ├── migration.go            # Migration runner
│   └── seeder.go              # Database seeder
├── repository/
│   ├── base_repository.go     # Base repository interface
│   ├── user_repository.go     # User repository implementation
│   └── product_repository.go  # Product repository implementation
```

**Dependencies:**
```go
// go.mod additions
gorm.io/gorm v1.25.5
gorm.io/driver/postgres v1.5.4
github.com/google/uuid v1.5.0
```

**Implementation Tasks:**
- [ ] Setup GORM connection with connection pooling
- [ ] Create base repository interface with CRUD operations
- [ ] Implement repository pattern for existing entities
- [ ] Add migration system
- [ ] Create database seeder for test data

---

### 1.2 Pagination System
**Files to create:**
```
app/
├── utils/
│   ├── pagination.go          # Pagination helper
│   └── query_builder.go       # Query builder utilities
├── models/
│   └── pagination.go          # Pagination request/response models
```

**Implementation:**
```go
type PaginationRequest struct {
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"pageSize" binding:"min=1,max=100"`
    SortBy   string `form:"sortBy"`
    Order    string `form:"order" binding:"oneof=asc desc"`
}

type PaginationResponse struct {
    Data       interface{} `json:"data"`
    TotalCount int64       `json:"totalCount"`
    Page       int         `json:"page"`
    PageSize   int         `json:"pageSize"`
    TotalPages int         `json:"totalPages"`
}
```

---

### 1.3 Structured Logging
**Files to create:**
```
app/
├── logger/
│   ├── logger.go              # Logger setup
│   ├── middleware.go          # Request logging middleware
│   └── fields.go              # Structured fields
```

**Dependencies:**
```go
github.com/rs/zerolog v1.31.0
```

**Implementation Tasks:**
- [ ] Setup zerolog with JSON output
- [ ] Create request/response logging middleware
- [ ] Add correlation ID tracking
- [ ] Implement log levels (debug, info, warn, error)
- [ ] Add file rotation support

---

### 1.4 Transaction Support
**Files to create:**
```
app/
├── database/
│   └── transaction.go         # Transaction manager
├── services/
│   └── base_service.go        # Base service with transaction support
```

**Implementation:**
```go
type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}
```

---

### 1.5 Consistent Response Format
**Files to create:**
```
app/
├── models/
│   ├── response.go            # Standard response structure
│   └── error_response.go      # Error response structure
├── utils/
│   └── response_helper.go     # Response helper functions
```

**Implementation:**
```go
type APIResponse struct {
    Success  bool        `json:"success"`
    Data     interface{} `json:"data,omitempty"`
    Metadata interface{} `json:"metadata,omitempty"`
    Error    *ErrorInfo  `json:"error,omitempty"`
    Message  string      `json:"message,omitempty"`
}

type ErrorInfo struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

---

### 1.6 Environment Configuration
**Files to create:**
```
app/
├── config/
│   ├── database.go            # Database config
│   ├── server.go              # Server config
│   ├── jwt.go                 # JWT config
│   └── loader.go              # Config loader
├── .env.example               # Environment template
├── .env.development           # Development config
├── .env.staging               # Staging config
└── .env.production            # Production config
```

**Dependencies:**
```go
github.com/spf13/viper v1.18.2
github.com/joho/godotenv v1.5.1
```

---

## Phase 2: i18n & RBAC (Week 3-4) 🌐

### 2.1 Internationalization (i18n)
**Files to create:**
```
app/
├── i18n/
│   ├── i18n.go                # i18n setup
│   ├── middleware.go          # Language detection middleware
│   └── locales/
│       ├── ja.json            # Japanese translations
│       └── en.json            # English translations
├── utils/
│   └── translator.go          # Translation helper
```

**Dependencies:**
```go
github.com/nicksnyder/go-i18n/v2 v2.4.0
golang.org/x/text v0.14.0
```

**Implementation Tasks:**
- [ ] Setup go-i18n with JSON translation files
- [ ] Create language detection middleware (Accept-Language header)
- [ ] Implement fallback to Japanese
- [ ] Add translation helper functions
- [ ] Translate all error messages
- [ ] Translate validation messages

---

### 2.2 Role-Based Access Control (RBAC)
**Files to create:**
```
app/
├── models/
│   ├── role.go                # Role model
│   ├── permission.go          # Permission model
│   └── user_role.go           # User-Role relationship
├── middleware/
│   ├── rbac.go                # RBAC middleware
│   └── permission.go          # Permission checker
├── services/
│   ├── role_service.go        # Role service
│   └── permission_service.go  # Permission service
├── handlers/
│   └── role_handler.go        # Role management handler
```

**Implementation:**
```go
type Role struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key"`
    Name        string    `gorm:"unique;not null"`
    Description string
    Permissions []Permission `gorm:"many2many:role_permissions"`
}

const (
    RoleAppUser   = "app_user"
    RoleShopStaff = "shop_staff"
    RoleAdmin     = "admin"
)
```

---

### 2.3 Localized Error Handling
**Files to create:**
```
app/
├── errors/
│   ├── app_error.go           # Custom error types
│   ├── error_codes.go         # Error code constants
│   └── error_handler.go       # Global error handler
├── middleware/
│   └── error_middleware.go    # Error handling middleware
```

**Implementation:**
```go
type AppError struct {
    Code       string
    Message    string
    StatusCode int
    Details    interface{}
}

const (
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeForbidden       = "FORBIDDEN"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeInternalServer  = "INTERNAL_SERVER_ERROR"
)
```

---

## Phase 3: Advanced Features (Week 5-6) ⚡

### 3.1 Search Functionality
**Files to create:**
```
app/
├── utils/
│   ├── search.go              # Search helper
│   └── filter.go              # Filter builder
├── models/
│   └── search_request.go      # Search request model
```

**Implementation:**
```go
type SearchRequest struct {
    Query    string                 `form:"q"`
    Filters  map[string]interface{} `form:"filters"`
    Page     int                    `form:"page"`
    PageSize int                    `form:"pageSize"`
}

// PostgreSQL ILIKE search
func (r *UserRepository) Search(query string) ([]User, error) {
    var users []User
    err := r.db.Where("name ILIKE ? OR email ILIKE ?", 
        "%"+query+"%", "%"+query+"%").Find(&users).Error
    return users, err
}
```

---

### 3.2 Core Entities - Shop System
**Files to create:**
```
app/
├── models/
│   ├── shop.go                # Shop model
│   ├── tag.go                 # Tag model
│   ├── user_tag.go            # UserTag model
│   └── device.go              # Device model
├── repository/
│   ├── shop_repository.go
│   ├── tag_repository.go
│   ├── user_tag_repository.go
│   └── device_repository.go
├── services/
│   ├── shop_service.go
│   ├── tag_service.go
│   └── device_service.go
├── handlers/
│   ├── shop_handler.go
│   ├── tag_handler.go
│   └── device_handler.go
```

**Models:**
```go
type Shop struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key"`
    Name        string    `gorm:"not null"`
    Description string
    Address     string
    Phone       string
    Email       string
    OwnerID     uuid.UUID
    Owner       User      `gorm:"foreignKey:OwnerID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Tag struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key"`
    Name        string    `gorm:"not null"`
    Type        string    // BLE, RFID, NFC
    DeviceID    string    `gorm:"unique"`
    ShopID      uuid.UUID
    Shop        Shop      `gorm:"foreignKey:ShopID"`
    Users       []User    `gorm:"many2many:user_tags"`
    CreatedAt   time.Time
}

type UserTag struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    UserID    uuid.UUID
    TagID     uuid.UUID
    User      User      `gorm:"foreignKey:UserID"`
    Tag       Tag       `gorm:"foreignKey:TagID"`
    AssignedAt time.Time
    ExpiresAt  *time.Time
}
```

---

### 3.3 File Upload & Storage
**Files to create:**
```
app/
├── storage/
│   ├── storage.go             # Storage interface
│   ├── local.go               # Local storage
│   ├── gcs.go                 # Google Cloud Storage
│   ├── s3.go                  # AWS S3 storage (optional)
│   └── minio.go               # MinIO storage (optional)
├── handlers/
│   └── upload_handler.go      # File upload handler
├── middleware/
│   └── upload.go              # Upload middleware
```

**Dependencies:**
```go
cloud.google.com/go/storage v1.36.0
google.golang.org/api v0.154.0
github.com/aws/aws-sdk-go v1.49.0 // optional
github.com/minio/minio-go/v7 v7.0.66 // optional
```

---

### 3.4 Query Optimization
**Files to create:**
```
app/
├── database/
│   ├── query_optimizer.go     # Query optimization utilities
│   └── indexes.go             # Index definitions
```

**Implementation Tasks:**
- [ ] Add database indexes for frequently queried fields
- [ ] Implement eager loading for related entities
- [ ] Add query result caching
- [ ] Optimize N+1 query problems
- [ ] Add database query logging

---

## Phase 4: Communication & Scheduling (Week 7-8) 📧

### 4.1 Email Integration
**Files to create:**
```
app/
├── email/
│   ├── mailer.go              # Email sender
│   ├── templates/
│   │   ├── verification.html  # Email verification template
│   │   ├── reset_password.html
│   │   └── notification.html
│   └── queue.go               # Email queue
├── services/
│   └── email_service.go       # Email service
```

**Dependencies:**
```go
github.com/go-mail/mail v2.3.1
github.com/matcornic/hermes/v2 v2.1.0
```

---

### 4.2 Cron Jobs & Scheduling
**Files to create:**
```
app/
├── scheduler/
│   ├── scheduler.go           # Scheduler setup
│   ├── jobs/
│   │   ├── cleanup_job.go     # Cleanup expired data
│   │   ├── analytics_job.go   # Generate analytics
│   │   └── notification_job.go # Send scheduled notifications
│   └── registry.go            # Job registry
```

**Dependencies:**
```go
github.com/go-co-op/gocron v1.37.0
```

---

### 4.3 Docker Support
**Files to create:**
```
Dockerfile
docker-compose.yml
docker-compose.dev.yml
docker-compose.prod.yml
.dockerignore
scripts/
├── docker-entrypoint.sh
└── wait-for-postgres.sh
```

**Docker Compose Services:**
- API service
- PostgreSQL database (or Cloud SQL proxy)
- Redis cache (or Cloud Memorystore)
- Google Cloud Storage (or MinIO for local dev)
- Nginx reverse proxy (or Cloud Load Balancer)

---

## Phase 5: Advanced Features (Week 9-10) 🔔

### 5.1 Push Notifications - Firebase
**Files to create:**
```
app/
├── notifications/
│   ├── firebase.go            # Firebase setup
│   ├── push_notification.go   # Push notification sender
│   └── device_token.go        # Device token management
├── models/
│   └── notification.go        # Notification model
├── services/
│   └── notification_service.go
├── handlers/
│   └── notification_handler.go
```

**Dependencies:**
```go
firebase.google.com/go/v4 v4.13.0
```

---

### 5.2 Analytics & Tracking
**Files to create:**
```
app/
├── analytics/
│   ├── tracker.go             # Event tracker
│   ├── aggregator.go          # Data aggregator
│   └── reporter.go            # Report generator
├── models/
│   ├── analytics_event.go     # Event model
│   └── analytics_report.go    # Report model
├── services/
│   └── analytics_service.go
├── handlers/
│   └── analytics_handler.go
```

**Features:**
- Tag usage tracking
- User behavior analytics
- Shop performance metrics
- Custom event tracking

---

## Phase 6: IoT Integration (Week 11-12) 🏷️

### 6.1 Physical Tag Management
**Files to create:**
```
app/
├── iot/
│   ├── ble_handler.go         # Bluetooth Low Energy
│   ├── rfid_handler.go        # RFID handler
│   ├── nfc_handler.go         # NFC handler
│   └── device_manager.go      # Device management
├── services/
│   └── tag_tracking_service.go
```

**Dependencies:**
```go
github.com/go-ble/ble v0.0.0-20230130210458-dd4b07d15402
```

---

## Additional Improvements 🎯

### Caching Layer
**Files to create:**
```
app/
├── cache/
│   ├── redis.go               # Redis cache
│   ├── memorystore.go         # Google Cloud Memorystore
│   ├── memory.go              # In-memory cache
│   └── cache_manager.go       # Cache manager
```

**Dependencies:**
```go
github.com/redis/go-redis/v9 v9.4.0
github.com/patrickmn/go-cache v2.1.0
cloud.google.com/go/redis v1.13.0 // for Cloud Memorystore
```

---

### Rate Limiting
**Files to create:**
```
app/
├── middleware/
│   └── rate_limiter.go        # Rate limiting middleware
```

**Dependencies:**
```go
github.com/ulule/limiter/v3 v3.11.2
```

---

### API Versioning
**Files to create:**
```
app/
├── api/
│   ├── v1/
│   │   └── routes.go
│   └── v2/
│       └── routes.go
```

---

### Health Checks
**Files to create:**
```
app/
├── health/
│   ├── health.go              # Health check handler
│   ├── database.go            # Database health
│   ├── redis.go               # Redis health
│   └── external.go            # External services health
```

---

## Testing Strategy 🧪

### Unit Tests
- Service layer tests with mocks
- Repository tests with test database
- Utility function tests

### Integration Tests
- API endpoint tests
- Database integration tests
- External service integration tests

### E2E Tests
- Complete user flows
- Authentication flows
- RBAC scenarios

---

## Performance Targets 🎯

- API response time: < 100ms (p95)
- Database query time: < 50ms (p95)
- Concurrent requests: 1000+ req/s
- Memory usage: < 512MB under normal load
- CPU usage: < 50% under normal load

---

## Security Checklist 🔒

- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS protection
- [ ] CSRF protection
- [ ] Rate limiting
- [ ] JWT token expiration
- [ ] Password hashing with bcrypt
- [ ] HTTPS only in production
- [ ] Secure headers middleware
- [ ] API key authentication for external services
- [ ] Role-based access control
- [ ] Audit logging

---

## Monitoring & Observability 📊

**Tools to integrate:**
- Prometheus metrics
- Grafana dashboards
- ELK stack for log aggregation
- Sentry for error tracking
- APM (Application Performance Monitoring)

---

## Documentation 📚

**To create:**
- API documentation (Swagger/OpenAPI)
- Architecture documentation
- Deployment guide
- Development setup guide
- Contributing guidelines
- Code style guide

---

## Timeline Summary

| Phase | Duration | Focus |
|-------|----------|-------|
| Phase 1 | Week 1-2 | Core Infrastructure |
| Phase 2 | Week 3-4 | i18n & RBAC |
| Phase 3 | Week 5-6 | Advanced Features |
| Phase 4 | Week 7-8 | Communication & Scheduling |
| Phase 5 | Week 9-10 | Notifications & Analytics |
| Phase 6 | Week 11-12 | IoT Integration |

**Total: 12 weeks for complete implementation**

---

## Getting Started

1. Review current codebase structure
2. Set up development environment
3. Create feature branches for each phase
4. Implement features incrementally
5. Write tests alongside implementation
6. Document as you go
7. Code review before merging
8. Deploy to staging for testing
9. Production deployment after approval

---

## Notes

- Prioritize phases based on business requirements
- Some features can be implemented in parallel
- Maintain backward compatibility
- Follow Go best practices and idioms
- Keep dependencies minimal and up-to-date
- Regular security audits
- Performance testing at each phase
