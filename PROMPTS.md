# Implementation Prompts - Base Golang RESTful API

Tài liệu này chứa các prompt chi tiết để implement từng phase trong roadmap.

---

## 📋 PHASE 1: CORE INFRASTRUCTURE

### Prompt 1.1: Database Integration - PostgreSQL + GORM

```
Implement PostgreSQL database integration with GORM for the base Golang RESTful API project:

Requirements:
1. Create database connection setup with connection pooling
2. Implement base repository pattern with generic CRUD operations
3. Add migration system using GORM AutoMigrate
4. Create database seeder for test data
5. Implement health check for database connection

File structure:
app/
├── database/
│   ├── connection.go          # Database connection with pooling
│   ├── migration.go            # Migration runner
│   └── seeder.go              # Database seeder
├── repository/
│   ├── base_repository.go     # Generic repository interface
│   ├── user_repository.go     # User repository implementation
│   └── product_repository.go  # Product repository implementation

Technical details:
- Use GORM v1.25.5 and PostgreSQL driver
- Connection pooling: MaxOpenConns=25, MaxIdleConns=5, ConnMaxLifetime=5min
- Support for transactions
- Proper error handling with custom errors
- Context support for cancellation
- Use UUID for primary keys

Configuration:
- Load from environment variables
- Support for multiple environments (dev, staging, prod)
- Connection retry mechanism with exponential backoff

Please implement following Go best practices and the .cursorrules in the project.
```

---

### Prompt 1.2: Pagination System

```
Implement a comprehensive pagination system for the Golang API:

Requirements:
1. Create pagination request/response models
2. Implement pagination helper functions
3. Add sorting support (asc/desc)
4. Support for multiple sort fields
5. Calculate total pages automatically
6. Add cursor-based pagination option

File structure:
app/
├── models/
│   └── pagination.go          # Pagination models
├── utils/
│   ├── pagination.go          # Pagination helper
│   └── query_builder.go       # Query builder utilities

Features:
- Page-based pagination (page, pageSize)
- Cursor-based pagination for large datasets
- Sorting by multiple fields
- Default values: page=1, pageSize=10
- Max pageSize validation (max 100)
- Total count calculation
- HATEOAS links (first, last, next, prev)

Response format:
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "totalCount": 100,
    "totalPages": 10,
    "hasNext": true,
    "hasPrev": false
  },
  "links": {
    "first": "/api/users?page=1",
    "last": "/api/users?page=10",
    "next": "/api/users?page=2",
    "prev": null
  }
}

Integrate with GORM for database queries.
Follow Go best practices and project .cursorrules.
```

---

### Prompt 1.3: Structured Logging with Zerolog

```
Implement structured logging system using zerolog:

Requirements:
1. Setup zerolog with JSON output
2. Create request/response logging middleware
3. Add correlation ID tracking
4. Implement log levels (debug, info, warn, error, fatal)
5. Add file rotation support
6. Log to both console and file

File structure:
app/
├── logger/
│   ├── logger.go              # Logger setup and configuration
│   ├── middleware.go          # Request logging middleware
│   └── fields.go              # Structured fields helpers

Features:
- Structured JSON logging
- Correlation ID for request tracking
- Request/response logging with duration
- Error stack traces
- Log rotation (daily, max 7 days)
- Different log levels per environment
- Contextual logging with fields
- Performance metrics logging

Middleware logging format:
{
  "level": "info",
  "time": "2024-01-01T10:00:00Z",
  "correlation_id": "uuid",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "duration_ms": 45,
  "ip": "127.0.0.1",
  "user_agent": "...",
  "user_id": "uuid"
}

Configuration:
- Log level from environment
- File path configuration
- Console output for development
- JSON output for production

Follow Go best practices and project .cursorrules.
```

---

### Prompt 1.4: Transaction Support

```
Implement database transaction support with GORM:

Requirements:
1. Create transaction manager interface
2. Implement transaction wrapper functions
3. Add automatic rollback on error
4. Support nested transactions
5. Add transaction context propagation
6. Implement transaction timeout

File structure:
app/
├── database/
│   └── transaction.go         # Transaction manager
├── services/
│   └── base_service.go        # Base service with transaction support

Features:
- WithTransaction helper function
- Automatic commit/rollback
- Context-aware transactions
- Transaction timeout support
- Nested transaction handling
- Transaction isolation levels
- Savepoint support

Usage example:
err := txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // Create user
    if err := userRepo.Create(tx, user); err != nil {
        return err
    }
    
    // Create user profile
    if err := profileRepo.Create(tx, profile); err != nil {
        return err
    }
    
    return nil
})

Error handling:
- Automatic rollback on any error
- Proper error propagation
- Transaction deadlock detection
- Retry mechanism for deadlocks

Follow Go best practices and project .cursorrules.
```

---

### Prompt 1.5: Consistent Response Format

```
Implement consistent API response format:

Requirements:
1. Create standard response structure
2. Implement error response structure
3. Add response helper functions
4. Support for metadata
5. Standardize error codes
6. Add success/failure indicators

File structure:
app/
├── models/
│   ├── response.go            # Standard response structure
│   └── error_response.go      # Error response structure
├── utils/
│   └── response_helper.go     # Response helper functions

Response structures:

Success response:
{
  "success": true,
  "data": {...},
  "metadata": {
    "timestamp": "2024-01-01T10:00:00Z",
    "version": "v1"
  },
  "message": "Operation successful"
}

Error response:
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  },
  "metadata": {
    "timestamp": "2024-01-01T10:00:00Z",
    "correlation_id": "uuid"
  }
}

Helper functions:
- SuccessResponse(data, message)
- ErrorResponse(code, message, details)
- PaginatedResponse(data, pagination)
- CreatedResponse(data)
- NoContentResponse()
- ValidationErrorResponse(errors)

Error codes:
- VALIDATION_ERROR
- UNAUTHORIZED
- FORBIDDEN
- NOT_FOUND
- CONFLICT
- INTERNAL_SERVER_ERROR
- BAD_REQUEST
- SERVICE_UNAVAILABLE

Follow Go best practices and project .cursorrules.
```

---

### Prompt 1.6: Environment Configuration

```
Implement multi-environment configuration system:

Requirements:
1. Setup Viper for configuration management
2. Support for multiple environments (dev, staging, prod)
3. Load from .env files and environment variables
4. Create configuration structs for different components
5. Add configuration validation
6. Support for configuration hot-reload

File structure:
app/
├── config/
│   ├── config.go              # Main config structure
│   ├── database.go            # Database config
│   ├── server.go              # Server config
│   ├── jwt.go                 # JWT config
│   ├── redis.go               # Redis config
│   ├── storage.go             # Storage config
│   └── loader.go              # Config loader
├── .env.example               # Environment template
├── .env.development           # Development config
├── .env.staging               # Staging config
└── .env.production            # Production config

Configuration structure:
type Config struct {
    Environment string
    Server      ServerConfig
    Database    DatabaseConfig
    JWT         JWTConfig
    Redis       RedisConfig
    Storage     StorageConfig
    Email       EmailConfig
    Logging     LoggingConfig
}

Features:
- Environment variable override
- Default values
- Configuration validation
- Sensitive data masking in logs
- Configuration reload without restart
- Type-safe configuration access

Environment files:
.env.development:
- Debug mode enabled
- Verbose logging
- Local database
- Short JWT expiration

.env.production:
- Debug mode disabled
- Error logging only
- Production database with SSL
- Long JWT expiration
- Rate limiting enabled

Validation:
- Required fields check
- Format validation
- Range validation
- Dependency validation

Follow Go best practices and project .cursorrules.
```

