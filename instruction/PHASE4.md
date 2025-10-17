# Phase 4: Email & Notification APIs

## 📋 Mục lục

1. [Tổng quan](#1-tổng-quan)
2. [Email API](#2-email-api)
3. [Notification API](#3-notification-api)
4. [Observer Pattern: Khái niệm & Ứng dụng](#4-observer-pattern-khái-niệm--ứng-dụng)
5. [Khởi tạo & Wiring trong dự án](#5-khởi-tạo--wiring-trong-dự-án)
6. [Xử lý lỗi & bảo mật](#6-xử-lý-lỗi--bảo-mật)
7. [Ví dụ gọi API](#7-ví-dụ-gọi-api)
8. [Kiểm thử liên quan](#8-kiểm-thử-liên-quan)
9. [Gợi ý mở rộng](#9-gợi-ý-mở-rộng)

---

## 1. Tổng quan

Phase này tập trung vào hai nhóm API chính: Email và Notification. Email API phục vụ gửi email đơn, đồng bộ/bất đồng bộ, và theo lô. Notification API cung cấp thông báo theo nhiều kênh (email, in-app) dựa trên Observer Pattern để dễ mở rộng và tách rời kênh gửi.

Mã nguồn chính:
- Email handler: `internal/app/handlers/email_handler.go`
- Notification handler: `internal/app/handlers/notification_handler.go`
- Observers (các kênh): `internal/app/observers/*.go`
- Notification service (Subject): `internal/domain/services/notification_service.go`
- Email service (queue, worker): `internal/domain/services/email_service.go`
- Models: `internal/domain/models/email.go`, `internal/domain/models/notification.go`

---

## 2. Email API

Tất cả endpoints cần Bearer token và được group dưới `/api/v1/email` (xem `routes.setupEmailRoutes`).

Các endpoint hỗ trợ:

- POST `/api/v1/email/send`
  - Mô tả: Gửi email bất đồng bộ (đưa vào hàng đợi)
  - Body: `SendEmailRequest`
    - `to: string[]` (bắt buộc)
    - `subject: string` (bắt buộc)
    - `body: string` (bắt buộc)
    - `is_html: boolean`
  - Response 200: `EmailResponse { success, message, queue_id }`
  - Lỗi: 400 (payload sai), 500 (dịch vụ email lỗi)

- POST `/api/v1/email/send-sync`
  - Mô tả: Gửi email đồng bộ (chờ SMTP trả kết quả)
  - Body: `SendEmailRequest`
  - Response 200: `EmailResponse { success, message }`
  - Lỗi: 400, 500

- POST `/api/v1/email/send-bulk`
  - Mô tả: Gửi nhiều email bất đồng bộ
  - Body: `SendBulkEmailRequest { emails: SendEmailRequest[] }`
  - Response 200: `EmailResponse { success, message }`
  - Lỗi: 400, 500

- GET `/api/v1/email/status`
  - Mô tả: Trạng thái dịch vụ email, kích thước hàng đợi
  - Response 200: `EmailStatusResponse { queue_size, status, message }`

- GET `/api/v1/email/test-connection`
  - Mô tả: Kiểm tra kết nối SMTP
  - Response 200: `TestEmailConnectionResponse { success, message, connected }`
  - Lỗi: 500 khi không kết nối được

Triển khai nổi bật:
- `EmailService` duy trì queue và worker pool để xử lý gửi email bất đồng bộ (`SendAsync`).
- Với gửi đồng bộ, handler gọi `SendSync` để SMTP trả kết quả ngay.
- Có hỗ trợ template qua `TemplateManager` nếu cần (các hàm `SendTemplate`, `SendTemplateSync`).

---

## 3. Notification API

Các endpoint được định nghĩa trong `NotificationHandler` và dùng Bearer token. Định tuyến có thể cần thêm vào router (xem phần [Khởi tạo & Wiring](#5-khởi-tạo--wiring-trong-dự-án)).

Các endpoint dự kiến:

- POST `/api/v1/notifications/send`
  - Mô tả: Gửi thông báo đến 1 user theo kênh chỉ định.
  - Body: `SendNotificationRequest`
    - `user_id: string(UUID)` (bắt buộc)
    - `type: "email" | "in_app"` (bắt buộc)
    - `title: string` (bắt buộc)
    - `message: string` (bắt buộc)
    - `data: object` (tuỳ kênh)
      - Với email channel: cần `data.email` (bắt buộc), tuỳ chọn `data.html_body`.
  - Response 200: `NotificationResponse { id, user_id, type, status, title, message, sent_at, created_at }`
  - Lỗi: 400 (UUID/đầu vào sai), 500 (kênh không tồn tại, dữ liệu kênh thiếu…)

- POST `/api/v1/notifications/send-bulk`
  - Mô tả: Gửi nhiều thông báo.
  - Body: `SendBulkNotificationRequest { notifications: SendNotificationRequest[] }`
  - Response 200: `SuccessResponse { message }`

- GET `/api/v1/notifications/user/{user_id}`
  - Mô tả: Lấy danh sách thông báo của user (kênh in-app).
  - Response 200: `ListNotificationsResponse { notifications[], total, unread }`

- POST `/api/v1/notifications/mark-read`
  - Mô tả: Đánh dấu đã đọc 1 thông báo (in-app).
  - Body: `MarkAsReadRequest { notification_id: string(UUID) }`
  - Response 200: `SuccessResponse { message }`
  - Lỗi: 400, 404 (không tồn tại)

- GET `/api/v1/notifications/status`
  - Mô tả: Trạng thái dịch vụ notification và các kênh sẵn có.
  - Response 200: `NotificationServiceStatusResponse { email_channel, in_app_channel, status, message }`

Triển khai nổi bật:
- `NotificationService` là Subject, ánh xạ `type` -> Observer tương ứng và gọi `observer.Update(...)`.
- `EmailChannel` gửi qua `EmailService` (bất đồng bộ), yêu cầu `notification.Data["email"]`.
- `InAppChannel` ghi thông báo vào DB qua `NotificationRepository`, hỗ trợ `GetNotifications` và `MarkAsRead`.

---

## 4. Observer Pattern: Khái niệm & Ứng dụng

Observer Pattern cho phép Subject phát sự kiện tới nhiều Observer mà không phụ thuộc cụ thể vào từng kênh. Trong dự án:

- Subject: `NotificationService`
  - `Attach(observer)`/`Detach(type)` để đăng ký/hủy kênh.
  - `Notify(ctx, notification)` gọi `Update` của Observer theo `notification.Type`.

- Observers: các kênh gửi thông báo triển khai `interfaces.Observer`:
  - `EmailChannel` (email)
  - `InAppChannel` (lưu DB, hiển thị trong app)
  - `SMSChannel` (mẫu minh hoạ – có thể mở rộng thật qua nhà cung cấp SMS)

Tại sao dùng cho email/notification:
- Tách biệt logic gửi theo kênh khỏi business logic: không cần `switch/case` hay `if-else` phức tạp trong handler/service.
- Dễ mở rộng: thêm kênh mới (Push, Webhook, Slack…) chỉ cần implement Observer và `Attach`.
- Tái sử dụng & kiểm thử tốt: từng kênh test độc lập; Subject test mapping/luồng gọi.
- Tương thích xử lý bất đồng bộ: kênh như email có queue/worker riêng không ảnh hưởng kênh khác.

---

## 5. Khởi tạo & Wiring trong dự án

Các thành phần chính cần khởi tạo khi chạy server:

- Email client + service:
  - `utils.EmailClient` cấu hình từ `configs/.env` (`EMAIL_SMTP_*`).
  - `services.EmailService` tạo queue/worker để gửi async.
  - Tuỳ chọn `TemplateManager` để gửi theo template.

- Notification service + channels:
  - `services.NotificationService` (Subject).
  - `observers.EmailChannel` dùng `EmailService` và địa chỉ `From`.
  - `observers.InAppChannel` dùng `NotificationRepository` để lưu DB.
  - Gắn kênh: `notificationService.Attach(emailChannel)`, `notificationService.Attach(inAppChannel)`.

- Định tuyến:
  - Email routes đã có trong `routes.setupEmailRoutes`.
  - Notification routes: cần thêm group `/api/v1/notifications` và map tới `NotificationHandler` (xem `internal/app/handlers/notification_handler.go`).

Lưu ý: Nếu bạn chưa thấy Notification routes trong `routes.go`, hãy bổ sung wiring khi cần phát hành API notification.

---

## 6. Xử lý lỗi & bảo mật

- Bảo mật: Tất cả endpoints Email/Notification yêu cầu Bearer token (JWT). Middleware: `AuthMiddleware.Authenticate()`.
- Xác thực dữ liệu: dùng `ShouldBindJSON` + tags `binding:"required"` trong models. Sai định dạng trả `400` với `models.ErrorResponse`.
- Lỗi kênh/ngoại lệ:
  - Gửi email thất bại: `500` với "Failed to send email".
  - Notification kênh không tìm thấy hoặc dữ liệu kênh thiếu (ví dụ thiếu `data.email` khi type=email) có thể trả `500`. Khuyến nghị: chuẩn hoá thông điệp lỗi client-facing và log chi tiết server-side.

---

## 7. Ví dụ gọi API

Email – gửi bất đồng bộ:

```bash
curl -X POST "http://localhost:8001/api/v1/email/send" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "Welcome",
    "body": "Thanks for joining!",
    "is_html": false
  }'
```

Notification – gửi qua email channel:

```bash
curl -X POST "http://localhost:8001/api/v1/notifications/send" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "email",
    "title": "Welcome!",
    "message": "Thank you for joining us",
    "data": {"email": "user@example.com", "html_body": "<p>Hello</p>"}
  }'
```

Notification – lấy danh sách in-app:

```bash
curl -X GET "http://localhost:8001/api/v1/notifications/user/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

## 8. Kiểm thử liên quan

- Unit tests có sẵn cho Email handler: `test/unit/email_handler_test.go` (mock service, kiểm tra 200/400/500, queue size, test connection).
- Gợi ý thêm tests cho Notification handler:
  - Happy paths: send (email/in_app), send-bulk, get by user, mark-read.
  - Edge cases: user_id không hợp lệ, thiếu `data.email` với channel email, kênh không đăng ký.

---

## 9. Gợi ý mở rộng

- Thêm kênh mới: Push/Webhook/Slack chỉ cần implement `interfaces.Observer` và `Attach` vào `NotificationService`.
- Bổ sung retry/Dead-letter queue cho EmailService khi gửi thất bại.
- Chuẩn hoá lỗi client-facing (error codes) cho Notification.
- Hoàn thiện wiring Notification routes trong `routes.go` và `cmd/server/main.go` khi phát hành.
- Log/metrics dành riêng cho từng kênh để theo dõi hiệu năng.

---

Chốt lại: Email API xử lý gửi thư với queue/worker cho hiệu năng; Notification API dựa trên Observer Pattern để dễ mở rộng kênh. Cách tiếp cận này giúp tách biệt, dễ kiểm thử, và phù hợp hệ thống cần nhiều kênh thông báo.
