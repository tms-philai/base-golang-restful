# 🌳 Cấu Trúc Dự Án - Standard Go Layout

## Cấu trúc HOÀN CHỈNH theo golang-standards/project-layout

```
base-golang-restful/
│
├── cmd/                                    # Main applications
│   └── server/                             # Server application
│       └── main.go                         # Entry point (từ app/main.go)
│
├── internal/                               # Private application code
│   ├── app/                                # Application layer
│   │   ├── handlers/                       # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── email_handler.go
│   │   │   ├── notification_handler.go
│   │   │   ├── api_handler.go
│   │   │   └── monitoring_handler.go
│   │   │
│   │   ├── routes/                         # Route definitions
│   │   │   └── routes.go
│   │   │
│   │   └── middleware/                     # HTTP middlewares
│   │       ├── auth.go
│   │       ├── rbac.go
│   │       ├── versioning.go
│   │       ├── error_handler.go
│   │       ├── cors.go
│   │       └── setup.go
│   │
│   ├── domain/                             # Business domain
│   │   ├── models/                         # Domain models
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   ├── product.go
│   │   │   ├── auth.go
│   │   │   ├── email.go
│   │   │   ├── notification.go
│   │   │   ├── file.go
│   │   │   ├── refresh_token.go
│   │   │   ├── common.go
│   │   │   ├── response.go
│   │   │   ├── error_response.go
│   │   │   ├── pagination.go
│   │   │   └── api.go
│   │   │
│   │   ├── services/                       # Business logic
│   │   │   ├── base_service.go
│   │   │   ├── user_service.go
│   │   │   ├── product_service.go
│   │   │   ├── jwt_service.go
│   │   │   ├── token_service.go
│   │   │   ├── file_service.go
│   │   │   └── file_validator.go
│   │   │
│   │   └── repository/                     # Data access layer
│   │       ├── base_repository.go
│   │       ├── user_repository.go
│   │       ├── product_repository.go
│   │       ├── role_repository.go
│   │       ├── permission_repository.go
│   │       ├── refresh_token_repository.go
│   │       ├── file_repository.go
│   │       └── interfaces.go
│   │
│   └── pkg/                                # Internal shared packages
│       ├── auth/                           # Authentication
│       │   ├── jwt.go
│       │   ├── jwt_test.go.old
│       │   ├── jwt_test_note.go
│       │   ├── password.go
│       │   └── password_test.go
│       │
│       ├── cache/                          # Caching
│       │   ├── cache_interface.go
│       │   ├── cache_service.go
│       │   ├── query_cache.go
│       │   └── redis_client.go
│       │
│       ├── config/                         # Configuration
│       │   ├── config.go
│       │   ├── config_test.go
│       │   ├── server.go
│       │   ├── database.go
│       │   ├── jwt.go
│       │   ├── jwt_test.go
│       │   ├── redis.go
│       │   ├── redis_test.go
│       │   ├── storage.go
│       │   ├── email.go
│       │   └── logging.go
│       │
│       ├── database/                       # Database
│       │   ├── connection.go
│       │   ├── connection_test.go
│       │   ├── migration.go
│       │   ├── migration_test.go
│       │   ├── seeder.go
│       │   ├── seeder_test.go
│       │   ├── transaction.go
│       │   └── transaction_test.go
│       │
│       ├── email/                          # Email service
│       │   ├── email_client.go
│       │   ├── email_client_test.go
│       │   ├── email_service.go
│       │   ├── notification.go
│       │   ├── notification_test.go
│       │   ├── template_manager.go
│       │   └── template_manager_test.go
│       │
│       ├── storage/                        # File storage
│       │   ├── storage_interface.go
│       │   ├── local_storage.go
│       │   ├── local_storage_test.go
│       │   ├── s3_storage.go
│       │   ├── s3_storage_test.go
│       │   ├── gcs_storage.go
│       │   └── gcs_storage_test.go
│       │
│       ├── health/                         # Health checks
│       │   ├── health_check.go
│       │   └── health_check_test.go
│       │
│       ├── metrics/                        # Metrics
│       │   └── metrics.go
│       │
│       ├── logger/                         # Logging
│       │   ├── logger.go
│       │   ├── logger_test.go
│       │   ├── fields.go
│       │   ├── fields_test.go
│       │   ├── middleware.go
│       │   └── middleware_test.go
│       │
│       ├── i18n/                           # Internationalization
│       │   ├── i18n.go
│       │   ├── middleware.go
│       │   └── locales/
│       │       ├── en.json
│       │       ├── vi.json
│       │       └── ja.json
│       │
│       ├── errors/                         # Error handling
│       │   ├── app_error.go
│       │   ├── app_error_test.go
│       │   ├── localized_error.go
│       │   └── localized_error_test.go
│       │
│       └── utils/                          # Utilities
│           ├── string.go
│           ├── string_test.go
│           ├── password.go
│           ├── pagination.go
│           ├── pagination_test.go
│           ├── validator.go
│           └── validator_test.go
│
├── pkg/                                    # Public libraries (nếu cần export)
│   └── (empty for now)
│
├── api/                                    # API definitions
│   └── openapi/                            # OpenAPI/Swagger specs
│       ├── swagger.json
│       ├── swagger.yaml
│       └── docs.go
│
├── configs/                                # Configuration files
│   ├── .env.example
│   ├── .env.development
│   ├── .env.production
│   ├── .env.mailhog
│   └── app.yaml
│
├── deployments/                            # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile
│   │   ├── Dockerfile.dev
│   │   ├── docker-compose.yml
│   │   └── docker-compose.email-test.yml
│   │
│   └── kubernetes/                         # K8s configs (future)
│       ├── deployment.yaml
│       ├── service.yaml
│       ├── configmap.yaml
│       └── secret.yaml
│
├── scripts/                                # Build and maintenance scripts
│   ├── build.sh                            # Build binary
│   ├── test.sh                             # Run tests
│   ├── swagger.sh                          # Generate Swagger
│   ├── migrate.sh                          # Run migrations
│   ├── seed.sh                             # Seed database
│   └── docker-build.sh                     # Docker build helper
│
├── docs/                                   # Documentation
│   ├── architecture.md                     # System architecture
│   ├── api-documentation.md                # API guide
│   ├── deployment.md                       # Deployment guide
│   ├── development.md                      # Development guide
│   ├── testing.md                          # Testing strategy
│   ├── phase-guides/                       # Phase implementation guides
│   │   ├── phase6-email-notifications.md
│   │   ├── phase7-api-testing.md
│   │   └── phase8-monitoring-deployment.md
│   └── diagrams/                           # Architecture diagrams
│       ├── system-overview.mmd
│       ├── auth-flow.mmd
│       └── database-schema.mmd
│
├── test/                                   # Integration tests & test data
│   ├── integration/
│   │   └── api_integration_test.go
│   ├── examples/
│   │   └── api_versioning_test.go
│   ├── fixtures/                           # Test data
│   │   ├── users.json
│   │   └── products.json
│   └── mocks/                              # Test mocks
│       ├── jwt_service_mock.go
│       ├── user_service_mock.go
│       └── product_service_mock.go
│
├── web/                                    # Web assets (if needed)
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/                          # Email templates
│       └── emails/
│           ├── welcome.html
│           ├── welcome.txt
│           ├── reset-password.html
│           └── reset-password.txt
│
├── assets/                                 # Project assets
│   ├── logo.png
│   └── banner.png
│
├── build/                                  # Build artifacts (gitignored)
│   ├── bin/
│   │   └── server
│   └── package/
│
├── .github/                                # GitHub specific
│   └── workflows/
│       ├── ci.yml
│       ├── test.yml
│       └── deploy.yml
│
├── .vscode/                                # VSCode settings
│   ├── settings.json
│   ├── launch.json
│   └── tasks.json
│
├── go.mod                                  # Root level
├── go.sum
├── Makefile                                # Build commands
├── .gitignore
├── .editorconfig
├── .env                                    # Local dev (gitignored)
├── README.md
├── LICENSE
├── CHANGELOG.md
└── CONTRIBUTING.md
```