---

## 📋 PHASE 2: i18n & RBAC

### Prompt 2.1: Internationalization (i18n)

```
Implement internationalization support with Japanese and English:

Requirements:
1. Setup go-i18n with JSON translation files
2. Create language detection middleware
3. Implement fallback to Japanese
4. Add translation helper functions
5. Translate all error messages
6. Translate validation messages
7. Support for pluralization
8. Support for variable interpolation

File structure:
app/
├── i18n/
│   ├── i18n.go                # i18n setup
│   ├── middleware.go          # Language detection middleware
│   ├── translator.go          # Translation helper
│   └── locales/
│       ├── ja.json            # Japanese translations
│       └── en.json            # English translations
├── utils/
│   └── translator.go          # Translation utilities

Translation files structure:
ja.json:
{
  "errors": {
    "validation": {
      "required": "{{.Field}}は必須です",
      "email": "有効なメールアドレスを入力してください",
      "min_length": "{{.Field}}は{{.Min}}文字以上である必要があります"
    },
    "auth": {
      "unauthorized": "認証が必要です",
      "forbidden": "アクセス権限がありません",
      "invalid_credentials": "メールアドレスまたはパスワードが正しくありません"
    },
    "not_found": "{{.Resource}}が見つかりません"
  },
  "messages": {
    "success": {
      "created": "{{.Resource}}が正常に作成されました",
      "updated": "{{.Resource}}が正常に更新されました",
      "deleted": "{{.Resource}}が正常に削除されました"
    }
  }
}

Features:
- Language detection from Accept-Language header
- Language override via query parameter (?lang=en)
- Fallback to Japanese as default
- Context-aware translations
- Pluralization support
- Variable interpolation
- Date/time formatting per locale
- Number formatting per locale

Middleware:
- Detect language from header
- Store language in context
- Add language to response headers

Helper functions:
- T(key, vars) - Translate with variables
- TPlural(key, count, vars) - Translate with pluralization
- LocalizeError(err) - Localize error messages
- LocalizeValidation(errors) - Localize validation errors

Integration:
- Integrate with validation errors
- Integrate with custom errors
- Integrate with response messages

Follow Go best practices and project .cursorrules.
```

---

### Prompt 2.2: Role-Based Access Control (RBAC)

```
Implement comprehensive RBAC system with three user types:

Requirements:
1. Create role and permission models
2. Implement RBAC middleware
3. Add permission checker functions
4. Create role management service
5. Implement role assignment
6. Add permission caching
7. Support for hierarchical roles

File structure:
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
├── repository/
│   ├── role_repository.go
│   └── permission_repository.go

User types:
1. App User (app_user)
   - Basic user permissions
   - Can view own data
   - Can update own profile
   
2. Shop Staff (shop_staff)
   - Can manage shop data
   - Can view shop analytics
   - Can manage tags
   - Can view customers
   
3. Admin (admin)
   - Full system access
   - Can manage users
   - Can manage roles
   - Can view all data

Models:
type Role struct {
    ID          uuid.UUID
    Name        string
    Description string
    Permissions []Permission
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Permission struct {
    ID          uuid.UUID
    Name        string
    Resource    string
    Action      string
    Description string
}

type UserRole struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    RoleID    uuid.UUID
    User      User
    Role      Role
    AssignedAt time.Time
    ExpiresAt  *time.Time
}

Permission format: "resource:action"
Examples:
- "user:read"
- "user:write"
- "shop:manage"
- "tag:create"
- "analytics:view"

Middleware usage:
router.GET("/admin/users", 
    RequireAuth(),
    RequireRole("admin"),
    RequirePermission("user:read"),
    handler.GetUsers)

Features:
- Role hierarchy (admin > shop_staff > app_user)
- Permission inheritance
- Dynamic permission checking
- Permission caching with Redis
- Role assignment with expiration
- Audit logging for role changes

Helper functions:
- HasRole(userID, role) bool
- HasPermission(userID, permission) bool
- HasAnyRole(userID, roles) bool
- HasAllPermissions(userID, permissions) bool
- GetUserRoles(userID) []Role
- GetUserPermissions(userID) []Permission

Seeder:
- Create default roles
- Assign default permissions
- Create admin user

Follow Go best practices and project .cursorrules.
```

---

### Prompt 2.3: Localized Error Handling

```
Implement comprehensive localized error handling system:

Requirements:
1. Create custom error types
2. Define error code constants
3. Implement global error handler
4. Add error handling middleware
5. Integrate with i18n system
6. Add error logging
7. Support for error details

File structure:
app/
├── errors/
│   ├── app_error.go           # Custom error types
│   ├── error_codes.go         # Error code constants
│   ├── error_handler.go       # Global error handler
│   └── error_factory.go       # Error factory functions
├── middleware/
│   └── error_middleware.go    # Error handling middleware

Custom error structure:
type AppError struct {
    Code       string
    Message    string
    MessageKey string
    StatusCode int
    Details    interface{}
    Err        error
    Stack      string
}

Error codes:
const (
    // Validation errors (400)
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeInvalidInput    = "INVALID_INPUT"
    ErrCodeInvalidFormat   = "INVALID_FORMAT"
    
    // Authentication errors (401)
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeInvalidToken    = "INVALID_TOKEN"
    ErrCodeTokenExpired    = "TOKEN_EXPIRED"
    
    // Authorization errors (403)
    ErrCodeForbidden       = "FORBIDDEN"
    ErrCodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
    
    // Not found errors (404)
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
    
    // Conflict errors (409)
    ErrCodeConflict        = "CONFLICT"
    ErrCodeDuplicateEntry  = "DUPLICATE_ENTRY"
    
    // Server errors (500)
    ErrCodeInternalServer  = "INTERNAL_SERVER_ERROR"
    ErrCodeDatabaseError   = "DATABASE_ERROR"
    ErrCodeExternalService = "EXTERNAL_SERVICE_ERROR"
)

Error factory functions:
- NewValidationError(details)
- NewUnauthorizedError()
- NewForbiddenError()
- NewNotFoundError(resource)
- NewConflictError(resource)
- NewInternalError(err)

Error middleware:
- Catch panics and convert to errors
- Log errors with context
- Translate error messages
- Format error response
- Add correlation ID
- Hide sensitive information in production

Response format:
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "メールアドレスは必須です",
    "details": {
      "field": "email",
      "reason": "required"
    }
  },
  "metadata": {
    "timestamp": "2024-01-01T10:00:00Z",
    "correlation_id": "uuid"
  }
}

Features:
- Automatic error translation
- Stack trace in development
- Error aggregation
- Error rate monitoring
- Sentry integration ready
- Custom error details
- HTTP status code mapping

Integration:
- Integrate with i18n system
- Integrate with logging system
- Integrate with monitoring

Follow Go best practices and project .cursorrules.
```

