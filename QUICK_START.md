# 🚀 Quick Start Guide

## Base Golang RESTful API

Hướng dẫn nhanh để chạy và test dự án.

---

## 📋 Prerequisites

### Required
- Go 1.21 or higher
- PostgreSQL 15+ (for database features)
- Redis 7+ (for caching features)

### Optional
- Docker & Docker Compose (recommended)
- AWS Account (for S3 storage)
- GCP Account (for GCS storage)
- SMTP Server (for email features)

---

## 🐳 Quick Start with Docker (Recommended)

### 1. Clone Repository
\`\`\`bash
git clone https://github.com/tms-philai/base-golang-restful.git
cd base-golang-restful
\`\`\`

### 2. Setup Environment
\`\`\`bash
cp env.example .env
# Edit .env if needed
\`\`\`

### 3. Start with Docker Compose
\`\`\`bash
docker-compose up -d
\`\`\`

### 4. Check Status
\`\`\`bash
# Check health
curl http://localhost:8080/health

# Check metrics
curl http://localhost:8080/metrics

# Test API
curl http://localhost:8080/api/v1/test
\`\`\`

### 5. Stop Services
\`\`\`bash
docker-compose down
\`\`\`

---

## 💻 Local Development Setup

### 1. Install Go
\`\`\`bash
# Download from https://golang.org/dl/
# Or use package manager:

# Windows (with Chocolatey)
choco install golang

# macOS (with Homebrew)
brew install go

# Linux (Ubuntu/Debian)
sudo apt-get install golang-go
\`\`\`

### 2. Verify Installation
\`\`\`bash
go version
# Should output: go version go1.21.x
\`\`\`

### 3. Install Dependencies
\`\`\`bash
cd base-golang-restful
go mod download
\`\`\`

### 4. Setup PostgreSQL
\`\`\`bash
# Using Docker
docker run -d \\
  --name postgres \\
  -p 5432:5432 \\
  -e POSTGRES_USER=postgres \\
  -e POSTGRES_PASSWORD=postgres \\
  -e POSTGRES_DB=base_golang_restful \\
  postgres:15-alpine

# Or install locally from https://www.postgresql.org/download/
\`\`\`

### 5. Setup Redis
\`\`\`bash
# Using Docker
docker run -d \\
  --name redis \\
  -p 6379:6379 \\
  redis:7-alpine

# Or install locally from https://redis.io/download
\`\`\`

### 6. Configure Environment
\`\`\`bash
cp env.example .env

# Edit .env with your settings
nano .env  # or use your favorite editor
\`\`\`

### 7. Run Application
\`\`\`bash
go run main.go
\`\`\`

### 8. Access Application
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics
- **Swagger Docs**: http://localhost:8080/swagger/index.html

---

## 🧪 Running Tests

### Run All Tests
\`\`\`bash
go test ./app/... -v
\`\`\`

### Run with Coverage
\`\`\`bash
go test ./app/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
\`\`\`

### Run Specific Tests
\`\`\`bash
# Database tests
go test ./app/database/... -v

# Auth tests
go test ./app/auth/... -v

# Cache tests
go test ./app/cache/... -v
\`\`\`

---

## 📡 API Testing

### Using cURL

#### Health Check
\`\`\`bash
curl http://localhost:8080/health
\`\`\`

#### Test Endpoint
\`\`\`bash
curl http://localhost:8080/api/v1/test
\`\`\`

#### i18n Test (English)
\`\`\`bash
curl http://localhost:8080/api/v1/i18n-test
\`\`\`

#### i18n Test (Vietnamese)
\`\`\`bash
curl -H "Accept-Language: vi" http://localhost:8080/api/v1/i18n-test
\`\`\`

#### Echo Test
\`\`\`bash
curl -X POST http://localhost:8080/api/v1/echo \\
  -H "Content-Type: application/json" \\
  -d '{"message": "Hello World", "test": true}'
\`\`\`

#### API Versioning Test
\`\`\`bash
# Version in header
curl -H "API-Version: v2" http://localhost:8080/api/v1/test

# Version in query
curl "http://localhost:8080/api/v1/test?version=v2"
\`\`\`

### Using Postman

1. Import the API collection (if available)
2. Set base URL: `http://localhost:8080`
3. Test endpoints:
   - GET `/health`
   - GET `/metrics`
   - GET `/api/v1/test`
   - POST `/api/v1/echo`

---

## 🔧 Development Tools

### Generate Swagger Docs
\`\`\`bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init

# Docs will be available at /swagger/index.html
\`\`\`

### Format Code
\`\`\`bash
go fmt ./...
\`\`\`

### Lint Code
\`\`\`bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
\`\`\`

### Run with Hot Reload
\`\`\`bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
\`\`\`

---

## 🐛 Troubleshooting

### Go Command Not Found
\`\`\`bash
# Add Go to PATH
# Windows: Add C:\\Go\\bin to System PATH
# Linux/Mac: Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/usr/local/go/bin
\`\`\`

### Port Already in Use
\`\`\`bash
# Change port in .env
APP_PORT=8081

# Or kill process using port 8080
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
\`\`\`

### Database Connection Error
\`\`\`bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U postgres -d base_golang_restful

# Verify .env settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=base_golang_restful
\`\`\`

### Redis Connection Error
\`\`\`bash
# Check Redis is running
docker ps | grep redis

# Test connection
redis-cli ping
# Should return: PONG

# Verify .env settings
REDIS_HOST=localhost
REDIS_PORT=6379
\`\`\`

---

## 📚 Next Steps

1. **Read Documentation**
   - `README.md` - Project overview
   - `ROADMAP.md` - Development roadmap
   - `PROMPTS.md` - Implementation details
   - `PROJECT_COMPLETION_SUMMARY.md` - Feature list

2. **Explore Features**
   - Test all API endpoints
   - Try different languages (EN/VI)
   - Check health and metrics
   - Review Swagger documentation

3. **Customize**
   - Add your own endpoints
   - Implement business logic
   - Configure for your needs
   - Deploy to production

---

## 🎯 Production Deployment

### Build for Production
\`\`\`bash
# Build binary
go build -o app main.go

# Run binary
./app
\`\`\`

### Docker Deployment
\`\`\`bash
# Build image
docker build -t base-golang-restful:latest .

# Run container
docker run -d -p 8080:8080 base-golang-restful:latest
\`\`\`

### Cloud Deployment
See `GOOGLE_CLOUD_DEPLOYMENT.md` for GCP deployment guide.

---

## 💡 Tips

- Use Docker Compose for easiest setup
- Check logs in `./logs/app.log`
- Monitor metrics at `/metrics`
- Use Swagger UI for API exploration
- Run tests before committing
- Keep `.env` file secure (never commit)

---

## 🆘 Need Help?

- Check `TEST_RESULTS.md` for test documentation
- Review `PROJECT_COMPLETION_SUMMARY.md` for features
- See `ROADMAP.md` for architecture details
- Contact: support@example.com

---

**Happy Coding! 🚀**
