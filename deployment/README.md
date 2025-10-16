# 🐳 Docker Deployment Guide

Hướng dẫn triển khai dự án Base Golang RESTful API sử dụng Docker và Docker Compose.

## 📋 Mục lục

- [Cấu trúc thư mục](#cấu-trúc-thư-mục)
- [Các file Docker](#các-file-docker)
- [Hướng dẫn sử dụng](#hướng-dẫn-sử-dụng)
- [Cấu hình môi trường](#cấu-hình-môi-trường)
- [Troubleshooting](#troubleshooting)

## 📁 Cấu trúc thư mục

```
deployment/
├── Dockerfile                    # Docker image cho ứng dụng
├── docker-compose.yml           # Production environment
├── docker-compose.dev.yml       # Development environment
├── docker-compose.email-test.yml # Email testing environment
├── nginx.conf                   # Nginx reverse proxy config
└── README.md                    # Hướng dẫn này
```

## 🐳 Các file Docker

### 1. **Dockerfile**
- Multi-stage build với Go 1.24
- Security: chạy với non-root user
- Health check tích hợp
- Optimized cho production

### 2. **docker-compose.yml** (Production)
- PostgreSQL database
- Redis cache (optional)
- Nginx reverse proxy (optional)
- Health checks và restart policies
- Volume persistence

### 3. **docker-compose.dev.yml** (Development)
- Hot reload với volume mounting
- MailHog cho email testing
- Database riêng cho dev
- Debug-friendly configuration

### 4. **docker-compose.email-test.yml** (Email Testing)
- MailHog và Mailpit cho email testing
- Database test riêng biệt
- Cấu hình email testing

## 🚀 Hướng dẫn sử dụng

### **Development Environment**

```bash
# Chạy development environment
cd deployment
docker-compose -f docker-compose.dev.yml up -d

# Xem logs
docker-compose -f docker-compose.dev.yml logs -f app

# Dừng services
docker-compose -f docker-compose.dev.yml down
```

**Truy cập:**
- API: http://localhost:8001
- Swagger: http://localhost:8001/swagger/index.html
- MailHog UI: http://localhost:8025
- Database: localhost:5432

### **Production Environment**

```bash
# Chạy production environment
cd deployment
docker-compose up -d

# Xem logs
docker-compose logs -f app

# Dừng services
docker-compose down
```

**Truy cập:**
- API: http://localhost:8001
- Swagger: http://localhost:8001/swagger/index.html
- Nginx: http://localhost:80 (nếu enable)

### **Email Testing Environment**

```bash
# Chạy email testing environment
cd deployment
docker-compose -f docker-compose.email-test.yml up -d

# Test email functionality
curl -X POST http://localhost:8001/api/v1/email/send \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to": "test@example.com", "subject": "Test", "body": "Hello"}'
```

**Truy cập:**
- API: http://localhost:8001
- MailHog UI: http://localhost:8025
- Mailpit UI: http://localhost:8026

## ⚙️ Cấu hình môi trường

### **Environment Variables**

Tạo file `.env` trong thư mục `deployment/`:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8001

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=base_gin

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key

# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com
```

### **Database Initialization**

Tạo thư mục `init-scripts/` và thêm file SQL:

```sql
-- init-scripts/01-init.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### **Nginx Configuration**

File `nginx.conf` đã được cấu hình sẵn với:
- Rate limiting
- CORS headers
- Security headers
- Gzip compression
- SSL support (uncomment để enable)

## 🔧 Commands hữu ích

### **Docker Commands**

```bash
# Build image
docker build -f deployment/Dockerfile -t base-gin .

# Run container
docker run -p 8001:8001 base-gin

# Xem logs
docker logs -f base-gin-app

# Vào container
docker exec -it base-gin-app sh

# Xóa containers và volumes
docker-compose down -v
```

### **Database Commands**

```bash
# Backup database
docker exec base-gin-postgres pg_dump -U postgres base_gin > backup.sql

# Restore database
docker exec -i base-gin-postgres psql -U postgres base_gin < backup.sql

# Vào database
docker exec -it base-gin-postgres psql -U postgres -d base_gin
```

### **Development Commands**

```bash
# Hot reload (development)
docker-compose -f docker-compose.dev.yml up --build

# Rebuild specific service
docker-compose -f docker-compose.dev.yml up --build app

# Xem resource usage
docker stats

# Clean up
docker system prune -a
```

## 🐛 Troubleshooting

### **Common Issues**

#### 1. **Port already in use**
```bash
# Kiểm tra port đang sử dụng
lsof -i :8001
lsof -i :5432

# Thay đổi port trong docker-compose.yml
ports:
  - "8002:8001"  # Thay vì 8001:8001
```

#### 2. **Database connection failed**
```bash
# Kiểm tra database container
docker-compose logs postgres

# Kiểm tra network
docker network ls
docker network inspect deployment_app-network
```

#### 3. **Permission denied**
```bash
# Fix permissions
sudo chown -R $USER:$USER .
chmod -R 755 .
```

#### 4. **Out of disk space**
```bash
# Clean up Docker
docker system prune -a
docker volume prune
```

### **Health Checks**

```bash
# Kiểm tra health status
docker-compose ps

# Test API health
curl http://localhost:8001/health

# Test database connection
docker exec base-gin-postgres pg_isready -U postgres
```

### **Logs và Debugging**

```bash
# Xem logs của tất cả services
docker-compose logs

# Xem logs của service cụ thể
docker-compose logs app
docker-compose logs postgres

# Follow logs real-time
docker-compose logs -f app

# Debug container
docker exec -it base-gin-app sh
```

## 📊 Monitoring

### **Resource Monitoring**

```bash
# Xem resource usage
docker stats

# Xem disk usage
docker system df

# Xem volume usage
docker volume ls
docker volume inspect deployment_postgres-data
```

### **Application Monitoring**

- Health check endpoint: `/health`
- Swagger documentation: `/swagger/index.html`
- Database metrics: PostgreSQL logs
- Application logs: Container logs

## 🔒 Security Best Practices

1. **Environment Variables**: Không commit file `.env`
2. **Secrets**: Sử dụng Docker secrets hoặc external secret management
3. **Network**: Sử dụng custom networks
4. **User**: Chạy container với non-root user
5. **Updates**: Thường xuyên update base images
6. **SSL**: Enable HTTPS trong production
7. **Firewall**: Cấu hình firewall rules

## 📚 Tài liệu tham khảo

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Redis Docker Image](https://hub.docker.com/_/redis)
- [Nginx Docker Image](https://hub.docker.com/_/nginx)

---

**Chúc bạn deploy thành công! 🚀**
