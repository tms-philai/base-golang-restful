# JWT + RBAC Implementation Guide

## 📚 Tổng quan

Dự án đã được tích hợp đầy đủ **JWT Authentication** và **RBAC (Role-Based Access Control)** để quản lý xác thực và phân quyền người dùng.

---

## ✅ Các tính năng đã hoàn thiện

### 1. **JWT Authentication**
- ✅ Access Token & Refresh Token
- ✅ Token generation, validation, và refresh
- ✅ Auth middleware với `Authenticate()` và `OptionalAuthenticate()`
- ✅ Password hashing với bcrypt
- ✅ Auth endpoints: `/api/v1/auth/register`, `/login`, `/refresh`, `/profile`, `/change-password`

### 2. **RBAC (Role-Based Access Control)**
- ✅ User, Role, Permission models với many-to-many relationships
- ✅ Default roles: `admin`, `user`, `moderator`
- ✅ Default permissions cho user, role, permission resources
- ✅ RBAC middleware functions
- ✅ Auto-seeding default roles và permissions

### 3. **Integration**
- ✅ Middleware load full User với Roles và Permissions
- ✅ Auto-assign default "user" role khi đăng ký
- ✅ Routes được setup đầy đủ tại `/api/v1/auth/*`

---

## 🚀 API Endpoints

### Authentication

#### 1. Register (Đăng ký)
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### 2. Login (Đăng nhập)
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### 3. Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

#### 4. Get Profile (Yêu cầu authentication)
```http
GET /api/v1/auth/profile
Authorization: Bearer <access_token>
```

#### 5. Change Password (Yêu cầu authentication)
```http
POST /api/v1/auth/change-password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "current_password": "OldPass123!",
  "new_password": "NewPass123!"
}
```

---

## 🔐 Sử dụng RBAC Middleware

### 1. Require Authentication
```go
// Yêu cầu user phải đăng nhập
protected := r.Group("/api/v1")
protected.Use(authMiddleware.Authenticate())
{
    protected.GET("/profile", handler.GetProfile)
}
```

### 2. Require Specific Role
```go
import "base-gin/internal/app/middleware"

// Chỉ admin mới truy cập được
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", handler.ListUsers)
    admin.DELETE("/users/:id", handler.DeleteUser)
}
```

### 3. Require Any Role (OR condition)
```go
// User có role "admin" HOẶC "moderator" đều được
moderators := protected.Group("/moderate")
moderators.Use(middleware.RequireAnyRole([]string{"admin", "moderator"}))
{
    moderators.POST("/approve", handler.ApproveContent)
}
```

### 4. Require All Roles (AND condition)
```go
// User phải có CẢ HAI roles
special := protected.Group("/special")
special.Use(middleware.RequireAllRoles([]string{"admin", "moderator"}))
{
    special.POST("/action", handler.SpecialAction)
}
```

### 5. Require Permission
```go
// Dựa trên permission cụ thể
users := protected.Group("/users")
{
    users.GET("/:id", 
        middleware.RequirePermission("user.read"), 
        handler.GetUser)
        
    users.POST("/", 
        middleware.RequirePermission("user.create"), 
        handler.CreateUser)
        
    users.PUT("/:id", 
        middleware.RequirePermission("user.update"), 
        handler.UpdateUser)
        
    users.DELETE("/:id", 
        middleware.RequirePermission("user.delete"), 
        handler.DeleteUser)
}
```

### 6. Require Any Permission (OR condition)
```go
// User có BẤT KỲ permission nào trong list
readonly := protected.Group("/readonly")
readonly.Use(middleware.RequireAnyPermission([]string{"user.read", "role.read"}))
{
    readonly.GET("/data", handler.GetData)
}
```

### 7. Require All Permissions (AND condition)
```go
// User phải có TẤT CẢ permissions
admin := protected.Group("/admin")
admin.Use(middleware.RequireAllPermissions([]string{"user.create", "user.delete"}))
{
    admin.POST("/bulk-action", handler.BulkAction)
}
```

---

## 📦 Default Roles & Permissions

### Roles được tạo sẵn:

1. **admin** - Administrator
   - Có tất cả permissions
   
2. **user** - Basic User
   - `user.read` - Đọc thông tin user

3. **moderator** - Moderator
   - `user.read` - Đọc thông tin user
   - `user.list` - List users
   - `role.read` - Đọc thông tin role
   - `role.list` - List roles