---

## 📂 CHI TIẾT TỪNG THƯ MỤC

### **1. `/cmd` - Main Applications**

```
cmd/
└── server/
    └── main.go                 # 223 lines
        # - Load configuration
        # - Initialize dependencies
        # - Setup router
        # - Start server
```

**Mục đích:** Entry points cho các executables  
**Quy tắc:** Keep it simple, minimal logic  

---

### **2. `/internal` - Private Code**

#### **2.1. `/internal/app` - Application Layer**

```
internal/app/
├── handlers/                   # HTTP handlers (7 files)
│   ├── auth_handler.go         # 312 lines - Auth endpoints
│   ├── user_handler.go         # ~250 lines - User CRUD
│   ├── product_handler.go      # ~380 lines - Product CRUD
│   ├── email_handler.go        # 213 lines - Email endpoints
│   ├── notification_handler.go # ~220 lines - Notification endpoints
│   ├── api_handler.go          # 181 lines - API info/version
│   └── monitoring_handler.go   # 67 lines - Health/metrics
│
├── routes/                     # Routing (1 file)
│   └── routes.go               # 163 lines - All route definitions
│
└── middleware/                 # HTTP middlewares (6 files)
    ├── auth.go                 # 158 lines - JWT authentication
    ├── rbac.go                 # 221 lines - Role-based access
    ├── versioning.go           # 100 lines - API versioning
    ├── error_handler.go        # 97 lines - Error handling
    ├── cors.go                 # 96 lines - CORS
    └── setup.go                # 19 lines - Middleware setup
```