---

## 📋 PHASE 3: ADVANCED FEATURES

### Prompt 3.1: Search Functionality

```
Implement comprehensive search functionality with PostgreSQL:

Requirements:
1. Create search request/response models
2. Implement full-text search with PostgreSQL
3. Add ILIKE pattern matching
4. Support for multiple search fields
5. Implement filter builder
6. Add search result highlighting
7. Support for fuzzy search

File structure:
app/
├── models/
│   ├── search_request.go      # Search request model
│   └── search_response.go     # Search response model
├── utils/
│   ├── search.go              # Search helper
│   ├── filter.go              # Filter builder
│   └── query_builder.go       # Advanced query builder

Search request structure:
type SearchRequest struct {
    Query    string                 `form:"q" json:"query"`
    Fields   []string               `form:"fields" json:"fields"`
    Filters  map[string]interface{} `form:"filters" json:"filters"`
    Page     int                    `form:"page" json:"page"`
    PageSize int                    `form:"pageSize" json:"pageSize"`
    SortBy   string                 `form:"sortBy" json:"sortBy"`
    Order    string                 `form:"order" json:"order"`
}

Features:
1. ILIKE search (case-insensitive):
   - WHERE name ILIKE '%query%'
   - Support for multiple fields
   - OR conditions between fields

2. Full-text search:
   - PostgreSQL tsvector and tsquery
   - Ranking by relevance
   - Language-specific search (Japanese, English)

3. Advanced filters:
   - Exact match: field = value
   - Range: field BETWEEN min AND max
   - IN: field IN (value1, value2)
   - NULL check: field IS NULL
   - Date range: created_at BETWEEN start AND end

4. Filter operators:
   - eq (equal)
   - ne (not equal)
   - gt (greater than)
   - gte (greater than or equal)
   - lt (less than)
   - lte (less than or equal)
   - in (in array)
   - nin (not in array)
   - like (pattern match)
   - ilike (case-insensitive pattern match)

Example usage:
GET /api/users/search?q=john&fields=name,email&filters[status]=active&sortBy=created_at&order=desc

Implementation:
func (r *UserRepository) Search(req SearchRequest) ([]User, int64, error) {
    query := r.db.Model(&User{})
    
    // Apply search
    if req.Query != "" {
        conditions := []string{}
        for _, field := range req.Fields {
            conditions = append(conditions, 
                fmt.Sprintf("%s ILIKE ?", field))
        }
        query = query.Where(
            strings.Join(conditions, " OR "),
            "%"+req.Query+"%")
    }
    
    // Apply filters
    for key, value := range req.Filters {
        query = query.Where(key+" = ?", value)
    }
    
    // Apply sorting
    if req.SortBy != "" {
        query = query.Order(req.SortBy + " " + req.Order)
    }
    
    // Count total
    var total int64
    query.Count(&total)
    
    // Apply pagination
    var users []User
    err := query.
        Offset((req.Page - 1) * req.PageSize).
        Limit(req.PageSize).
        Find(&users).Error
    
    return users, total, err
}

Response format:
{
  "data": [...],
  "search": {
    "query": "john",
    "fields": ["name", "email"],
    "total_results": 45,
    "search_time_ms": 23
  },
  "pagination": {...}
}

Follow Go best practices and project .cursorrules.
```

---

### Prompt 3.2: Core Entities - Shop System

```
Implement complete shop management system with tags and devices:

Requirements:
1. Create models for Shop, Tag, UserTag, Device
2. Implement repositories for all entities
3. Create services with business logic
4. Implement handlers for API endpoints
5. Add many-to-many relationships
6. Support for tag assignment and tracking
7. Add device management

File structure:
app/
├── models/
│   ├── shop.go                # Shop model
│   ├── tag.go                 # Tag model
│   ├── user_tag.go            # UserTag model
│   ├── device.go              # Device model
│   └── tag_history.go         # Tag usage history
├── repository/
│   ├── shop_repository.go
│   ├── tag_repository.go
│   ├── user_tag_repository.go
│   └── device_repository.go
├── services/
│   ├── shop_service.go
│   ├── tag_service.go
│   ├── device_service.go
│   └── tag_tracking_service.go
├── handlers/
│   ├── shop_handler.go
│   ├── tag_handler.go
│   └── device_handler.go

Models:

1. Shop:
type Shop struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name        string    `gorm:"not null;index"`
    NameJa      string    `gorm:"column:name_ja"`
    NameEn      string    `gorm:"column:name_en"`
    Description string
    Address     string
    City        string
    PostalCode  string
    Phone       string
    Email       string    `gorm:"unique"`
    Website     string
    OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
    Owner       User      `gorm:"foreignKey:OwnerID"`
    Status      string    `gorm:"default:'active'"` // active, inactive, suspended
    Tags        []Tag     `gorm:"foreignKey:ShopID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

2. Tag:
type Tag struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name        string    `gorm:"not null;index"`
    Type        string    `gorm:"not null"` // BLE, RFID, NFC, QR
    DeviceID    string    `gorm:"unique;not null"`
    ShopID      uuid.UUID `gorm:"type:uuid;not null"`
    Shop        Shop      `gorm:"foreignKey:ShopID"`
    Status      string    `gorm:"default:'active'"` // active, inactive, lost, damaged
    BatteryLevel int      `gorm:"default:100"`
    LastSeen    *time.Time
    Users       []User    `gorm:"many2many:user_tags"`
    Metadata    datatypes.JSON
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

3. UserTag:
type UserTag struct {
    ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
    TagID      uuid.UUID `gorm:"type:uuid;not null;index"`
    User       User      `gorm:"foreignKey:UserID"`
    Tag        Tag       `gorm:"foreignKey:TagID"`
    AssignedBy uuid.UUID `gorm:"type:uuid"`
    AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    ExpiresAt  *time.Time
    Status     string    `gorm:"default:'active'"` // active, expired, revoked
    UsageCount int       `gorm:"default:0"`
    LastUsedAt *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

4. Device:
type Device struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    DeviceID     string    `gorm:"unique;not null"`
    Type         string    `gorm:"not null"` // reader, scanner, beacon
    Model        string
    Manufacturer string
    ShopID       uuid.UUID `gorm:"type:uuid"`
    Shop         Shop      `gorm:"foreignKey:ShopID"`
    Location     string
    Status       string    `gorm:"default:'online'"` // online, offline, maintenance
    LastPing     *time.Time
    FirmwareVersion string
    Metadata     datatypes.JSON
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

5. TagHistory:
type TagHistory struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    TagID     uuid.UUID `gorm:"type:uuid;not null;index"`
    UserID    uuid.UUID `gorm:"type:uuid;index"`
    ShopID    uuid.UUID `gorm:"type:uuid;index"`
    DeviceID  uuid.UUID `gorm:"type:uuid;index"`
    Action    string    `gorm:"not null"` // scan, assign, revoke, update
    Location  string
    Metadata  datatypes.JSON
    CreatedAt time.Time
}

API Endpoints:

Shops:
- POST   /api/shops                 - Create shop
- GET    /api/shops                 - List shops (paginated, searchable)
- GET    /api/shops/:id             - Get shop details
- PUT    /api/shops/:id             - Update shop
- DELETE /api/shops/:id             - Delete shop
- GET    /api/shops/:id/tags        - Get shop tags
- GET    /api/shops/:id/analytics   - Get shop analytics

Tags:
- POST   /api/tags                  - Create tag
- GET    /api/tags                  - List tags (paginated, filterable)
- GET    /api/tags/:id              - Get tag details
- PUT    /api/tags/:id              - Update tag
- DELETE /api/tags/:id              - Delete tag
- POST   /api/tags/:id/assign       - Assign tag to user
- POST   /api/tags/:id/revoke       - Revoke tag from user
- GET    /api/tags/:id/history      - Get tag usage history
- GET    /api/tags/:id/users        - Get users with this tag

User Tags:
- GET    /api/users/:id/tags        - Get user's tags
- POST   /api/users/:id/tags        - Assign tag to user
- DELETE /api/users/:id/tags/:tagId - Remove tag from user

Devices:
- POST   /api/devices               - Register device
- GET    /api/devices               - List devices
- GET    /api/devices/:id           - Get device details
- PUT    /api/devices/:id           - Update device
- DELETE /api/devices/:id           - Delete device
- POST   /api/devices/:id/ping      - Device heartbeat

Business logic:
1. Tag assignment:
   - Check tag availability
   - Check user eligibility
   - Set expiration if needed
   - Log assignment in history
   - Send notification

2. Tag scanning:
   - Verify tag is active
   - Verify user has access
   - Check expiration
   - Log scan in history
   - Update last seen
   - Trigger analytics

3. Shop management:
   - Only owner can manage shop
   - Admin can manage all shops
   - Staff can view shop data

Follow Go best practices and project .cursorrules.
```

---

### Prompt 3.3: File Upload & Storage

```
Implement file upload and storage system with multiple backends:

Requirements:
1. Create storage interface for multiple backends
2. Implement local file storage
3. Implement Google Cloud Storage (GCS)
4. Implement AWS S3 storage (optional)
5. Implement MinIO storage (optional)
6. Add file validation (type, size)
7. Generate thumbnails for images
8. Support for presigned URLs
9. Add file metadata storage

File structure:
app/
├── storage/
│   ├── storage.go             # Storage interface
│   ├── local.go               # Local file system storage
│   ├── gcs.go                 # Google Cloud Storage
│   ├── s3.go                  # AWS S3 storage (optional)
│   ├── minio.go               # MinIO storage (optional)
│   └── factory.go             # Storage factory
├── models/
│   └── file.go                # File metadata model
├── repository/
│   └── file_repository.go     # File repository
├── services/
│   └── file_service.go        # File service
├── handlers/
│   └── upload_handler.go      # Upload handler
├── middleware/
│   └── upload.go              # Upload middleware
└── utils/
    ├── image.go               # Image processing
    └── file_validator.go      # File validation

Storage interface:
type Storage interface {
    Upload(ctx context.Context, file io.Reader, path string, opts UploadOptions) (string, error)
    Download(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
    GetURL(ctx context.Context, path string, expires time.Duration) (string, error)
    Exists(ctx context.Context, path string) (bool, error)
    List(ctx context.Context, prefix string) ([]FileInfo, error)
}

File model:
type File struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FileName     string    `gorm:"not null"`
    OriginalName string    `gorm:"not null"`
    FilePath     string    `gorm:"not null;unique"`
    FileSize     int64     `gorm:"not null"`
    MimeType     string    `gorm:"not null"`
    Extension    string    `gorm:"not null"`
    Storage      string    `gorm:"not null"` // local, s3, minio
    Bucket       string
    URL          string
    ThumbnailURL string
    Width        int
    Height       int
    Hash         string    `gorm:"index"`
    UploaderID   uuid.UUID `gorm:"type:uuid"`
    Uploader     User      `gorm:"foreignKey:UploaderID"`
    IsPublic     bool      `gorm:"default:false"`
    Metadata     datatypes.JSON
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}

Features:
1. File validation:
   - Allowed file types (images, documents, videos)
   - Max file size (configurable per type)
   - File extension validation
   - MIME type validation
   - Virus scanning (optional)

2. Image processing:
   - Thumbnail generation (multiple sizes)
   - Image optimization
   - Format conversion
   - EXIF data extraction
   - Watermark (optional)

3. Upload options:
   - Public/private access
   - Custom path
   - Custom filename
   - Metadata
   - Expiration time

4. Security:
   - Sanitize filenames
   - Generate unique filenames (UUID)
   - Validate file content
   - Rate limiting
   - User quota

API Endpoints:
- POST   /api/upload                - Upload single file
- POST   /api/upload/multiple       - Upload multiple files
- GET    /api/files                 - List user files
- GET    /api/files/:id             - Get file details
- GET    /api/files/:id/download    - Download file
- DELETE /api/files/:id             - Delete file
- GET    /api/files/:id/url         - Get presigned URL

Configuration:
type StorageConfig struct {
    Provider    string // local, gcs, s3, minio
    Local       LocalConfig
    GCS         GCSConfig
    S3          S3Config
    MinIO       MinIOConfig
    MaxFileSize int64
    AllowedTypes []string
}

type GCSConfig struct {
    ProjectID           string
    Bucket              string
    CredentialsFile     string
    CredentialsJSON     string
    PublicURL           string
    SignedURLExpiration time.Duration
}

type S3Config struct {
    Region          string
    Bucket          string
    AccessKeyID     string
    SecretAccessKey string
    Endpoint        string
    UseSSL          bool
}

Usage example:
// Upload file
file, err := c.FormFile("file")
url, err := fileService.Upload(ctx, file, UploadOptions{
    IsPublic: true,
    GenerateThumbnail: true,
})

// Get presigned URL
url, err := fileService.GetPresignedURL(ctx, fileID, 1*time.Hour)

Follow Go best practices and project .cursorrules.
```

---

## 📋 PHASE 4: COMMUNICATION & SCHEDULING

### Prompt 4.1: Email Integration

```
Implement comprehensive email system with templates:

Requirements:
1. Setup SMTP email sender
2. Create HTML email templates
3. Implement email queue system
4. Add email tracking
5. Support for attachments
6. Implement email verification flow
7. Add password reset emails

File structure:
app/
├── email/
│   ├── mailer.go              # Email sender
│   ├── queue.go               # Email queue
│   ├── tracker.go             # Email tracking
│   ├── templates/
│   │   ├── base.html          # Base template
│   │   ├── verification.html  # Email verification
│   │   ├── reset_password.html
│   │   ├── welcome.html
│   │   ├── notification.html
│   │   └── invoice.html
│   └── builder.go             # Email builder
├── models/
│   ├── email.go               # Email model
│   └── email_log.go           # Email log model
├── services/
│   └── email_service.go       # Email service
└── handlers/
    └── email_handler.go       # Email webhook handler

