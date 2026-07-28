# SLICE 3 — Quota, ownership, CORS (+ optionalAuth)

**Truy ngược:** vision (ẩn danh 10 link/ngày, user unlimited, FE gọi API kèm cookie) · DESIGN ADR-013, ADR-014 · gộp T2.7 (optionalAuth) từ slice 2.
**Nhánh:** `feature/frontend-integration`. Mỗi ticket 1 commit, dừng review.

## Mục tiêu slice
Xong slice này: request nhận diện được user (nếu có token); `encode` áp quota ẩn danh 10/ngày theo IP và gắn ownership khi đăng nhập; `GET /api/v1/links` liệt kê link của user; CORS cho FE gọi kèm cookie.

## Kiến trúc (ADR-013/014)
- `links` += `user_id` (nullable FK), `creator_ip` (text nullable).
- Đăng nhập → set `user_id`, bỏ qua quota. Ẩn danh → set `creator_ip`, đếm trong ngày (UTC) ≥10 → **429**.
- IP = `RemoteAddr` (không tin XFF). CORS: go-chi/cors + `ALLOWED_ORIGINS` cụ thể + credentials.

## Data model (migration mới)
```sql
-- 000004_add_link_ownership
ALTER TABLE links ADD COLUMN user_id    UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE links ADD COLUMN creator_ip TEXT;
CREATE INDEX idx_links_user_id    ON links(user_id);
CREATE INDEX idx_links_creator_ip ON links(creator_ip, created_at);
```
*(user_id `ON DELETE SET NULL`: xoá user thì link ẩn danh hoá, không mất. Index creator_ip+created_at phục vụ đếm quota.)*

## Endpoints
| Method | Path | Auth | Ghi chú |
|---|---|---|---|
| POST | `/api/v1/encode` | optional | user→ownership+unlimited; ẩn danh→quota 10/ngày theo IP (429 nếu vượt) |
| GET | `/api/v1/links` | **required** | list link của user hiện tại |

## Config thêm
`ALLOWED_ORIGINS` (CSV, vd `http://localhost:5173`).

## Deps mới
`github.com/go-chi/cors`.

---

## TICKET-3.1 — Schema + sqlc cho ownership & quota
**Phụ thuộc:** —
### Chỉ dẫn
- Migration `000004_add_link_ownership.{up,down}.sql` (ALTER ADD COLUMN + index; down = DROP COLUMN).
- `data/queries/links.sql`: cập nhật `CreateLink` (thêm `user_id`, `creator_ip` vào INSERT); thêm `CountLinksByCreatorIPSince :one` (WHERE creator_ip=$1 AND created_at >= $2), `ListLinksByUserID :many` (WHERE user_id=$1 ORDER BY created_at DESC).
- `make sqlc` → commit. `CreateLinkParams` giờ có UserID (uuid nullable → `*uuid.UUID`? sqlc: nullable uuid → override có sẵn? kiểm: uuid nullable → `uuid.NullUUID` hoặc `*uuid.UUID`. Nếu ra NullUUID thì repo map).
### Acceptance
- [ ] migrate up/down sạch; sqlc build; go build xanh.

## TICKET-3.2 — optionalAuth + requireAuth + user context
**Phụ thuộc:** slice 2 (pkg/auth ParseAccessToken)
### Chỉ dẫn
- `internal/pkg/authctx`: `WithUser(ctx, userID uuid.UUID) ctx` + `UserFromContext(ctx) (uuid.UUID, bool)` (key private).
- Middleware `handler`/hoặc `httpserver`: `optionalAuth(secret []byte)` — parse `Authorization: Bearer` → `pkgauth.ParseAccessToken` → `WithUser`; thiếu/lỗi → next không chặn. `requireAuth` — 401 nếu `!UserFromContext`.
- Đặt ở đâu: middleware cần `secret` → nhận qua tham số khi dựng router (main truyền `cfg.JWTSecret`). optionalAuth áp cho group `/api/v1`; requireAuth chỉ bọc `/links`.
### Acceptance
- [ ] Unit test: token hợp lệ→context có user; thiếu/hỏng→qua mà không user; requireAuth thiếu user→401.

## TICKET-3.3 — Repo link: ownership + count + list
**Phụ thuộc:** 3.1
### Chỉ dẫn
- `model.Link` += `UserID *uuid.UUID`, `CreatorIP *string` (nullable).
- `repository/link`: `Create` set user_id/creator_ip; `CountByCreatorIPSince(ctx, ip string, since time.Time) (int, error)`; `ListByUserID(ctx, userID uuid.UUID) ([]model.Link, error)`. Cập nhật interface + mock (mockery).
- DB test: create có/không user_id; count theo ip+ngày; list theo user (thứ tự created_at desc).
### Acceptance
- [ ] `make test` xanh cho repo link.

## TICKET-3.4 — Controller link: quota + ownership + list
**Phụ thuộc:** 3.3
### Chỉ dẫn
- `apperror` += `TooManyRequests(429)`.
- `Encode` đổi chữ ký: nhận `userID *uuid.UUID` + `clientIP string`. Logic: nếu userID != nil → tạo link với user_id (bỏ qua quota). Nếu nil → `CountByCreatorIPSince(ip, startOfTodayUTC)` ≥ `MaxAnonLinksPerDay` (=10, hằng model) → 429; else tạo với creator_ip. Dedup vẫn trước (link đã tồn tại → trả cũ, KHÔNG tính quota).
- Thêm use-case `ListByUser(ctx, userID) ([]model.Link, error)`.
- Mock test: user→no quota+ownership; ẩn danh dưới hạn→tạo; ẩn danh chạm hạn→429; dedup không tính quota.
### Acceptance
- [ ] `go test controller/link` xanh (mock).

## TICKET-3.5 — Handler: encode ctx-user+IP, GET /links, gắn middleware
**Phụ thuộc:** 3.2, 3.4
### Chỉ dẫn
- `encode` handler: lấy `userID` từ `authctx.UserFromContext`, `clientIP` từ `RemoteAddr` (helper `clientIP(r)` tách host khỏi `host:port`). Truyền vào `ctrl.Encode`.
- Handler mới `ListLinks` (GET /api/v1/links): `requireAuth` → userID từ context → `ctrl.ListByUser` → JSON `[]{code, short_url, long_url, created_at}`.
- Router: `optionalAuth` áp cho group `/api/v1`; `/links` bọc `requireAuth`. main truyền secret.
- Test: encode với/không user (mock controller, giả context); 429 propagate; list requireAuth (thiếu→401, có→list).
### Acceptance
- [ ] `go test handler/...` xanh; swagger cập nhật (encode note quota, /links).

## TICKET-3.6 — CORS + config + wiring
**Phụ thuộc:** 3.5
### Chỉ dẫn
- `config` += `AllowedOrigins []string` (parse CSV `ALLOWED_ORIGINS`, required). Test + compose env.
- `httpserver.New`: thêm `cors.Handler` (go-chi/cors) với AllowedOrigins, AllowCredentials true, methods/headers, trước routes. main truyền origins.
- `.env`/compose: `ALLOWED_ORIGINS=http://localhost:5173`.
### Acceptance
- [ ] preflight OPTIONS trả đúng header (test); `make up` FE gọi được (Kyle verify).

---

## Thứ tự: 3.1 → 3.2 → 3.3 → 3.4 → 3.5 → 3.6.

## Ngoài scope slice 3 (tương lai)
Xoá/sửa link của user (management UI), rolling-window quota, tin XFF sau trusted proxy, SameSite=None+CSRF khi FE khác-site.
