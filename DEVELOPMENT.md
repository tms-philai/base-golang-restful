# 🚀 Hướng Dẫn Development - Base Golang RESTful API

## ✅ Đã Setup

Hệ thống của bạn đã được cấu hình và đang chạy với:

- **Go**: 1.25.1 ✅
- **PostgreSQL**: Docker container (port 5433) ✅
- **Redis**: Docker container (port 6379) ✅
- **Server**: Running on port 8080 ✅

---

## 📋 Quick Start

### 1. Sử dụng Script Tiện Ích (Khuyến nghị)

```bash
# Start tất cả (DB + Redis + App)
./scripts/dev.sh start

# Stop tất cả
./scripts/dev.sh stop

# Restart app (sau khi sửa code)
./scripts/dev.sh restart

# Xem status
./scripts/dev.sh status

# Xem logs
./scripts/dev.sh logs

# Xem tất cả lệnh
./scripts/dev.sh help
```

### 2. Quản Lý Thủ Công

#### Start Services

```bash
# Start PostgreSQL
docker run -d --name postgres-dev \
  -p 5433:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=base_golang_restful \
  postgres:15-alpine

# Start Redis
docker run -d --name redis-dev \
  -p 6379:6379 \
  redis:7-alpine
```

#### Build & Run App

```bash
cd app

# Build
go build -o bin/server main_simple.go

# Run foreground
./bin/server

# Run background
./bin/server > logs/server.log 2>&1 &

# Stop
pkill -f "bin/server"
```

---

## 🌐 Endpoints Có Sẵn

| Endpoint | URL | Mô tả |
|----------|-----|-------|
| **Health Check** | http://localhost:8080/health | Kiểm tra trạng thái hệ thống |
| **Swagger UI** | http://localhost:8080/swagger/index.html | API Documentation |
| **API Test** | http://localhost:8080/api/v1/test | Test endpoint |
| **API Ping** | http://localhost:8080/api/v1/ping | Ping pong |

---

## 🧪 Test API

### Sử dụng cURL

```bash
# Health check
curl http://localhost:8080/health

# Test API
curl http://localhost:8080/api/v1/test

# Test với pretty print
curl -s http://localhost:8080/health | python3 -m json.tool
```

### Sử dụng Browser

Mở trình duyệt và truy cập:
- Swagger: http://localhost:8080/swagger/index.html
- Health: http://localhost:8080/health

---

## 📁 Cấu Trúc Dự Án

```
base-golang-restful/
├── app/                       # Main application
│   ├── main.go               # Full app (có lỗi, đang fix)
│   ├── main_simple.go        # Simple version (đang dùng) ✅
│   ├── bin/                  # Binary files
│   │   └── server           # Compiled server
│   ├── logs/                 # Application logs
│   │   └── server.log
│   ├── uploads/              # File uploads
│   ├── config/               # Configuration
│   ├── models/               # Data models
│   ├── services/             # Business logic
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # Middlewares
│   ├── repository/           # Data access
│   ├── auth/                 # Authentication
│   ├── cache/                # Caching
│   ├── database/             # Database utilities
│   ├── docs/                 # Swagger docs
│   ├── email/                # Email service
│   ├── errors/               # Error handling
│   ├── i18n/                 # Internationalization
│   ├── logger/               # Logging
│   └── utils/                # Utilities
├── scripts/
│   └── dev.sh               # Development script
├── .env                      # Environment config ✅
├── docker-compose.yml        # Docker compose (không dùng)
├── Dockerfile                # Dockerfile (không dùng)
└── DEVELOPMENT.md            # Tài liệu này
```

---

## ⚙️ Cấu Hình (.env)

File `.env` trong root directory:

```bash
# Application
APP_PORT=8080
APP_ENV=development

# Database (PostgreSQL trên Docker)
DB_HOST=localhost
DB_PORT=5433              # ⚠️ Note: 5433 vì tránh conflict
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=base_golang_restful

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=7d

# Storage
STORAGE_TYPE=local
STORAGE_BASE_PATH=./uploads
```