Email model:
type Email struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    To          string    `gorm:"not null"`
    From        string    `gorm:"not null"`
    Subject     string    `gorm:"not null"`
    Body        string    `gorm:"type:text"`
    BodyHTML    string    `gorm:"type:text"`
    Template    string
    TemplateData datatypes.JSON
    Attachments datatypes.JSON
    Status      string    `gorm:"default:'pending'"` // pending, sent, failed, bounced
    Priority    int       `gorm:"default:5"`
    ScheduledAt *time.Time
    SentAt      *time.Time
    ErrorMessage string
    Attempts    int       `gorm:"default:0"`
    MaxAttempts int       `gorm:"default:3"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

Configuration:
type EmailConfig struct {
    Provider    string // smtp, sendgrid, ses
    SMTP        SMTPConfig
    From        string
    FromName    string
    ReplyTo     string
    MaxRetries  int
    RetryDelay  time.Duration
}

type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    UseTLS   bool
}

Features:
1. Email templates:
   - HTML templates with Go template engine
   - Base layout with header/footer
   - Responsive design
   - Inline CSS
   - Variable interpolation
   - Localization support

2. Email queue:
   - Async email sending
   - Priority queue
   - Retry mechanism
   - Rate limiting
   - Batch sending

3. Email tracking:
   - Delivery status
   - Open tracking
   - Click tracking
   - Bounce handling
   - Unsubscribe handling

4. Email types:
   a. Verification email:
      - Send verification link
      - Token expiration
      - Resend functionality
   
   b. Password reset:
      - Send reset link
      - Token expiration (1 hour)
      - One-time use token
   
   c. Welcome email:
      - Send after registration
      - Include getting started guide
   
   d. Notification email:
      - Tag assignment
      - Tag expiration warning
      - Shop updates
   
   e. Invoice email:
      - PDF attachment
      - Payment details

API Endpoints:
- POST /api/email/send           - Send email
- POST /api/email/verify         - Verify email
- POST /api/email/resend         - Resend verification
- POST /api/email/reset-password - Send password reset
- GET  /api/email/logs           - Get email logs
- POST /api/webhooks/email       - Email webhook (bounces, opens)

Service methods:
- SendVerificationEmail(user)
- SendPasswordResetEmail(user)
- SendWelcomeEmail(user)
- SendNotificationEmail(user, notification)
- SendInvoiceEmail(user, invoice)

Template example (verification.html):
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Greeting}}</h1>
    <p>{{.Message}}</p>
    <a href="{{.VerificationURL}}">{{.ButtonText}}</a>
    <p>{{.Footer}}</p>
</body>
</html>

Usage example:
err := emailService.SendVerificationEmail(ctx, user, EmailOptions{
    Language: "ja",
    ExpiresIn: 24 * time.Hour,
})

Follow Go best practices and project .cursorrules.
```

---

### Prompt 4.2: Cron Jobs & Scheduling

```
Implement task scheduling system with cron jobs:

Requirements:
1. Setup gocron scheduler
2. Create job registry
3. Implement cleanup jobs
4. Add analytics generation jobs
5. Create notification jobs
6. Add job monitoring
7. Support for one-time and recurring jobs

File structure:
app/
├── scheduler/
│   ├── scheduler.go           # Scheduler setup
│   ├── registry.go            # Job registry
│   ├── monitor.go             # Job monitoring
│   └── jobs/
│       ├── cleanup_job.go     # Cleanup expired data
│       ├── analytics_job.go   # Generate analytics
│       ├── notification_job.go # Send scheduled notifications
│       ├── backup_job.go      # Database backup
│       ├── report_job.go      # Generate reports
│       └── sync_job.go        # Sync external data
├── models/
│   └── scheduled_job.go       # Job model
└── services/
    └── job_service.go         # Job service

Job model:
type ScheduledJob struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name        string    `gorm:"unique;not null"`
    Description string
    Schedule    string    `gorm:"not null"` // cron expression
    JobType     string    `gorm:"not null"` // recurring, one-time
    Handler     string    `gorm:"not null"`
    Params      datatypes.JSON
    Status      string    `gorm:"default:'active'"` // active, paused, stopped
    LastRunAt   *time.Time
    NextRunAt   *time.Time
    RunCount    int       `gorm:"default:0"`
    FailCount   int       `gorm:"default:0"`
    LastError   string
    Timeout     int       `gorm:"default:300"` // seconds
    Retries     int       `gorm:"default:3"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

Jobs to implement:

1. Cleanup Job (Daily at 2 AM):
   - Delete expired tokens
   - Delete old logs (> 30 days)
   - Delete expired user tags
   - Clean up temporary files
   - Archive old records

2. Analytics Job (Every hour):
   - Calculate tag usage statistics
   - Generate shop analytics
   - Calculate user engagement metrics
   - Update dashboard data
   - Generate trending reports

3. Notification Job (Every 15 minutes):
   - Send pending notifications
   - Check tag expiration warnings
   - Send reminder emails
   - Process notification queue

4. Backup Job (Daily at 3 AM):
   - Database backup
   - File storage backup
   - Export to S3
   - Cleanup old backups

5. Report Job (Weekly on Monday):
   - Generate weekly reports
   - Send reports to admins
   - Generate shop performance reports
   - User activity reports

6. Sync Job (Every 30 minutes):
   - Sync with external systems
   - Update device status
   - Refresh cache
   - Update exchange rates

Scheduler setup:
type Scheduler struct {
    cron      *gocron.Scheduler
    jobs      map[string]*gocron.Job
    logger    *zerolog.Logger
    db        *gorm.DB
}

func (s *Scheduler) Start() error {
    // Register jobs
    s.RegisterJob("cleanup", "0 2 * * *", CleanupJob)
    s.RegisterJob("analytics", "0 * * * *", AnalyticsJob)
    s.RegisterJob("notifications", "*/15 * * * *", NotificationJob)
    s.RegisterJob("backup", "0 3 * * *", BackupJob)
    s.RegisterJob("reports", "0 9 * * 1", ReportJob)
    s.RegisterJob("sync", "*/30 * * * *", SyncJob)
    
    // Start scheduler
    s.cron.StartAsync()
    return nil
}

Job interface:
type Job interface {
    Name() string
    Execute(ctx context.Context) error
    OnSuccess(ctx context.Context)
    OnError(ctx context.Context, err error)
}

Job implementation example:
type CleanupJob struct {
    db     *gorm.DB
    logger *zerolog.Logger
}

func (j *CleanupJob) Execute(ctx context.Context) error {
    j.logger.Info().Msg("Starting cleanup job")
    
    // Delete expired tokens
    if err := j.cleanupTokens(ctx); err != nil {
        return err
    }
    
    // Delete old logs
    if err := j.cleanupLogs(ctx); err != nil {
        return err
    }
    
    // Delete expired tags
    if err := j.cleanupTags(ctx); err != nil {
        return err
    }
    
    j.logger.Info().Msg("Cleanup job completed")
    return nil
}

Features:
- Cron expression support
- Job timeout handling
- Automatic retry on failure
- Job locking (prevent concurrent runs)
- Job monitoring and logging
- Job status tracking
- Manual job trigger
- Job pause/resume
- Job history

Monitoring:
- Job execution time
- Success/failure rate
- Error tracking
- Alert on failures
- Dashboard for job status