### Permissions được tạo sẵn:

#### User Permissions
- `user.create` - Create new users
- `user.read` - View user details
- `user.update` - Update user information
- `user.delete` - Delete users
- `user.list` - List all users

#### Role Permissions
- `role.create` - Create new roles
- `role.read` - View role details
- `role.update` - Update role information
- `role.delete` - Delete roles
- `role.list` - List all roles

#### Permission Permissions
- `permission.read` - View permission details
- `permission.list` - List all permissions

---

## 🔧 Cách kiểm tra Role/Permission trong Code

### Trong Handler
```go
func (h *Handler) SomeAction(c *gin.Context) {
    // Lấy user từ context
    userInterface, exists := middleware.GetCurrentUser(c)
    if !exists {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }
    
    user := userInterface.(*models.User)
    
    // Kiểm tra role
    if user.HasRole("admin") {
        // Admin logic
    }
    
    // Kiểm tra permission
    if user.HasPermission("user.delete") {
        // Allow delete
    }
    
    // Kiểm tra nhiều roles
    if user.HasAnyRole([]string{"admin", "moderator"}) {
        // Logic
    }
}
```

### Helper Functions
```go
// Trong middleware package
userID, exists := middleware.GetUserID(c)
email, exists := middleware.GetUserEmail(c)
user, exists := middleware.GetCurrentUser(c)
claims, exists := middleware.GetTokenClaims(c)
```

---

## 🗄️ Database Schema

### Tables
- `users` - User accounts
- `roles` - System roles
- `permissions` - System permissions
- `user_roles` - User-Role many-to-many
- `role_permissions` - Role-Permission many-to-many
- `refresh_tokens` - Refresh token storage

### Migrations & Seeding
Migrations và seeding sẽ tự động chạy khi khởi động server:
```go
// Trong main.go
database.AutoMigrate(db)       // Tạo tables
database.SeedDefaultRoles(db)  // Seed default data
```

---

## 🧪 Testing

### Test với cURL

#### 1. Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "first_name": "Test",
    "last_name": "User"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }'
```

#### 3. Access Protected Route
```bash
# Lưu token từ login response
TOKEN="eyJhbGc..."

curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📝 Ghi chú quan trọng

1. **Password Requirements:**
   - Tối thiểu 8 ký tự
   - Phải có chữ hoa
   - Phải có chữ thường
   - Phải có số
   - Phải có ký tự đặc biệt

2. **Token Expiration:**
   - Access Token: 15 phút (có thể config)
   - Refresh Token: 7 ngày (có thể config)

3. **Default Role:**
   - User mới sẽ tự động được gán role `user`
   - Admin cần được gán thủ công hoặc qua database

4. **RBAC Context:**
   - Middleware tự động load User với Roles và Permissions
   - Không cần query thêm trong handlers

---

## 🎯 Next Steps (Tùy chọn mở rộng)

### 1. Thêm Role Management Endpoints
```go
// Tạo handlers cho role management
POST   /api/v1/roles              - Create role
GET    /api/v1/roles              - List roles
GET    /api/v1/roles/:id          - Get role details
PUT    /api/v1/roles/:id          - Update role
DELETE /api/v1/roles/:id          - Delete role
POST   /api/v1/roles/:id/permissions - Assign permissions
```

### 2. Thêm User Role Assignment
```go
POST   /api/v1/users/:id/roles    - Assign role to user
DELETE /api/v1/users/:id/roles/:roleId - Remove role from user
```

### 3. Thêm Permission Management
```go
GET    /api/v1/permissions        - List permissions
POST   /api/v1/permissions        - Create permission
```

### 4. Custom Error Handlers
Tùy chỉnh error responses trong RBAC middleware:
```go
config := middleware.RBACConfig{
    UnauthorizedHandler: func(c *gin.Context) {
        c.JSON(401, gin.H{"error": "Please login"})
        c.Abort()
    },
    ForbiddenHandler: func(c *gin.Context) {
        c.JSON(403, gin.H{"error": "Access denied"})
        c.Abort()
    },
}

protected.Use(middleware.RequireRole("admin", config))
```

---

## 📞 Support

Nếu có vấn đề hoặc câu hỏi:
1. Kiểm tra logs trong terminal
2. Kiểm tra Swagger docs tại: `http://localhost:8080/swagger/index.html`
3. Kiểm tra database có đầy đủ tables và data

---

**Chúc bạn coding vui vẻ! 🚀**