---

## 🔧 Development Workflow

### 1. Sửa Code

```bash
# Edit files trong app/
vim app/services/user_service.go
```

### 2. Rebuild & Restart

```bash
# Option 1: Dùng script
./scripts/dev.sh restart

# Option 2: Manual
cd app
pkill -f "bin/server"
go build -o bin/server main_simple.go
./bin/server &
```

### 3. Test

```bash
curl http://localhost:8080/api/v1/test
```

### 4. Xem Logs

```bash
tail -f app/logs/server.log
```

---

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Tìm process đang dùng port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Connection Error

```bash
# Kiểm tra PostgreSQL container
docker ps | grep postgres-dev

# Xem logs
docker logs postgres-dev

# Restart
docker restart postgres-dev
```

### Redis Connection Error

```bash
# Kiểm tra Redis container
docker ps | grep redis-dev

# Test connection
docker exec -it redis-dev redis-cli ping
# Kết quả: PONG
```

### Server Không Start

```bash
# Xem logs chi tiết
cat app/logs/server.log

# Kiểm tra Go dependencies
cd app && go mod tidy
```

---

## 📊 Database Management

### Truy cập PostgreSQL

```bash
# Sử dụng psql trong container
docker exec -it postgres-dev psql -U postgres -d base_golang_restful

# Hoặc từ host (nếu có psql)
psql -h localhost -p 5433 -U postgres -d base_golang_restful
```

### Truy cập Redis

```bash
# Sử dụng redis-cli trong container
docker exec -it redis-dev redis-cli

# Các lệnh Redis hữu ích
> PING
> KEYS *
> GET key_name
> FLUSHALL  # Xóa tất cả keys (cẩn thận!)
```

---

## 🚀 Production Build

Khi sẵn sàng deploy production:

```bash
cd app

# Build với optimizations
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -o bin/server-prod \
  main.go  # Sử dụng main.go đầy đủ khi đã fix

# File binary nhỏ hơn, tối ưu cho production
```

---

## 📝 Notes

### Hiện Tại Đang Dùng

- ✅ **main_simple.go**: Version đơn giản, ổn định
- ❌ **main.go**: Version đầy đủ, có một số lỗi handlers

### Cần Sửa (TODO)

1. Fix handlers trong `main.go` để dùng full features
2. Thêm database migrations
3. Implement authentication endpoints
4. Thêm tests

### Services Đang Chạy

```bash
# Kiểm tra
./scripts/dev.sh status

# Kết quả:
# 🐳 Docker Containers:
#    postgres-dev: Up X minutes
#    redis-dev: Up X minutes
# 🚀 Go Server:
#    ✅ Running (PID: XXXX)
```

---

## 💡 Tips

1. **Hot Reload**: Cài `air` để auto-reload khi sửa code
   ```bash
   go install github.com/cosmtrek/air@latest
   cd app && air
   ```

2. **Database GUI**: Sử dụng TablePlus, DBeaver, hoặc pgAdmin
   - Host: localhost
   - Port: 5433
   - User: postgres
   - Pass: postgres
   - DB: base_golang_restful

3. **API Testing**: Dùng Postman hoặc Insomnia
   - Import Swagger JSON từ `/swagger/doc.json`

4. **Logs**: Luôn theo dõi logs khi develop
   ```bash
   ./scripts/dev.sh logs
   ```

---

## 🎯 Quick Commands

```bash
# Start everything
./scripts/dev.sh start

# Check status
./scripts/dev.sh status

# Watch logs
./scripts/dev.sh logs

# Restart after code changes
./scripts/dev.sh restart

# Stop everything
./scripts/dev.sh stop

# Clean up containers
./scripts/dev.sh services:rm
```

---

**Happy Coding! 🚀**

Nếu có vấn đề gì, check logs hoặc chạy `./scripts/dev.sh status` để xem hệ thống.