API Endpoints:
- GET    /api/jobs              - List scheduled jobs
- POST   /api/jobs              - Create job
- GET    /api/jobs/:id          - Get job details
- PUT    /api/jobs/:id          - Update job
- DELETE /api/jobs/:id          - Delete job
- POST   /api/jobs/:id/run      - Trigger job manually
- POST   /api/jobs/:id/pause    - Pause job
- POST   /api/jobs/:id/resume   - Resume job
- GET    /api/jobs/:id/history  - Get job execution history

Follow Go best practices and project .cursorrules.
```

---

### Prompt 4.3: Docker Support & Google Cloud Deployment

```
Implement complete Docker containerization with Google Cloud support:

Requirements:
1. Create multi-stage Dockerfile
2. Create docker-compose for development
3. Create docker-compose for production
4. Add health checks
5. Configure volumes
6. Setup networking
7. Add environment-specific configs
8. Support for Google Cloud Run deployment
9. Support for Google Kubernetes Engine (GKE)
10. Cloud SQL Proxy integration

Files to create:
Dockerfile
docker-compose.yml
docker-compose.dev.yml
docker-compose.prod.yml
.dockerignore
cloudbuild.yaml              # Google Cloud Build
k8s/
├── deployment.yaml          # Kubernetes deployment
├── service.yaml             # Kubernetes service
├── ingress.yaml             # Kubernetes ingress
└── configmap.yaml           # Kubernetes configmap
scripts/
├── docker-entrypoint.sh
├── wait-for-postgres.sh
├── cloud-sql-proxy.sh       # Cloud SQL proxy script
└── init-db.sh

Dockerfile (multi-stage):
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./app/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/app/config ./config
COPY --from=builder /app/app/i18n ./i18n
COPY --from=builder /app/app/email/templates ./email/templates

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./main"]

docker-compose.yml (development):
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go-api
    ports:
      - "8080:8080"
    environment:
      - ENV=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=app_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - .:/app
      - ./uploads:/app/uploads
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    container_name: postgres
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=app_db
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./scripts/init-db.sh:/docker-entrypoint-initdb.d/init-db.sh
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    networks:
      - app-network
    restart: unless-stopped

  # For local development - use MinIO
  # For production - use Google Cloud Storage
  minio:
    image: minio/minio:latest
    container_name: minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3
    networks:
      - app-network
    restart: unless-stopped
    profiles:
      - local  # Only start in local development

  # Cloud SQL Proxy for connecting to Cloud SQL
  cloud-sql-proxy:
    image: gcr.io/cloudsql-docker/gce-proxy:latest
    container_name: cloud-sql-proxy
    command:
      - "/cloud_sql_proxy"
      - "-instances=PROJECT_ID:REGION:INSTANCE_NAME=tcp:0.0.0.0:5432"
    ports:
      - "5432:5432"
    networks:
      - app-network
    restart: unless-stopped
    profiles:
      - cloud  # Only start when using Cloud SQL

  nginx:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - app-network
    restart: unless-stopped

volumes:
  postgres-data:
  redis-data:
  minio-data:

networks:
  app-network:
    driver: bridge

docker-compose.prod.yml:
version: '3.8'

services:
  app:
    image: your-registry/go-api:latest
    container_name: go-api-prod
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - JWT_SECRET=${JWT_SECRET}
    volumes:
      - /var/log/app:/app/logs
      - /var/uploads:/app/uploads
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    networks:
      - app-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

.dockerignore:
.git
.gitignore
.env*
*.md
Dockerfile*
docker-compose*
.vscode
.idea
*.log
tmp/
vendor/
test/
coverage*
*.test

scripts/docker-entrypoint.sh:
#!/bin/sh
set -e

echo "Waiting for PostgreSQL..."
./wait-for-postgres.sh postgres

echo "Running migrations..."
./main migrate

echo "Starting application..."
exec ./main

scripts/wait-for-postgres.sh:
#!/bin/sh
set -e

host="$1"
shift

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"

cloudbuild.yaml (Google Cloud Build):
steps:
  # Build the container image
  - name: 'gcr.io/cloud-builders/docker'
    args:
      - 'build'
      - '-t'
      - 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
      - '-t'
      - 'gcr.io/$PROJECT_ID/go-api:latest'
      - '.'
  
  # Push the container image to Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    args:
      - 'push'
      - 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
  
  # Deploy to Cloud Run
  - name: 'gcr.io/cloud-builders/gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'go-api'
      - '--image'
      - 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
      - '--region'
      - 'asia-northeast1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'
      - '--add-cloudsql-instances'
      - '$PROJECT_ID:asia-northeast1:postgres-instance'
      - '--set-env-vars'
      - 'ENV=production,DB_HOST=/cloudsql/$PROJECT_ID:asia-northeast1:postgres-instance'

images:
  - 'gcr.io/$PROJECT_ID/go-api:$COMMIT_SHA'
  - 'gcr.io/$PROJECT_ID/go-api:latest'

k8s/deployment.yaml (GKE Deployment):
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-api
  labels:
    app: go-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-api
  template:
    metadata:
      labels:
        app: go-api
    spec:
      serviceAccountName: go-api-sa
      containers:
      - name: go-api
        image: gcr.io/PROJECT_ID/go-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        - name: DB_HOST
          value: "127.0.0.1"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: cloudsql-db-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: cloudsql-db-credentials
              key: password
        - name: DB_NAME
          value: "app_db"
        - name: GCS_BUCKET
          value: "PROJECT_ID-uploads"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      
      # Cloud SQL Proxy sidecar
      - name: cloud-sql-proxy
        image: gcr.io/cloudsql-docker/gce-proxy:latest
        command:
          - "/cloud_sql_proxy"
          - "-instances=PROJECT_ID:REGION:INSTANCE_NAME=tcp:5432"
        securityContext:
          runAsNonRoot: true

k8s/service.yaml:
apiVersion: v1
kind: Service
metadata:
  name: go-api-service
spec:
  type: LoadBalancer
  selector:
    app: go-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080

Makefile additions:
docker-build:
	docker build -t go-api:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-shell:
	docker-compose exec app sh

docker-prod-deploy:
	docker-compose -f docker-compose.prod.yml up -d

# Google Cloud commands
gcloud-build:
	gcloud builds submit --config cloudbuild.yaml

gcloud-deploy-run:
	gcloud run deploy go-api \
		--image gcr.io/$(PROJECT_ID)/go-api:latest \
		--region asia-northeast1 \
		--platform managed \
		--allow-unauthenticated

gcloud-deploy-gke:
	kubectl apply -f k8s/

gcloud-logs:
	gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=go-api" --limit 50

Follow Go best practices and project .cursorrules.
```

---

## 📋 PHASE 5: NOTIFICATIONS & ANALYTICS

### Prompt 5.1: Push Notifications - Firebase

```
Implement Firebase Cloud Messaging for push notifications:

Requirements:
1. Setup Firebase Admin SDK
2. Implement device token management
3. Create notification service
4. Add notification templates
5. Support for topic-based notifications
6. Implement notification scheduling
7. Add notification history

File structure:
app/
├── notifications/
│   ├── firebase.go            # Firebase setup
│   ├── push_notification.go   # Push notification sender
│   ├── device_token.go        # Device token management
│   ├── topics.go              # Topic management
│   └── templates.go           # Notification templates
├── models/
│   ├── notification.go        # Notification model
│   ├── device_token.go        # Device token model
│   └── notification_log.go    # Notification log
├── repository/
│   ├── notification_repository.go
│   └── device_token_repository.go
├── services/
│   └── notification_service.go
└── handlers/
    └── notification_handler.go

