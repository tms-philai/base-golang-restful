# Base Golang RESTful API with Gin Framework

A clean, production-ready RESTful API built with Go and the Gin web framework.

## Features

- High-performance HTTP web framework (Gin)
- Clean project structure
- RESTful API endpoints
- JSON request/response handling
- Middleware support
- Error handling

## Prerequisites

- Go 1.19 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <your-repository-url>
cd base-golang-restful
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/ping` - Returns a simple pong response

### User Management
- **GET** `/hello/:name` - Greets a user by name
- **POST** `/users` - Creates a new user

## Example Usage

### Health Check
```bash
curl http://localhost:8080/ping
```

Response:
```json
{
  "message": "pong"
}
```

### Greet User
```bash
curl http://localhost:8080/hello/John
```

Response:
```json
{
  "message": "Hello John"
}
```

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

Response:
```json
{
  "message": "User created successfully",
  "user": {
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

## Project Structure

```
.
├── main.go              # Application entry point
├── go.mod              # Go module file
├── go.sum              # Go dependencies checksum
├── .cursorrules        # Cursor IDE rules for Go development
├── binding/            # Request binding and validation
├── codec/              # JSON encoding/decoding
├── internal/           # Private application code
├── render/             # Response rendering
└── README.md           # This file
```

## Development

This project follows Go best practices and conventions. The `.cursorrules` file contains development guidelines for consistent code style.

### Building

```bash
go build -o app main.go
```

### Running
```bash
./app
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

This project is open source and available under the MIT License.