#### **2.2. `/internal/domain` - Business Logic**

```
internal/domain/
├── models/                     # Domain models (14 files)
│   ├── user.go                 # 145 lines - User model
│   ├── role.go                 # 98 lines - Role model
│   ├── permission.go           # 54 lines - Permission model
│   ├── product.go              # 67 lines - Product model
│   ├── auth.go                 # 115 lines - Auth DTOs
│   ├── email.go                # 39 lines - Email DTOs
│   ├── notification.go         # 64 lines - Notification DTOs
│   ├── file.go                 # 85 lines - File model
│   ├── refresh_token.go        # 49 lines - Refresh token model
│   ├── common.go               # 48 lines - Common types
│   ├── response.go             # 111 lines - Response wrappers
│   ├── error_response.go       # 67 lines - Error responses
│   ├── pagination.go           # 124 lines - Pagination
│   └── api.go                  # 63 lines - API info models
│
├── services/                   # Business logic (7 files)
│   ├── base_service.go         # 78 lines - Base service
│   ├── user_service.go         # 253 lines - User business logic
│   ├── product_service.go      # ~280 lines - Product logic
│   ├── jwt_service.go          # ~150 lines - JWT operations
│   ├── token_service.go        # 142 lines - Token management
│   ├── file_service.go         # ~200 lines - File operations
│   └── file_validator.go       # 89 lines - File validation
│
└── repository/                 # Data access (9 files)
    ├── base_repository.go      # 56 lines - Base repository
    ├── user_repository.go      # 147 lines - User data access
    ├── product_repository.go   # ~180 lines - Product data access
    ├── role_repository.go      # 89 lines - Role data access
    ├── permission_repository.go # 76 lines - Permission data access
    ├── refresh_token_repository.go # 123 lines - Token data access
    ├── file_repository.go      # 98 lines - File data access
    ├── query_builder.go        # ~100 lines - Query builder
    └── interfaces.go           # 45 lines - Repository interfaces
```

#### **2.3. `/internal/pkg` - Shared Internal Packages**