Models:

1. DeviceToken:
type DeviceToken struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
    User      User      `gorm:"foreignKey:UserID"`
    Token     string    `gorm:"unique;not null"`
    Platform  string    `gorm:"not null"` // ios, android, web
    DeviceID  string
    AppVersion string
    IsActive  bool      `gorm:"default:true"`
    LastUsedAt time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

2. Notification:
type Notification struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Title       string    `gorm:"not null"`
    Body        string    `gorm:"not null"`
    Type        string    `gorm:"not null"` // tag_assigned, tag_expiring, shop_update
    Priority    string    `gorm:"default:'normal'"` // low, normal, high
    ImageURL    string
    ActionURL   string
    Data        datatypes.JSON
    TargetType  string    `gorm:"not null"` // user, topic, all
    TargetID    string
    ScheduledAt *time.Time
    SentAt      *time.Time
    Status      string    `gorm:"default:'pending'"` // pending, sent, failed
    SuccessCount int      `gorm:"default:0"`
    FailureCount int      `gorm:"default:0"`
    CreatedBy   uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

3. NotificationLog:
type NotificationLog struct {
    ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    NotificationID uuid.UUID `gorm:"type:uuid;not null;index"`
    UserID         uuid.UUID `gorm:"type:uuid;index"`
    DeviceTokenID  uuid.UUID `gorm:"type:uuid;index"`
    Status         string    `gorm:"not null"` // sent, failed, clicked, dismissed
    ErrorMessage   string
    SentAt         time.Time
    ClickedAt      *time.Time
    CreatedAt      time.Time
}

Firebase setup:
type FirebaseService struct {
    app    *firebase.App
    client *messaging.Client
    logger *zerolog.Logger
}

func NewFirebaseService(credentialsPath string) (*FirebaseService, error) {
    opt := option.WithCredentialsFile(credentialsPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return nil, err
    }
    
    client, err := app.Messaging(context.Background())
    if err != nil {
        return nil, err
    }
    
    return &FirebaseService{
        app:    app,
        client: client,
    }, nil
}

Notification types:

1. Tag Assigned:
   - Title: "新しいタグが割り当てられました"
   - Body: "{{shop_name}}からタグが割り当てられました"
   - Action: Open tag details

2. Tag Expiring:
   - Title: "タグの有効期限が近づいています"
   - Body: "{{tag_name}}は{{days}}日後に期限切れになります"
   - Action: Renew tag

3. Shop Update:
   - Title: "{{shop_name}}からのお知らせ"
   - Body: "{{message}}"
   - Action: Open shop page

4. Promotion:
   - Title: "特別オファー"
   - Body: "{{offer_details}}"
   - Action: Open promotion

Features:

1. Send to single user:
func (s *NotificationService) SendToUser(ctx context.Context, userID uuid.UUID, notification Notification) error {
    // Get user device tokens
    tokens, err := s.deviceTokenRepo.GetActiveTokensByUserID(userID)
    if err != nil {
        return err
    }
    
    // Send to all devices
    for _, token := range tokens {
        err := s.sendToToken(ctx, token.Token, notification)
        if err != nil {
            s.logFailure(notification.ID, userID, token.ID, err)
        } else {
            s.logSuccess(notification.ID, userID, token.ID)
        }
    }
    
    return nil
}

2. Send to topic:
func (s *NotificationService) SendToTopic(ctx context.Context, topic string, notification Notification) error {
    message := &messaging.Message{
        Topic: topic,
        Notification: &messaging.Notification{
            Title: notification.Title,
            Body:  notification.Body,
            ImageURL: notification.ImageURL,
        },
        Data: notification.Data,
    }
    
    _, err := s.firebase.client.Send(ctx, message)
    return err
}

3. Send to multiple users:
func (s *NotificationService) SendToMultipleUsers(ctx context.Context, userIDs []uuid.UUID, notification Notification) error {
    // Get all device tokens
    tokens, err := s.deviceTokenRepo.GetActiveTokensByUserIDs(userIDs)
    if err != nil {
        return err
    }
    
    // Batch send (max 500 per batch)
    return s.batchSend(ctx, tokens, notification)
}

4. Schedule notification:
func (s *NotificationService) ScheduleNotification(ctx context.Context, notification Notification, scheduledAt time.Time) error {
    notification.ScheduledAt = &scheduledAt
    notification.Status = "scheduled"
    return s.notificationRepo.Create(&notification)
}

Topics:
- all_users
- shop_{shop_id}
- tag_type_{type}
- user_role_{role}

API Endpoints:
- POST   /api/notifications/send          - Send notification
- POST   /api/notifications/schedule      - Schedule notification
- GET    /api/notifications               - List notifications
- GET    /api/notifications/:id           - Get notification details
- DELETE /api/notifications/:id           - Delete notification
- POST   /api/device-tokens               - Register device token
- DELETE /api/device-tokens/:id           - Remove device token
- POST   /api/notifications/:id/click     - Track notification click
- GET    /api/notifications/history       - Get user notification history

Configuration:
type NotificationConfig struct {
    FirebaseCredentialsPath string
    DefaultPriority         string
    BatchSize               int
    RetryAttempts           int
    RetryDelay              time.Duration
}

Follow Go best practices and project .cursorrules.
```

---

### Prompt 5.2: Analytics & Tracking

```
Implement comprehensive analytics and tracking system:

Requirements:
1. Create event tracking system
2. Implement data aggregation
3. Generate analytics reports
4. Track tag usage
5. Track user behavior
6. Generate shop analytics
7. Create dashboard metrics

File structure:
app/
├── analytics/
│   ├── tracker.go             # Event tracker
│   ├── aggregator.go          # Data aggregator
│   ├── reporter.go            # Report generator
│   ├── metrics.go             # Metrics calculator
│   └── dashboard.go           # Dashboard data
├── models/
│   ├── analytics_event.go     # Event model
│   ├── analytics_report.go    # Report model
│   └── analytics_metric.go    # Metric model
├── repository/
│   ├── analytics_repository.go
│   └── metrics_repository.go
├── services/
│   └── analytics_service.go
└── handlers/
    └── analytics_handler.go

Models:

1. AnalyticsEvent:
type AnalyticsEvent struct {
    ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    EventType  string    `gorm:"not null;index"` // tag_scan, user_login, shop_view
    EventName  string    `gorm:"not null"`
    UserID     uuid.UUID `gorm:"type:uuid;index"`
    ShopID     uuid.UUID `gorm:"type:uuid;index"`
    TagID      uuid.UUID `gorm:"type:uuid;index"`
    DeviceID   uuid.UUID `gorm:"type:uuid;index"`
    SessionID  string    `gorm:"index"`
    IPAddress  string
    UserAgent  string
    Platform   string
    Location   string
    Properties datatypes.JSON
    Timestamp  time.Time `gorm:"index"`
    CreatedAt  time.Time
}

