# Phase 5: Testing Strategy (Unit + Integration)

## 📋 Mục lục

1. [Tổng quan](#1-tổng-quan)
2. [Cấu trúc thư mục test](#2-cấu-trúc-thư-mục-test)
3. [Chuẩn bị môi trường test](#3-chuẩn-bị-môi-trường-test)
4. [Chạy test](#4-chạy-test)
5. [Unit Test Guidelines](#5-unit-test-guidelines)
6. [Integration Test Guidelines](#6-integration-test-guidelines)
7. [Helpers & Utilities](#7-helpers--utilities)
8. [Coverage & Báo cáo](#8-coverage--báo-cáo)
9. [CI/CD tích hợp test](#9-cicd-tích-hợp-test)
10. [Troubleshooting](#10-troubleshooting)
11. [Gợi ý mở rộng](#11-gợi-ý-mở-rộng)

---

## 1. Tổng quan

Phase này chuẩn hoá cách viết và chạy test cho dự án, bám sát cấu trúc thực tế trong `test/`. Bao gồm unit test cho handlers và integration test kiểm tra end-to-end với DB thật.

Mục tiêu:
- Cấu trúc, quy ước và công cụ test nhất quán
- Hướng dẫn setup DB test, sinh `.env.test` tự động
- Mẫu chạy test bằng Makefile và lệnh Go
- Mục tiêu coverage, cách tạo báo cáo

---

## 2. Cấu trúc thư mục test

Nguồn thực tế: `test/README.md`

```
test/
  unit/                    # Unit tests for individual components
    auth_handler_test.go
    user_handler_test.go
    product_handler_test.go
    email_handler_test.go
    file_handler_test.go
    monitoring_handler_test.go
  integration/             # Integration tests for API endpoints
    auth_integration_test.go
    user_integration_test.go
    product_integration_test.go
  helpers/                 # Test helper functions and utilities
    test_helpers.go
  config/                  # Test configuration
    test_config.go
  setup/                   # Test setup and teardown
    database_setup.go
  dependencies.go          # Ensure test deps
  README.md
```

---

## 3. Chuẩn bị môi trường test

### Database PostgreSQL

Tạo DB test và user:
```sql
CREATE DATABASE test_base_gin;
CREATE USER test_user WITH PASSWORD 'test_password';
GRANT ALL PRIVILEGES ON DATABASE test_base_gin TO test_user;
```

### Tự động sinh `.env.test`

`test/config/test_config.go` tự tạo `configs/.env.test` nếu chưa có với cấu hình mặc định (DB/JWT/SMTP/Storage...). Không cần thao tác tay nếu chạy test qua Makefile.

### Makefile hỗ trợ

Sử dụng các target:
- `make test-setup`: tạo thư mục cần thiết (configs, tmp)
- `make test-db-setup`: in hướng dẫn khởi tạo DB
- `make dev-setup`: tải deps + test-setup + test-db-setup

---

## 4. Chạy test

### Qua Makefile

```bash
make test                 # Tất cả tests
make test-unit            # Chỉ unit
make test-integration     # Chỉ integration
make test-coverage        # Tất cả + coverage.html
make test-unit-coverage   # Unit + coverage.html
make test-integration-coverage # Integration + coverage.html
make test-race            # Phát hiện race condition
make test-verbose         # Verbose output
```

### Trực tiếp bằng Go

```bash
go test ./...                    # Tất cả
go test ./test/unit/...          # Unit
go test ./test/integration/...   # Integration
go test -v -run TestAuthHandler_Register ./test/unit
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## 5. Unit Test Guidelines

Phạm vi: kiểm thử handler/service riêng lẻ với mock/stub. Ví dụ: `test/unit/file_handler_test.go` dùng `testify/mock`.

Nguyên tắc:
- Dùng table-driven tests cho nhiều scenario
- Dùng `gin.SetMode(gin.TestMode)` và `httptest.NewRecorder()`
- Mock service dependencies với `testify/mock`
- Kiểm tra status code, body, header, và schema response (dùng helpers)
- Test edge cases: thiếu auth, payload invalid, lỗi service, resource không tồn tại, conflict...

Snippet khởi tạo (rút gọn):
```go
gin.SetMode(gin.TestMode)
mockSvc := new(MockUserService)
h := handlers.NewUserHandler(mockSvc)
r := gin.New()
r.POST("/users", h.CreateUser)
w := httptest.NewRecorder()
req, _ := helpers.CreateJSONRequest("POST", "/users", payload)
r.ServeHTTP(w, req)
helpers.AssertJSONResponse(t, w, http.StatusCreated)
```

---

## 6. Integration Test Guidelines

Phạm vi: kiểm thử flow API hoàn chỉnh với DB thật. Sử dụng setup trong `test/setup/database_setup.go`:

- `TestMain` chạy trước/sau toàn bộ test: connect DB, migrate, seed roles/permissions, cleanup
- `SetupTestEnvironment`/`CleanupTestEnvironment`: chuẩn bị môi trường và dọn dẹp giữa các test
- `TruncateTestTables`: đảm bảo môi trường sạch giữa test cases

Gợi ý flow:
- Auth: register -> login -> refresh -> profile -> change-password
- User: admin list/create/delete; user self-update; forbidden cases
- Product: create (auth), update/delete (admin/owner), list/search/sort/paginate
- File: upload -> get -> download -> delete (owner)

Lưu ý:
- Sử dụng token thật từ login (không mock) để test middleware
- Đảm bảo dọn dẹp dữ liệu tạo mới trong teardown

---

## 7. Helpers & Utilities

Nguồn: `test/helpers/test_helpers.go`

- Builders: `CreateTestUser()`, `CreateTestAdmin()`, `CreateTestProduct()`
- Request: `CreateJSONRequest()`, `CreateMultipartRequest()`
- Assertions: `AssertJSONResponse()`, `AssertErrorResponse()`, `AssertSuccessResponse()`, `AssertAuthResponse()`, `AssertUserResponse()`, `AssertProductResponse()`, `AssertFileResponse()`, `AssertPaginationResponse()`
- Middleware mocks: `MockAuthMiddleware()`, `MockOptionalAuthMiddleware()`, `MockAdminMiddleware()`
- Misc: `SetupTestRouter()`, `GenerateTestJWT()` (placeholder), helpers con trỏ `StringPtr`, `Float64Ptr`

---

## 8. Coverage & Báo cáo

Mục tiêu coverage (tham khảo trong `test/README.md`):
- Handlers > 90%
- Services > 85%
- Repositories > 80%
- Overall > 85%

Sinh report:
```bash
make test-coverage
# Tạo coverage.html ở root, mở để xem chi tiết
```

Hiển thị nhanh trên terminal:
```bash
make test-coverage-term
```

---

## 9. CI/CD tích hợp test

Workflow mẫu (đã có trong README test): sử dụng service `postgres`, chạy `make ci-test` (deps + setup + toàn bộ tests với coverage).

Tips CI:
- Cache module Go để tăng tốc
- Đảm bảo health-check Postgres trước khi chạy test
- Xuất artifact `coverage.html` nếu cần

---

## 10. Troubleshooting

- Kết nối DB: kiểm tra Postgres chạy và tài khoản DB đúng
- Cổng 8001: tránh trùng cổng khi boot server trong integration
- Quyền thư mục: `LOCAL_STORAGE_PATH` (mặc định `/tmp/test_files`) có quyền ghi
- Env: chắc chắn `.env.test` đã được tạo (Makefile `test-setup` hoặc auto qua test_config)
- Race condition: sử dụng `make test-race`

---

## 11. Gợi ý mở rộng

- Bổ sung integration test cho File endpoints (multipart upload/cleanup)
- Thêm test cho RBAC (role/permission) ở từng route
- Tạo mocks tự động bằng `mockgen` cho services/repositories
- Thêm benchmark (đã có target `make benchmark`)
- Báo cáo JUnit/XML cho CI, thu thập coverage đa gói

Chúc bạn đạt coverage “xanh mướt” và test vững chắc! ✅