```
internal/pkg/
├── auth/                       # Authentication (4 files)
│   ├── jwt.go                  # 155 lines - JWT manager
│   ├── jwt_test.go.old         # Old tests
│   ├── password.go             # 84 lines - Password hashing
│   └── password_test.go        # Tests
│
├── cache/                      # Caching (4 files)
│   ├── cache_interface.go      # Interface definitions
│   ├── cache_service.go        # Cache service
│   ├── query_cache.go          # Query result caching
│   └── redis_client.go         # Redis client
│
├── config/                     # Configuration (10 files)
│   ├── config.go               # 125 lines - Main config
│   ├── config_test.go
│   ├── server.go               # Server config
│   ├── database.go             # Database config
│   ├── jwt.go                  # JWT config
│   ├── jwt_test.go
│   ├── redis.go                # Redis config
│   ├── redis_test.go
│   ├── storage.go              # Storage config
│   ├── email.go                # Email config
│   └── logging.go              # Logging config
│
├── database/                   # Database (8 files)
│   ├── connection.go           # 126 lines - DB connection
│   ├── connection_test.go
│   ├── migration.go            # 134 lines - Migrations
│   ├── migration_test.go
│   ├── seeder.go               # 57 lines - Data seeding
│   ├── seeder_test.go
│   ├── transaction.go          # 153 lines - Transaction support
│   └── transaction_test.go
│
├── email/                      # Email service (8 files)
│   ├── email_client.go         # 138 lines - SMTP client
│   ├── email_client_test.go
│   ├── email_service.go        # 147 lines - Email service
│   ├── notification.go         # 186 lines - Notification system
│   ├── notification_test.go
│   ├── template_manager.go     # 190 lines - Template engine
│   └── template_manager_test.go
│
├── storage/                    # File storage (6 files)
│   ├── storage_interface.go    # Storage interface
│   ├── local_storage.go        # 222 lines - Local storage
│   ├── local_storage_test.go
│   ├── s3_storage.go           # 227 lines - AWS S3
│   ├── s3_storage_test.go
│   ├── gcs_storage.go          # 201 lines - Google Cloud Storage
│   └── gcs_storage_test.go
│
├── health/                     # Health checks (2 files)
│   ├── health_check.go         # 146 lines - Health checker
│   └── health_check_test.go
│
├── metrics/                    # Metrics (1 file)
│   └── metrics.go              # 124 lines - Metrics collection
│
├── logger/                     # Structured logging (6 files)
│   ├── logger.go               # Logger implementation
│   ├── logger_test.go
│   ├── fields.go               # Log fields
│   ├── fields_test.go
│   ├── middleware.go           # Logging middleware
│   └── middleware_test.go
│
├── i18n/                       # Internationalization (4 files)
│   ├── i18n.go                 # 188 lines - i18n manager
│   ├── i18n_test.go
│   ├── middleware.go           # 99 lines - i18n middleware
│   └── middleware_test.go
│
├── errors/                     # Error handling (4 files)
│   ├── app_error.go            # Application errors
│   ├── app_error_test.go
│   ├── localized_error.go      # Localized errors
│   └── localized_error_test.go
│
└── utils/                      # Utilities (7 files)
    ├── string.go               # String utilities
    ├── string_test.go
    ├── password.go             # 67 lines - Password utils
    ├── pagination.go           # Pagination helpers
    ├── pagination_test.go
    ├── validator.go            # Validation helpers
    └── validator_test.go
```

---

### **3. `/api` - API Definitions**

```
api/
└── openapi/                    # OpenAPI/Swagger
    ├── swagger.json            # Generated
    ├── swagger.yaml            # Generated
    └── docs.go                 # Generated
```

---

### **4. `/configs` - Configuration**

```
configs/
├── .env.example                # Template env file
├── .env.development            # Dev environment
├── .env.production             # Production environment
├── .env.docker                 # Docker environment
├── .env.mailhog                # MailHog testing
└── app.yaml                    # Application config (optional)
```

---

### **5. `/deployments` - Deployment**

```
deployments/
├── docker/
│   ├── Dockerfile              # Production Dockerfile
│   ├── Dockerfile.dev          # Development Dockerfile
│   ├── docker-compose.yml      # Full stack
│   ├── docker-compose.email-test.yml  # Email testing
│   └── .dockerignore
│
└── kubernetes/                 # K8s manifests
    ├── namespace.yaml
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    ├── configmap.yaml
    ├── secret.yaml
    └── hpa.yaml                # Horizontal Pod Autoscaler
```

---

### **6. `/scripts` - Scripts**

```
scripts/
├── build.sh                    # Build binary
├── test.sh                     # Run tests
├── test-coverage.sh            # Test with coverage
├── swagger.sh                  # Generate Swagger docs
├── migrate.sh                  # Database migrations
├── seed.sh                     # Seed database
├── docker-build.sh             # Build Docker image
├── docker-push.sh              # Push to registry
├── lint.sh                     # Run linters
└── deploy.sh                   # Deploy to production
```

