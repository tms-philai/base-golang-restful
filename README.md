# Base Golang RESTful API

Dự án RESTful API cơ bản được xây dựng bằng Go và Gin framework, được tùy chỉnh từ source code gốc của [gin-gonic/gin](https://github.com/gin-gonic/gin).

## 🚀 Tính năng

- **High-performance**: Sử dụng Gin framework với hiệu suất cao
- **RESTful API**: Các endpoint API tuân thủ chuẩn REST
- **JSON Support**: Xử lý request/response JSON
- **Middleware**: Hỗ trợ middleware cho logging, recovery, CORS, v.v.
- **Validation**: Validation dữ liệu đầu vào với struct tags
- **Multiple Formats**: Hỗ trợ JSON, XML, YAML, TOML, Protobuf
- **Clean Architecture**: Cấu trúc code rõ ràng, dễ bảo trì

## 📋 Yêu cầu hệ thống

- **Go**: 1.23.0 hoặc cao hơn
- **Git**: Để clone và quản lý source code

## 🛠️ Cài đặt

1. **Clone repository:**
```bash
git clone <your-repository-url>
cd base-golang-restful
```

2. **Cài đặt dependencies:**
```bash
go mod tidy
```

3. **Chạy ứng dụng:**
```bash
go run main.go
```

Server sẽ khởi động tại `http://localhost:8080`

## 📚 API Endpoints

### Health Check
```http
GET /ping
```
**Response:**
```json
{
  "message": "pong"
}
```

### Greeting
```http
GET /hello/:name
```
**Parameters:**
- `name` (string): Tên người dùng

**Response:**
```json
{
  "message": "Hello John"
}
```

### User Management
```http
POST /users
```
**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "user": {
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

## 🧪 Test API

### Sử dụng curl

**Health Check:**
```bash
curl http://localhost:8080/ping
```

**Greeting:**
```bash
curl http://localhost:8080/hello/John
```

**Create User:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### Sử dụng Postman
Import các endpoint trên vào Postman để test dễ dàng hơn.

## 📁 Cấu trúc dự án

```
base-golang-restful/
├── main.go                 # Entry point của ứng dụng
├── go.mod                  # Go module dependencies
├── go.sum                  # Checksum của dependencies
├── .cursorrules           # Quy tắc phát triển cho Cursor IDE
├── README.md              # Tài liệu dự án
│
├── binding/               # Request binding và validation
│   ├── binding.go         # Core binding logic
│   ├── json.go           # JSON binding
│   ├── form.go           # Form binding
│   ├── xml.go            # XML binding
│   └── ...               # Các format khác
│
├── codec/                 # Encoding/Decoding
│   └── json/             # JSON codec implementations
│       ├── api.go        # JSON API interface
│       ├── sonic.go      # Sonic JSON (high performance)
│       └── ...           # Các implementation khác
│
├── render/               # Response rendering
│   ├── json.go          # JSON rendering
│   ├── xml.go           # XML rendering
│   ├── html.go          # HTML rendering
│   └── ...              # Các format khác
│
├── internal/            # Private packages
│   ├── bytesconv/       # Byte conversion utilities
│   └── fs/              # File system utilities
│
└── Core Gin Files       # Các file core của Gin framework
    ├── gin.go           # Main Gin engine
    ├── context.go       # Request context
    ├── routergroup.go   # Router group
    └── ...              # Các file core khác
```

## 🔧 Development

### Build ứng dụng
```bash
go build -o app main.go
```

### Chạy binary
```bash
./app
```

### Format code
```bash
go fmt ./...
```

### Kiểm tra lỗi
```bash
go vet ./...
```

## 🎯 Tính năng nâng cao

Dự án này bao gồm toàn bộ source code của Gin framework, cho phép:

- **Custom Middleware**: Tạo middleware tùy chỉnh
- **Multiple Binding**: Hỗ trợ nhiều format binding (JSON, XML, Form, etc.)
- **Custom Rendering**: Render response theo nhiều format
- **Performance Optimization**: Tối ưu hóa hiệu suất với các codec khác nhau
- **Extensibility**: Dễ dàng mở rộng chức năng

## 📖 Tài liệu tham khảo

- [Gin Documentation](https://gin-gonic.com/)
- [Go Documentation](https://golang.org/doc/)
- [RESTful API Design](https://restfulapi.net/)

## 🤝 Đóng góp

1. Fork repository
2. Tạo feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Tạo Pull Request

## 📝 Ghi chú

- Dự án này được tùy chỉnh từ source code gốc của Gin framework
- Đã loại bỏ các file test và documentation không cần thiết
- Bao gồm `.cursorrules` để hỗ trợ development với Cursor IDE
- Tuân thủ Go best practices và coding conventions

## 📄 License

Dự án này sử dụng MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.