2. AnalyticsReport:
type AnalyticsReport struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ReportType  string    `gorm:"not null"` // daily, weekly, monthly
    ReportName  string    `gorm:"not null"`
    Period      string    `gorm:"not null"` // 2024-01, 2024-W01
    ShopID      uuid.UUID `gorm:"type:uuid;index"`
    Data        datatypes.JSON
    GeneratedAt time.Time
    CreatedAt   time.Time
}

3. AnalyticsMetric:
type AnalyticsMetric struct {
    ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    MetricName string    `gorm:"not null;index"`
    MetricType string    `gorm:"not null"` // counter, gauge, histogram
    Value      float64   `gorm:"not null"`
    Tags       datatypes.JSON
    Timestamp  time.Time `gorm:"index"`
    CreatedAt  time.Time
}

Event types to track:

1. User events:
   - user_registered
   - user_login
   - user_logout
   - profile_updated
   - password_changed

2. Tag events:
   - tag_created
   - tag_assigned
   - tag_scanned
   - tag_revoked
   - tag_expired

3. Shop events:
   - shop_created
   - shop_viewed
   - shop_updated
   - shop_deleted

4. Device events:
   - device_registered
   - device_online
   - device_offline
   - device_scan

5. System events:
   - api_request
   - api_error
   - job_executed
   - email_sent

Metrics to calculate:

1. Tag metrics:
   - Total tags created
   - Active tags
   - Tags scanned today/week/month
   - Average scans per tag
   - Most popular tags
   - Tag usage by shop
   - Tag usage by time of day

2. User metrics:
   - Total users
   - Active users (DAU, WAU, MAU)
   - New users
   - User retention rate
   - User engagement score
   - Average session duration
   - User churn rate

3. Shop metrics:
   - Total shops
   - Active shops
   - Shop performance score
   - Tags per shop
   - Customers per shop
   - Revenue per shop (if applicable)

4. System metrics:
   - API requests per second
   - Average response time
   - Error rate
   - Database query time
   - Cache hit rate

Analytics service methods:

1. Track event:
func (s *AnalyticsService) TrackEvent(ctx context.Context, event AnalyticsEvent) error {
    // Validate event
    if err := s.validateEvent(event); err != nil {
        return err
    }
    
    // Save event
    if err := s.repo.CreateEvent(&event); err != nil {
        return err
    }
    
    // Update real-time metrics
    s.updateMetrics(event)
    
    return nil
}

2. Generate report:
func (s *AnalyticsService) GenerateReport(ctx context.Context, reportType string, period string) (*AnalyticsReport, error) {
    // Get events for period
    events, err := s.repo.GetEventsByPeriod(reportType, period)
    if err != nil {
        return nil, err
    }
    
    // Aggregate data
    data := s.aggregateEvents(events)
    
    // Create report
    report := &AnalyticsReport{
        ReportType:  reportType,
        ReportName:  fmt.Sprintf("%s Report - %s", reportType, period),
        Period:      period,
        Data:        data,
        GeneratedAt: time.Now(),
    }
    
    return report, s.repo.CreateReport(report)
}

3. Get dashboard metrics:
func (s *AnalyticsService) GetDashboardMetrics(ctx context.Context, shopID *uuid.UUID) (DashboardMetrics, error) {
    metrics := DashboardMetrics{}
    
    // Get tag metrics
    metrics.TotalTags = s.getTotalTags(shopID)
    metrics.ActiveTags = s.getActiveTags(shopID)
    metrics.TodayScans = s.getTodayScans(shopID)
    
    // Get user metrics
    metrics.TotalUsers = s.getTotalUsers(shopID)
    metrics.ActiveUsersToday = s.getActiveUsersToday(shopID)
    metrics.NewUsersThisWeek = s.getNewUsersThisWeek(shopID)
    
    // Get trends
    metrics.ScanTrend = s.getScanTrend(shopID, 7) // Last 7 days
    metrics.UserTrend = s.getUserTrend(shopID, 7)
    
    return metrics, nil
}

4. Get tag usage analytics:
func (s *AnalyticsService) GetTagUsageAnalytics(ctx context.Context, tagID uuid.UUID, period string) (*TagAnalytics, error) {
    analytics := &TagAnalytics{
        TagID: tagID,
    }
    
    // Get scan count
    analytics.TotalScans = s.getTagScanCount(tagID, period)
    
    // Get unique users
    analytics.UniqueUsers = s.getTagUniqueUsers(tagID, period)
    
    // Get scan times
    analytics.ScansByHour = s.getTagScansByHour(tagID, period)
    analytics.ScansByDay = s.getTagScansByDay(tagID, period)
    
    // Get locations
    analytics.TopLocations = s.getTagTopLocations(tagID, period)
    
    return analytics, nil
}

API Endpoints:
- POST   /api/analytics/track            - Track event
- GET    /api/analytics/dashboard        - Get dashboard metrics
- GET    /api/analytics/reports          - List reports
- POST   /api/analytics/reports/generate - Generate report
- GET    /api/analytics/reports/:id      - Get report
- GET    /api/analytics/tags/:id         - Get tag analytics
- GET    /api/analytics/shops/:id        - Get shop analytics
- GET    /api/analytics/users/:id        - Get user analytics
- GET    /api/analytics/trends           - Get trends
- GET    /api/analytics/export           - Export analytics data

Dashboard response:
{
  "metrics": {
    "tags": {
      "total": 1250,
      "active": 980,
      "scans_today": 456,
      "trend": [120, 145, 178, 156, 189, 234, 456]
    },
    "users": {
      "total": 5430,
      "active_today": 234,
      "new_this_week": 45,
      "trend": [180, 195, 210, 198, 220, 245, 234]
    },
    "shops": {
      "total": 45,
      "active": 38
    }
  },
  "top_tags": [...],
  "top_shops": [...],
  "recent_activity": [...]
}

Aggregation job:
- Run hourly to aggregate events
- Calculate daily/weekly/monthly metrics
- Generate reports
- Update dashboard cache

Follow Go best practices and project .cursorrules.
```

---

## 🎯 Usage Instructions

### How to use these prompts:

1. **Sequential Implementation:**
   - Start with Phase 1 prompts
   - Complete each phase before moving to the next
   - Test thoroughly after each implementation

2. **Parallel Implementation:**
   - Some features can be developed in parallel
   - Coordinate with team members
   - Merge carefully to avoid conflicts

3. **Customization:**
   - Adjust prompts based on specific requirements
   - Add or remove features as needed
   - Maintain consistency with existing code

4. **Testing:**
   - Write tests alongside implementation
   - Follow TDD approach where possible
   - Maintain test coverage above 80%

5. **Documentation:**
   - Update API documentation
   - Add code comments
   - Update README with new features

6. **Code Review:**
   - Review code before merging
   - Follow Go best practices
   - Adhere to project .cursorrules

---

## 📝 Notes

- All prompts follow Go best practices
- Implement error handling properly
- Use context for cancellation
- Add proper logging
- Write comprehensive tests
- Document public APIs
- Follow SOLID principles
- Keep code DRY
- Use dependency injection
- Implement graceful shutdown

---

## 🔄 Continuous Improvement

After implementing each phase:
1. Gather feedback
2. Identify improvements
3. Refactor if needed
4. Update documentation
5. Plan next phase

---

**Happy Coding! 🚀**