---

### **7. `/docs` - Documentation**

```
docs/
├── README.md                   # Documentation overview
├── ARCHITECTURE.md             # System architecture
├── API_DOCUMENTATION.md        # API usage guide
├── DEVELOPMENT.md              # Development guide
├── DEPLOYMENT.md               # Deployment guide
├── TESTING.md                  # Testing strategy
├── GOOGLE_CLOUD_DEPLOYMENT.md  # GCP deployment
├── REFACTORING_GUIDE.md        # This guide
├── PROJECT_COMPLETION_SUMMARY.md
│
├── phase-guides/               # Implementation guides
│   ├── PHASE1_CORE.md
│   ├── PHASE2_I18N_RBAC.md
│   ├── PHASE3_AUTH.md
│   ├── PHASE4_FILE.md
│   ├── PHASE5_CACHE.md
│   ├── PHASE6_EMAIL.md
│   ├── PHASE7_API.md
│   └── PHASE8_MONITORING.md
│
└── diagrams/                   # Mermaid diagrams
    ├── system-overview.mmd
    ├── authentication-flow.mmd
    ├── database-erd.mmd
    ├── email-notification-flow.mmd
    └── deployment-architecture.mmd
```

---

### **8. `/test` - Additional Tests**

```
test/
├── integration/                # Integration tests
│   ├── api_integration_test.go
│   ├── auth_flow_test.go
│   └── database_test.go
│
├── e2e/                        # End-to-end tests
│   └── user_journey_test.go
│
├── fixtures/                   # Test data
│   ├── users.json
│   ├── products.json
│   └── test_config.yaml
│
├── mocks/                      # Generated mocks
│   ├── jwt_service_mock.go
│   ├── user_service_mock.go
│   └── product_service_mock.go
│
└── testhelpers/                # Test utilities
    ├── test_helpers.go         # 247 lines
    └── integration_test.go     # 156 lines
```

---

### **9. `/web` - Web Assets**

```
web/
├── static/                     # Static files
│   ├── css/
│   ├── js/
│   └── images/
│
└── templates/                  # Templates
    └── emails/                 # Email templates
        ├── welcome.html
        ├── welcome.txt
        ├── reset-password.html
        └── reset-password.txt
```

---

### **10. Root Level Files**

```
/
├── go.mod                      # Go module
├── go.sum                      # Dependencies
├── Makefile                    # Build automation
├── README.md                   # Project overview
├── LICENSE                     # MIT License
├── CHANGELOG.md                # Version history
├── CONTRIBUTING.md             # Contribution guide
├── .gitignore                  # Git ignore
├── .editorconfig               # Editor config
├── .env                        # Local environment (gitignored)
└── .dockerignore               # Docker ignore
```

---

## 📊 STATISTICS

### Tổng quan:
- **Total Directories:** ~40
- **Total Go Files:** ~120
- **Total Lines of Code:** ~15,000
- **Test Files:** ~30
- **Documentation Files:** ~15

### Phân bố Code:

| Layer | Files | Lines | Purpose |
|-------|-------|-------|---------|
| **Handlers** | 7 | ~1,800 | HTTP endpoints |
| **Services** | 7 | ~1,500 | Business logic |
| **Repository** | 9 | ~1,200 | Data access |
| **Models** | 14 | ~1,000 | Data structures |
| **Middleware** | 6 | ~750 | HTTP middlewares |
| **Infrastructure** | 50+ | ~8,000 | Auth, DB, Email, Cache, etc. |
| **Tests** | 30+ | ~2,000 | Unit & integration tests |

---

## 🎯 MAPPING - File Hiện Tại → Vị Trí Mới

### Application Layer:
```
app/main.go                     → cmd/server/main.go
app/handlers/*.go               → internal/app/handlers/*.go
app/routes/*.go                 → internal/app/routes/*.go
app/middleware/*.go             → internal/app/middleware/*.go
```

### Domain Layer:
```
app/models/*.go                 → internal/domain/models/*.go
app/services/*.go               → internal/domain/services/*.go
app/repository/*.go             → internal/domain/repository/*.go
```

### Infrastructure:
```
app/auth/*.go                   → internal/pkg/auth/*.go
app/cache/*.go                  → internal/pkg/cache/*.go
app/config/*.go                 → internal/pkg/config/*.go
app/database/*.go               → internal/pkg/database/*.go
app/email/*.go                  → internal/pkg/email/*.go
app/storage/*.go                → internal/pkg/storage/*.go
app/health/*.go                 → internal/pkg/health/*.go
app/metrics/*.go                → internal/pkg/metrics/*.go
app/logger/*.go                 → internal/pkg/logger/*.go
app/i18n/*.go                   → internal/pkg/i18n/*.go
app/errors/*.go                 → internal/pkg/errors/*.go
app/utils/*.go                  → internal/pkg/utils/*.go
```

### API & Configs:
```
app/docs/*                      → api/openapi/*
app/.env                        → configs/.env.example
app/locales/*.json              → internal/pkg/i18n/locales/*.json
```

### Deployment:
```
Dockerfile                      → deployments/docker/Dockerfile
docker-compose.yml              → deployments/docker/docker-compose.yml
```

### Documentation:
```
*.md files                      → docs/*.md (except README.md)
```

### Data & Logs:
```
app/uploads/                    → web/static/uploads/ (hoặc giữ data/ ở root)
app/logs/                       → logs/ (ở root, gitignored)
```

---

## 🔧 Import Changes - Cheat Sheet

Khi refactor, thay đổi imports như sau:

```go
// OLD IMPORTS:
import (
    "base-golang-restful-app/handlers"
    "base-golang-restful-app/models"
    "base-golang-restful-app/services"
    "base-golang-restful-app/repository"
    "base-golang-restful-app/middleware"
    "base-golang-restful-app/auth"
    "base-golang-restful-app/config"
    "base-golang-restful-app/database"
    "base-golang-restful-app/email"
)

// NEW IMPORTS:
import (
    "base-golang-restful/internal/app/handlers"
    "base-golang-restful/internal/domain/models"
    "base-golang-restful/internal/domain/services"
    "base-golang-restful/internal/domain/repository"
    "base-golang-restful/internal/app/middleware"
    "base-golang-restful/internal/pkg/auth"
    "base-golang-restful/internal/pkg/config"
    "base-golang-restful/internal/pkg/database"
    "base-golang-restful/internal/pkg/email"
)
```

---

## 📝 .gitignore (Cập nhật)

```gitignore
# Binaries
bin/
*.exe
server

# Environment
.env
.env.local

# Logs
logs/
*.log

# Uploads
uploads/
web/static/uploads/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Test
coverage.out
coverage.html

# Build
build/

# Temporary
*.tmp
*.bak
*.old
```

---

## 🎊 KẾT QUẢ SAU REFACTOR

### Project Root sẽ trông như thế này:

```bash
$ ls -la

drwxr-xr-x  api/
drwxr-xr-x  assets/
drwxr-xr-x  cmd/
drwxr-xr-x  configs/
drwxr-xr-x  deployments/
drwxr-xr-x  docs/
drwxr-xr-x  internal/
drwxr-xr-x  scripts/
drwxr-xr-x  test/
drwxr-xr-x  web/
-rw-r--r--  .gitignore
-rw-r--r--  .editorconfig
-rw-r--r--  go.mod
-rw-r--r--  go.sum
-rw-r--r--  Makefile
-rw-r--r--  README.md
-rw-r--r--  LICENSE
-rw-r--r--  CHANGELOG.md
```

**Clean, professional, và đúng chuẩn Go! 🎉**

---

## 💡 LƯU Ý QUAN TRỌNG

1. **Backup trước khi refactor:**
   ```bash
   cp -r app app.backup
   ```

2. **Test sau mỗi bước di chuyển**

3. **Commit thường xuyên:**
   ```bash
   git add .
   git commit -m "refactor: move handlers to internal/app/handlers"
   ```

4. **Dùng IDE để refactor imports** (GoLand, VSCode có built-in support)

5. **Không rush** - làm từng layer một

---

Bây giờ bạn có toàn bộ cây cấu trúc! Bắt đầu từ bước nào bạn muốn, tôi sẽ hỗ trợ! 🚀

