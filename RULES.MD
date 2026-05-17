# Contributing to go-echo-boilerplate

This guide defines the canonical coding patterns, architecture conventions, and workflow for this Go Echo boilerplate. It serves as a reference for both team members and AI-assisted development.

**Before writing code, look at how the codebase already does it.** The best contribution mirrors existing patterns rather than inventing new ones.

**Stack**: Go 1.2x | Echo v4 | PostgreSQL (pgx/v5) | Squirrel | Viper | Redis (go-redis/v9) | Docker

---

## Table of Contents

- [Philosophy](#philosophy)
- [Architecture Overview](#architecture-overview)
- [Code Style Guide](#code-style-guide)
  - [Enums](#enums)
    - [Roles](#roles)
    - [Log Levels & Log Types](#log-levels--log-types)
    - [Error Types](#error-types)
  - [Domain Entities](#domain-entities)
  - [DAO Layer](#dao-layer)
  - [Repository Layer](#repository-layer)
  - [Service Layer](#service-layer)
  - [DTOs](#dtos)
    - [Request DTOs](#request-dtos)
    - [Response DTOs](#response-dtos)
    - [Pagination DTOs](#pagination-dtos)
  - [Handlers](#handlers)
  - [Routes](#routes)
  - [Error Handling](#error-handling)
  - [Logging](#logging)
  - [Permissions & Group Permissions](#permissions--group-permissions)
  - [Validation](#validation)
  - [Request & Response Standards](#request--response-standards)
- [Docker](#docker)
- [Common Pitfalls](#common-pitfalls)
- [Git Workflow](#git-workflow)
- [Database Migrations](#database-migrations)

---

## Philosophy

1. **Relentless simplicity** — Don't add abstractions, interfaces, or indirection unless you genuinely need them. If it's used once, inline it. Three similar lines are better than a premature abstraction.

2. **Follow the patterns** — Handlers are thin HTTP adapters. Services hold business logic. Repositories handle data access. DAOs map database rows. Don't blur these lines.

3. **Idiomatic Go** — Use standard library conventions, explicit error handling, struct embedding, and composition over inheritance. Don't fight the language.

4. **Code reads like a book** — Don't be clever — be clear. A reader should understand intent and business logic at a glance.

5. **Docker always** — Every service, database, and tool runs in Docker. There is no "run locally without Docker" path. If it isn't in `docker-compose.yml`, it doesn't exist.

---

## Architecture Overview

### Data Flow

```
Request:  JSON → RequestDTO → dto.ToDomain() → Service(domain) → Repo(domain) → DAO → SQL
Response: SQL → DAO → dao.ToDomain() → Repo(domain) → Service(domain) → dto.FromDomain() → ResponseDTO → JSON
```

### Hard Rules

- **Domain entities have zero struct tags** — no `json`, no `db`, no `validate`. JSONB serialization is handled by DAOs.
- **DTOs own JSON serialization** (`json` tags + `validate` tags). All API input/output goes through DTOs.
- **DAOs own database scanning**. All DB row mapping goes through DAOs.
- **Services only accept and return domain types.** Never DTOs, never DAOs.
- **Handlers never pass domain entities directly to responses** — always convert via DTO mappers.
- **JSON field names use `snake_case`** in all request and response DTOs.
- **Every response — success or error — goes through the shared response helpers.** No raw `c.JSON()` calls in handlers.
- **Every request and every response is logged.** The logger middleware handles this automatically; do not add manual request/response logs in handlers.

### Why Every Module Uses DAO → Domain → DTO

The `DAO → Domain → DTO` conversion chain may feel like overhead for simple CRUD modules where the structs look nearly identical. **Keep it consistent anyway.** Each conversion is a struct field copy — nanoseconds compared to the milliseconds spent on database queries, network I/O, and JSON marshaling. The real costs of skipping it are human, not computational:

- If one module skips the DAO, contributors won't know whether it was a deliberate shortcut or a mistake.
- Services grow business logic over time. Retrofitting the layer boundary later is harder than having it from the start.
- Multiple consumers (HTTP handlers, WebSocket handlers, background workers) can reuse the same service when it speaks domain types.

If you ever hit a measured, genuine hot path (batch imports, bulk exports), optimize that specific method. But that should be a profiler-driven decision, not a default.

### Project Layout

```
cmd/
  api/main.go                          # Entry point
  seed/main.go                         # Seeds permissions + group_permissions + admin role
internal/
  config/                              # Viper config loading
  domain/                              # Entities (zero tags)
  auth/                                # Authentication module (special layout)
    application/                       # Services (JWT, sessions, RBAC checks)
    domain/                            # Auth provider interfaces
    infra/
      http/                            # Handlers + routes + dto/
      persistence/                     # Repo implementations
  modules/                             # Feature modules
    {module}/
      application/                     # Services + repository interfaces
      infra/
        http/                          # Handlers, routes, dto/
        persistence/                   # Repo implementations + dao/
  platform/                            # database/, cache/, storage/
  server/                              # Server bootstrap + route wiring
  shared/
    app_errors/                        # Typed AppError + error codes
    enums/                             # roles.go, log_level.go, log_type.go, error_type.go
    middleware/                        # Auth, logger, permission guard
    response/                          # OK, Created, Paginated, Error helpers
    constants/
pkg/
  logger/                              # Logger interface + console + DB implementations
  req_ctx/                             # BindAndValidate, ParsePagination, context helpers
  query/                               # Filter/sort helpers for Squirrel
configs/                               # YAML config files
docs/                                  # Swagger/OpenAPI generated docs
migrations/                            # SQL migration files (golang-migrate)
docker-compose.yml
Dockerfile
Makefile
```

### Module Isolation Rules

- Each module owns its own repository, service, handler, DTOs, DAOs, and routes.
- Modules **never** import each other. Cross-module communication goes through domain interfaces.
- Shared code lives in `internal/shared/` or `pkg/`.

---

## Code Style Guide

Every code example below is the **canonical pattern**. Follow these exactly when adding new modules or modifying existing ones.

---

### Enums

**Location**: `internal/shared/enums/`

All project-wide enums live here. Never define a role, log level, log type, or error type string inline in business logic — always import from this package.

---

#### Roles

**Location**: `internal/shared/enums/roles.go`

Roles represent the coarse-grained identity of a user. Fine-grained access is handled by permissions — roles are for convenience grouping only (e.g., seeding an admin group, UI display).

```go
package enums

type Role string

const (
    RoleAdmin     Role = "admin"      // Full system access — assigned during seed
    RoleModerator Role = "moderator"  // Can manage content, cannot manage users or permissions
    RoleUser      Role = "user"       // Default role for registered users
    RoleGuest     Role = "guest"      // Unauthenticated or anonymous access
)

func (r Role) String() string { return string(r) }

func (r Role) IsValid() bool {
    switch r {
    case RoleAdmin, RoleModerator, RoleUser, RoleGuest:
        return true
    }
    return false
}
```

**Rules**:
- A user has one role. Fine-grained access comes from the permission system, not more roles.
- `RoleAdmin` is assigned only by the seed. There is no "promote to admin" API endpoint.
- `RoleGuest` is the implicit role for any request with no valid JWT — middleware sets this automatically.
- When adding a new role: add it here, update `IsValid()`, add a migration to create its default `group_permission`, and update the seed.

---

#### Log Levels & Log Types

**Location**: `internal/shared/enums/log_level.go` and `internal/shared/enums/log_type.go`

```go
// log_level.go
package enums

type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"  // Query tracing, internal state — never in production
    LogLevelInfo  LogLevel = "info"   // Successful business events (user created, payment processed)
    LogLevelWarn  LogLevel = "warn"   // Non-fatal failures (retry triggered, deprecated path hit)
    LogLevelError LogLevel = "error"  // Unexpected failures — always include the original error
)

func (l LogLevel) String() string { return string(l) }
```

```go
// log_type.go
package enums

type LogType string

const (
    LogTypeRequest  LogType = "request"   // Incoming HTTP request — logged by middleware
    LogTypeResponse LogType = "response"  // Outgoing HTTP response — logged by middleware
    LogTypeSystem   LogType = "system"    // App lifecycle: startup, shutdown, migrations
    LogTypeAudit    LogType = "audit"     // Security-sensitive actions: login, permission change
)

func (t LogType) String() string { return string(t) }
```

---

#### Error Types

**Location**: `internal/shared/enums/error_type.go`

```go
package enums

type ErrorType string

const (
    ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"    // 400 — bad input format
    ErrorTypeAuth         ErrorType = "AUTH_ERROR"           // 401 — missing or expired token
    ErrorTypeForbidden    ErrorType = "FORBIDDEN_ERROR"      // 403 — valid token, no permission
    ErrorTypeNotFound     ErrorType = "NOT_FOUND_ERROR"      // 404 — resource does not exist
    ErrorTypeConflict     ErrorType = "CONFLICT_ERROR"       // 409 — duplicate or bad state transition
    ErrorTypeRateLimit    ErrorType = "RATE_LIMIT_ERROR"     // 429 — too many requests
    ErrorTypeDatabase     ErrorType = "DATABASE_ERROR"       // 500 — DB failure (hidden from client)
    ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"       // 500 — unexpected server error
)

func (e ErrorType) String() string { return string(e) }
```

---

### Domain Entities

**Location**: `internal/domain/{entity}.go`

Domain entities have **zero struct tags**. They contain the entity struct, type aliases for IDs, factory functions, and state-machine methods.

```go
package domain

import (
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/boilerplate/internal/shared/enums"
)

type UserID uuid.UUID
type UserStatus string

const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
    UserStatusBanned   UserStatus = "banned"
)

type User struct {
    ID           UserID
    Name         string
    Email        string
    Phone        string
    Role         enums.Role
    PasswordHash string
    Status       UserStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func NewUser(name, email, phone string) *User {
    return &User{
        ID:     UserID(uuid.New()),
        Name:   name,
        Email:  email,
        Phone:  phone,
        Role:   enums.RoleUser,
        Status: UserStatusActive,
    }
}

func (u UserID) String() string       { return uuid.UUID(u).String() }
func ParseUserID(s string) (UserID, error) {
    id, err := uuid.Parse(s)
    if err != nil {
        return UserID{}, err
    }
    return UserID(id), nil
}
```

**State machine methods** validate transitions:

```go
func (u *User) Deactivate() error {
    if u.Status == UserStatusBanned {
        return app_errors.New(app_errors.ErrCodeConflict, "banned users cannot be deactivated")
    }
    u.Status = UserStatusInactive
    return nil
}
```

**Key conventions**:
- Type aliases for IDs (`UserID`, `ProductID`) with `String()` and `Parse` methods.
- Nullable fields use pointer types only when `nil` has business meaning.
- Repository interfaces live in the application layer (see [Service Layer](#service-layer)).

---

### DAO Layer

**Location**: `internal/modules/{module}/infra/persistence/dao/{module}_dao.go`

DAOs are the **only structs that touch database rows**. Each DAO provides `ToDomain()` (method) and `FromDomain()` (package-level function):

```go
package dao

import (
    "time"

    "github.com/google/uuid"

    "github.com/yourorg/boilerplate/internal/domain"
    "github.com/yourorg/boilerplate/internal/shared/enums"
)

type UserDAO struct {
    ID           uuid.UUID
    Name         string
    Email        string
    Phone        string
    Role         string    // Raw DB string, converted in ToDomain()
    PasswordHash string
    Status       string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (d *UserDAO) ToDomain() *domain.User {
    return &domain.User{
        ID:           domain.UserID(d.ID),
        Name:         d.Name,
        Email:        d.Email,
        Phone:        d.Phone,
        Role:         enums.Role(d.Role),
        PasswordHash: d.PasswordHash,
        Status:       domain.UserStatus(d.Status),
        CreatedAt:    d.CreatedAt,
        UpdatedAt:    d.UpdatedAt,
    }
}

func FromDomain(u *domain.User) UserDAO {
    return UserDAO{
        ID:           uuid.UUID(u.ID),
        Name:         u.Name,
        Email:        u.Email,
        Phone:        u.Phone,
        Role:         u.Role.String(),
        PasswordHash: u.PasswordHash,
        Status:       string(u.Status),
        CreatedAt:    u.CreatedAt,
        UpdatedAt:    u.UpdatedAt,
    }
}
```

**JSONB columns** — The DAO owns `json` tags for JSONB marshaling. The domain entity stays tag-free:

```go
// Domain — no tags
type UserPreferences struct {
    Language string
    Timezone string
}

// DAO — owns json tags for JSONB marshaling
type UserPreferencesDAO struct {
    Language string `json:"language"`
    Timezone string `json:"timezone"`
}
```

---

### Repository Layer

**Location**: `internal/modules/{module}/infra/persistence/{module}_repo_impl.go`

**Always use Squirrel for query building. No raw SQL strings.** Repositories scan into DAOs, then convert to domain entities.

```go
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *UserRepoImpl) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
    query, args, err := psql.Select(
        "id", "name", "email", "phone", "role", "password_hash", "status", "created_at", "updated_at",
    ).From("users").Where(sq.Eq{"id": uuid.UUID(id)}).ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    var d dao.UserDAO
    err = r.db.QueryRow(ctx, query, args...).Scan(
        &d.ID, &d.Name, &d.Email, &d.Phone, &d.Role, &d.PasswordHash, &d.Status, &d.CreatedAt, &d.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, app_errors.NotFound("user")
        }
        return nil, app_errors.DatabaseFailure(err)
    }

    return d.ToDomain(), nil
}
```

**Repository rules**:
- Scan into `dao.XxxDAO`, then call `d.ToDomain()`. Use `dao.FromDomain()` for inserts/updates.
- Map `pgx.ErrNoRows` to `app_errors.NotFound("resource_name")`.
- `defer rows.Close()` after the error check, and always check `rows.Err()`.
- Use transactions (`db.Begin` + `defer tx.Rollback(ctx)`) for multi-statement operations.
- Every method takes `context.Context` as the first parameter.
- Wrap errors: `fmt.Errorf("failed to {verb} {noun}: %w", err)`.

---

### Service Layer

**Location**: `internal/modules/{module}/application/{module}_service.go`

Services contain business logic. They **only accept and return domain types** — never import DTO or DAO packages. The repository interface lives in the same package:

```go
package application

type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id domain.UserID) error
}

type UserService struct {
    repository UserRepository
    logger     logger.Logger
}

func (s *UserService) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
    if limit <= 0  { limit = 20 }
    if limit > 100 { limit = 100 }
    return s.repository.FindAll(ctx, limit, offset)
}
```

**Service rules**:
- Enforce business constraints (pagination limits, access control, state validation).
- Fetch-before-update: always `FindByID` before `Update`.
- Log meaningful actions with structured key-value pairs.

---

### DTOs

**Location**: `internal/modules/{module}/infra/http/dto/`

DTOs are the **only structs with `json` and `validate` tags**. JSON field names use **`snake_case`**.

---

#### Request DTOs

```go
// Create — non-pointer fields, validate:"required" on mandatory ones
type CreateUserReq struct {
    Name  string `json:"name"  validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
    Phone string `json:"phone" validate:"required,e164"`
}

func (r *CreateUserReq) ToDomain() *domain.User {
    return domain.NewUser(r.Name, r.Email, r.Phone)
}

// Update — pointer fields for partial update
type UpdateUserReq struct {
    Name  *string `json:"name,omitempty"  validate:"omitempty,min=2,max=100"`
    Email *string `json:"email,omitempty" validate:"omitempty,email"`
    Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
}

func (r *UpdateUserReq) ToDomain(existing *domain.User) {
    if r.Name  != nil { existing.Name  = *r.Name  }
    if r.Email != nil { existing.Email = *r.Email }
    if r.Phone != nil { existing.Phone = *r.Phone }
}
```

---

#### Response DTOs

Every module exposes a single-item response function and a list response function. **Always use these — never build a response map inline in a handler.**

```go
type UserRes struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Role      string    `json:"role"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Single item
func UserResFromDomain(u *domain.User) UserRes {
    return UserRes{
        ID:        u.ID.String(),
        Name:      u.Name,
        Email:     u.Email,
        Phone:     u.Phone,
        Role:      u.Role.String(),
        Status:    string(u.Status),
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}

// List — always returns []UserRes, never nil
func UserListResFromDomain(users []*domain.User) []UserRes {
    res := make([]UserRes, len(users))
    for i := range users {
        res[i] = UserResFromDomain(users[i])
    }
    return res
}
```

---

#### Pagination DTOs

**Location**: `pkg/req_ctx/pagination.go`

Pagination request parsing and normalization is done once in `req_ctx` — never manually in a handler.

```go
// PaginationReq is populated by reqctx.ParsePagination(c).
// All normalization (defaults, clamping) happens here — handlers receive clean values.
type PaginationReq struct {
    Page    int `query:"page"`
    PerPage int `query:"per_page"`
}

func ParsePagination(c echo.Context) PaginationReq {
    var p PaginationReq
    _ = (&echo.DefaultBinder{}).BindQueryParams(c, &p)
    if p.Page    <= 0   { p.Page    = 1   }
    if p.PerPage <= 0   { p.PerPage = 20  }
    if p.PerPage > 100  { p.PerPage = 100 }
    return p
}

func (p PaginationReq) Offset() int { return (p.Page - 1) * p.PerPage }
```

The paginated response shape is built by `response.Paginated()` — handlers never build it manually:

```json
{
    "success": true,
    "data": [ ... ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 143,
        "total_pages": 8
    }
}
```

Usage in a handler:

```go
func (h *UserHandler) ListUsers(c echo.Context) error {
    p := reqctx.ParsePagination(c)
    users, total, err := h.service.FindAll(c.Request().Context(), p.PerPage, p.Offset())
    if err != nil {
        return err
    }
    return response.Paginated(c, dto.UserListResFromDomain(users), p.Page, p.PerPage, total)
}
```

---

### Handlers

**Location**: `internal/modules/{module}/infra/http/{module}_handler.go`

Handlers are thin — bind input, call service, convert result via DTO. No business logic. No logging (the middleware handles it).

```go
// @Summary     Create a user
// @Tags        users
// @Security    BearerAuth
// @Param       request body dto.CreateUserReq true "User data"
// @Success     201 {object} object{success=bool,data=dto.UserRes}
// @Router      /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req dto.CreateUserReq
    if err := reqctx.BindAndValidate(c, &req); err != nil {
        return err
    }
    user := req.ToDomain()
    if err := h.service.CreateUser(c.Request().Context(), user); err != nil {
        return err
    }
    return response.Created(c, dto.UserResFromDomain(user))
}

// @Summary     Get a user by ID
// @Tags        users
// @Security    BearerAuth
// @Param       id path string true "User ID"
// @Success     200 {object} object{success=bool,data=dto.UserRes}
// @Router      /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
    userID, err := domain.ParseUserID(c.Param("id"))
    if err != nil {
        return app_errors.InvalidInput()
    }
    user, err := h.service.FindByID(c.Request().Context(), userID)
    if err != nil {
        return err
    }
    return response.OK(c, dto.UserResFromDomain(user))
}
```

**Key conventions**:
- Use `reqctx.BindAndValidate(c, &req)` for body parsing + validation.
- Use `c.Request().Context()` — **never** `context.Background()`.
- Return errors directly — the error middleware converts them to the canonical error response shape.
- Include Swagger annotations on every handler.
- **Never call `c.JSON()` directly.** Always use `response.OK`, `response.Created`, `response.Paginated`.

---

### Routes

**Location**: `internal/modules/{module}/infra/http/routes.go`

Each module exposes a `RegisterRoutes` function that wires dependencies bottom-up (repo → service → handler) and registers endpoints.

```go
func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, log logger.Logger) {
    repo    := persistence.NewUserRepoImpl(db)
    service := application.NewUserService(repo, log)
    handler := NewUserHandler(service)

    g := e.Group("/users", middleware.Auth())
    g.GET("",     handler.ListUsers)
    g.POST("",    handler.CreateUser)
    g.GET("/:id", handler.GetUser)
    g.PUT("/:id", handler.UpdateUser)
    g.DELETE("/:id", handler.DeleteUser)
}
```

REST conventions: `GET /`, `GET /:id`, `POST /`, `PUT /:id`, `DELETE /:id`.

---

### Error Handling

Errors flow up untouched until the Echo error middleware catches and formats them. No handler ever inspects or wraps an error from a service — it just returns it.

```
Repository → app_errors.NotFound(), app_errors.DatabaseFailure(), or fmt.Errorf()
Service    → app_errors.New() or fmt.Errorf("failed to {action}: %w", err)
Handler    → return err   ← that's it
Middleware → catches everything, maps to canonical error response
```

#### Canonical error response shape

Every error response — whether it comes from a DB failure, a missing resource, a bad JWT, a validation failure, or a panic — produces exactly this JSON:

```json
{
    "type":      "NOT_FOUND_ERROR",
    "status":    404,
    "message":   "user not found",
    "path":      "/users/abc-123",
    "method":    "GET",
    "timestamp": "2024-11-03T14:22:05+05:00"
}
```

| Field | Description |
|---|---|
| `type` | One of the `ErrorType` enum values — used by clients to branch on error kind |
| `status` | HTTP status code, also set on the response |
| `message` | Human-readable, user-facing message — **always present, never empty** |
| `path` | Request path, taken from `c.Request().URL.Path` |
| `method` | HTTP method (`GET`, `POST`, etc.) |
| `timestamp` | RFC3339 with timezone offset — human-readable, sortable |

The error middleware (`internal/shared/middleware/error_middleware.go`) is the **only place** that produces this envelope. Never build it manually anywhere else.

#### AppError

```go
// Using named constructors — prefer these
app_errors.NotFound("user")                                           // type=NOT_FOUND_ERROR, 404, "user not found"
app_errors.InvalidInput()                                             // type=VALIDATION_ERROR, 400
app_errors.DatabaseFailure(err)                                       // type=DATABASE_ERROR, 500, cause hidden from client
app_errors.Forbidden("you do not have permission to perform this action")

// Using the generic constructor for custom cases
app_errors.New(app_errors.ErrCodeConflict, "email is already registered").WithCause(err)
```

#### DB error handling

The error middleware also catches raw `pgx` / database errors that bubble up without being wrapped in an `AppError`. It maps them automatically:

| DB error code | ErrorType | HTTP Status | Message shown to client |
|---|---|---|---|
| `23505` unique violation | `CONFLICT_ERROR` | 409 | Detail from `pg_exception_detail` |
| `23503` foreign key | `VALIDATION_ERROR` | 400 | Generic "related resource not found" |
| `23502` not null | `VALIDATION_ERROR` | 400 | Generic "missing required field" |
| `23514` check constraint | `VALIDATION_ERROR` | 400 | Generic "value out of allowed range" |
| Any other DB error | `DATABASE_ERROR` | 500 | "an internal error occurred" (detail hidden) |

The original DB error is always attached as `.WithCause(err)` for internal logging — it is **never** sent to the client.

#### Error message style

Messages in `AppError` are **user-facing**. Keep them:
- **Lowercase** — `"user not found"`, not `"User Not Found"`
- **No periods** — `"email is already registered"`, not `"Email is already registered."`
- **Descriptive** — `"email is already registered"`, not `"UNIQUE constraint violation"`
- **No internals** — SQL errors, file paths, stack traces go in `.WithCause(err)`, never in the message

#### Error wrapping convention

Always use the pattern `"failed to {verb} {noun}: %w"`:

```go
return fmt.Errorf("failed to create user: %w", err)    // Good
return err                                               // Bad — no context
return fmt.Errorf("something went wrong: %w", err)      // Bad — vague
```

---

### Logging

**Location**: `pkg/logger/`

The project uses a custom `Logger` interface with four levels. Every level has both a console output and an optional DB-write path. The middleware handles request and response logging automatically — do not add logging in handlers.

#### Logger interface

```go
type Logger interface {
    Debug(msg string, keysAndValues ...any)
    Info(msg string, keysAndValues ...any)
    Warn(msg string, keysAndValues ...any)
    Error(msg string, keysAndValues ...any)

    // DB-writing variants — use these for audit trails and security events
    DBDebug(ctx context.Context, msg string, keysAndValues ...any)
    DBInfo(ctx context.Context, msg string, keysAndValues ...any)
    DBWarn(ctx context.Context, msg string, keysAndValues ...any)
    DBError(ctx context.Context, msg string, keysAndValues ...any)
}
```

The four levels and when to use each:

| Level | Use for | DB write by default? |
|---|---|---|
| `Debug` | Query tracing, internal state dumps — disable in production | No |
| `Info` | Successful business events (user created, payment completed) | Optional |
| `Warn` | Non-fatal failures (retry triggered, deprecated API hit) | Optional |
| `Error` | Unexpected failures — always attach the original error | Yes |

#### Log entry stored in DB

Every DB-write log call stores one row in the `system_logs` table:

```json
{
    "id":          "uuid",
    "level":       "error",
    "type":        "system",
    "message":     "failed to send OTP",
    "context":     { "user_id": "abc", "phone": "+99361234567" },
    "created_at":  "2024-11-03T14:22:05+05:00"
}
```

#### Request / response logging

The logger middleware (`internal/shared/middleware/logger_middleware.go`) automatically logs every incoming request and every outgoing response. You never add this manually in handlers.

Request log (type = `request`, level = `info`):
```json
{
    "type":      "request",
    "method":    "POST",
    "path":      "/users",
    "ip":        "192.168.1.1",
    "user_agent": "Mozilla/5.0 ...",
    "timestamp": "2024-11-03T14:22:05+05:00"
}
```

Response log (type = `response`, level = `info` or `error` on 5xx):
```json
{
    "type":      "response",
    "method":    "POST",
    "path":      "/users",
    "status":    201,
    "latency_ms": 34,
    "timestamp": "2024-11-03T14:22:05+05:00"
}
```

#### Usage in services

```go
// Console only — for day-to-day business events
s.logger.Info("user created", "user_id", user.ID)
s.logger.Warn("rate limit approaching", "user_id", user.ID, "count", count)
s.logger.Error("failed to save user", "error", err)
s.logger.Debug("FindByID query", "query", query, "args", args)

// DB write — for audit trail or security-sensitive events
s.logger.DBInfo(ctx, "user role changed", "user_id", user.ID, "new_role", role)
s.logger.DBWarn(ctx, "failed login attempt", "email", email, "ip", ip)
s.logger.DBError(ctx, "payment processing failed", "order_id", id, "error", err)
```

**Rules**:
- Log in services (business events) and repositories (`Debug` only). Not in handlers.
- Use structured pairs: `"user_id", user.ID` — not `fmt.Sprintf("id=%s", user.ID)`.
- **Never log passwords, tokens, OTP codes, or raw PII.** Log entity IDs only.
- Use `DBxxx` variants for: logins, permission changes, role changes, any action with a `LogTypeAudit` semantic.

---

### Permissions & Group Permissions

**Location**: `internal/modules/permission/` and `internal/modules/group_permission/`

This project uses a two-layer RBAC system: **permissions** (atomic capabilities) and **group permissions** (named bundles of permission IDs assigned to users or roles).

#### Permission shape

A permission represents one atomic action in the system. Permissions are defined in code and seeded into the database on startup — they are never created via an API.

```go
// Domain entity — internal/domain/permission.go
type Permission struct {
    ID          int
    Module      string   // e.g. "users", "orders", "reports"
    Action      string   // e.g. "get all users", "delete order"
    Name        string   // e.g. "users.get_all", "orders.delete" — used in guards
    Description string
    CreatedAt   time.Time
}
```

The seed defines permissions like this:

```go
var Permissions = []domain.Permission{
    {Module: "users",   Action: "get all users",     Name: "users.get_all",    Description: "List all users with pagination"},
    {Module: "users",   Action: "get user by id",    Name: "users.get_one",    Description: "Fetch a single user by ID"},
    {Module: "users",   Action: "create user",       Name: "users.create",     Description: "Create a new user"},
    {Module: "users",   Action: "update user",       Name: "users.update",     Description: "Update user fields"},
    {Module: "users",   Action: "delete user",       Name: "users.delete",     Description: "Delete a user permanently"},
    {Module: "orders",  Action: "get all orders",    Name: "orders.get_all",   Description: "List all orders with pagination"},
    // ... one entry per protected endpoint
}
```

**Permission naming rules**:
- `Module` is the plural noun for the resource group: `"users"`, `"orders"`, `"reports"`.
- `Action` is a human-readable sentence: `"get all users"`, `"delete order"`.
- `Name` is the machine key used in permission guards: `"{module}.{verb}"` — lowercase, underscored, no spaces.

#### Group permissions

A group permission is a named bundle of permission IDs. It is assigned to a user (or to a role during seed). A user's effective permissions are the union of all permissions in their group.

```go
// Domain entity — internal/domain/group_permission.go
type GroupPermission struct {
    ID            int
    Name          string   // e.g. "admin", "moderator", "read_only"
    PermissionIDs []int    // The permission IDs this group grants
    CreatedAt     time.Time
}
```

The seed creates the default groups:

```go
var GroupPermissions = []GroupPermissionSeed{
    {
        Name: "admin",
        // Admin group gets every permission — populated dynamically from the permissions list
        Permissions: AllPermissionNames,
    },
    {
        Name: "moderator",
        Permissions: []string{"users.get_all", "users.get_one", "orders.get_all"},
    },
    {
        Name: "read_only",
        Permissions: []string{"users.get_all", "users.get_one"},
    },
}
```

#### How permission checking works

A permission guard middleware checks the authenticated user's group permissions before the handler runs:

```go
// Applied per-route in routes.go
g.DELETE("/:id", handler.DeleteUser, middleware.RequirePermission("users.delete"))
```

The middleware resolves: JWT → `user_id` → `group_permission_id` → `permission_ids` → check `"users.delete"` is in the set. If not, it returns a `403 FORBIDDEN_ERROR` immediately.

**Rules**:
- Never check permissions inside a handler or service. Use the middleware.
- New endpoints must have a corresponding seed entry in the permissions list.
- Never hard-code permission names as raw strings in more than one place — define them as constants in the seed file and reference them from routes.
- Changing a permission `Name` requires a migration — the name is the stable identifier.

---

### Validation

All request validation goes through `reqctx.BindAndValidate(c, &req)` which does JSON binding + struct validation via `go-playground/validator/v10`. Validation errors are automatically converted to the canonical error response shape.

**Where each type of validation belongs:**

| Layer | Validates | Example |
|---|---|---|
| DTO struct tags | Format | `required`, `email`, `e164`, `min=N`, `max=N` |
| Handler | Path/query params | Parse UUID, convert types |
| Service | Business rules | "email already registered", "limit max 100" |
| Domain methods | State machine | "can only deactivate an active user" |

Don't put business rules in validation tags. Tags handle format ("is this a valid email?"), services handle meaning ("is this email already registered?").

---

### Request & Response Standards

#### Success response shapes

```json
// Single resource
{
    "success": true,
    "data": { "id": "...", "name": "...", ... }
}

// List with pagination
{
    "success": true,
    "data": [ { ... }, { ... } ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 143,
        "total_pages": 8
    }
}

// Delete
{
    "success": true,
    "data": null
}
```

#### Error response shape

```json
{
    "type":      "NOT_FOUND_ERROR",
    "status":    404,
    "message":   "user not found",
    "path":      "/users/abc-123",
    "method":    "GET",
    "timestamp": "2024-11-03T14:22:05+05:00"
}
```

#### Response helper functions

**Location**: `internal/shared/response/`

These are the only functions that produce HTTP responses. Use them everywhere.

| Function | Status | Use for |
|---|---|---|
| `response.OK(c, dto)` | 200 | Get one, update |
| `response.Created(c, dto)` | 201 | Create |
| `response.Paginated(c, data, page, perPage, total)` | 200 | List endpoints |
| `response.OK(c, nil)` | 200 | Delete |

```go
// Bad — never do this
return c.JSON(200, map[string]any{"data": user})
return c.JSON(200, "user updated successfully")

// Good
return response.OK(c, dto.UserResFromDomain(user))
return response.Created(c, dto.UserResFromDomain(user))
return response.Paginated(c, dto.UserListResFromDomain(users), p.Page, p.PerPage, total)
return response.OK(c, nil)
```

**Always return the resource DTO from mutations.** Returning a string message forces clients to make a second GET to update local state.

---

## Docker

**Docker is not optional.** Every developer, CI pipeline, and deployment environment runs the application via Docker. There is no supported path to run the application directly on a host machine.

### Files

```
Dockerfile             # Multi-stage build: builder → final minimal image
docker-compose.yml     # App + PostgreSQL + Redis — everything needed to run
.dockerignore
```

### `docker-compose.yml` structure

```yaml
version: "3.9"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file: configs/.env
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

### `Dockerfile` structure

```dockerfile
# Stage 1 — build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Stage 2 — final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY configs/ configs/
EXPOSE 8080
CMD ["./server"]
```

### Makefile targets

```makefile
up:           ## Start all services
	docker compose up --build -d

down:         ## Stop all services
	docker compose down

logs:         ## Tail app logs
	docker compose logs -f app

seed:         ## Run the seed (permissions + group_permissions)
	docker compose exec app ./server seed

migrate-up:   ## Apply all pending migrations
	docker compose exec app ./server migrate up

migrate-down: ## Rollback last migration
	docker compose exec app ./server migrate down

migrate-create name=?:  ## Create a new migration file pair
	docker compose exec app ./server migrate create $(name)
```

**Rules**:
- Never run `go run` or `./server` directly on the host — always via `docker compose exec` or `make`.
- The `app` service waits for the `db` healthcheck before starting. Do not remove this dependency.
- `configs/.env` is gitignored. `configs/.env.example` is committed and must be kept in sync.
- All environment variables are loaded from `.env` via Viper. No hardcoded config values in code.

---

## Common Pitfalls

### Always propagate request context

**Never use `context.Background()` in request code.** It breaks timeouts, cancellation, and tracing.

```go
// Bad
_, err = r.db.Exec(context.Background(), query, args...)

// Good
_, err = r.db.Exec(ctx, query, args...)
```

### Nil slices marshal to `null`, not `[]`

Client code doing `data.items.map(...)` will crash on `null`. Always initialize:

```go
users := make([]*domain.User, 0)
```

### SELECT column order must match Scan order

Positional scanning means a mismatch won't error — it will silently put wrong data in wrong fields. When adding a column, add it in the **same position** in both `Select()` and `Scan()`. When reviewing PRs, always cross-check these two lists.

### `defer rows.Close()` goes AFTER the error check

If `Query()` returns an error, `rows` is nil. Deferring `Close()` before the check will panic.

```go
rows, err := r.db.Query(ctx, query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()
```

### Enforce pagination limits

Never trust client-provided `per_page`. `reqctx.ParsePagination` handles clamping — never do it again in a handler or service.

### Import ordering

Group in three blocks: standard library, external packages, internal packages.

```go
import (
    "context"
    "fmt"

    sq "github.com/Masterminds/squirrel"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/yourorg/boilerplate/internal/domain"
    "github.com/yourorg/boilerplate/internal/shared/app_errors"
)
```

### Never check permissions in handlers or services

Permission checking belongs to middleware only. A handler that calls `service.HasPermission(...)` is a pattern violation — add the `middleware.RequirePermission("x.y")` to the route instead.

---

## Git Workflow

### Branch Naming

```
feature/short-description     # New features
fix/short-description          # Bug fixes
refactor/short-description     # Code restructuring
docs/short-description         # Documentation only
chore/short-description        # Tooling, dependencies, CI
```

### Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/): `{type}: {short description in imperative mood}`

```
feat: add user filtering by status
fix: handle null email in user list query
refactor: extract user DAO from mapper file
chore: update go dependencies
```

Keep under 72 characters. Imperative mood: "add", "fix", "remove" — not "added", "fixes", "removed".

### Pull Requests

- Branch from `develop`, PR back to `develop`.
- Keep PRs focused — one feature or fix per PR.

---

## Database Migrations

**Location**: `migrations/`
**Tool**: [golang-migrate](https://github.com/golang-migrate/migrate)

Files follow `{sequence}_{description}.{up,down}.sql` naming (6-digit zero-padded sequence):

```
migrations/
  000001_create_users_table.up.sql
  000001_create_users_table.down.sql
  000002_create_permissions_table.up.sql
  000002_create_permissions_table.down.sql
  000003_create_group_permissions_table.up.sql
  000003_create_group_permissions_table.down.sql
  000004_create_system_logs_table.up.sql
  000004_create_system_logs_table.down.sql
```

**Rules**:
- Every up migration **must** have a corresponding down migration.
- Down migrations must be safe to re-run (`IF EXISTS`, `IF NOT EXISTS`).
- Never modify an existing migration that has been merged to `develop`. Create a new one.
- Add `NOT NULL` with a `DEFAULT` when adding required columns to existing tables.

### Core migration examples

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    phone         VARCHAR(20)  NOT NULL UNIQUE,
    role          VARCHAR(20)  NOT NULL DEFAULT 'user',
    password_hash TEXT         NOT NULL,
    status        VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

```sql
-- 000002_create_permissions_table.up.sql
CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL      PRIMARY KEY,
    module      VARCHAR(50) NOT NULL,
    action      TEXT        NOT NULL,
    name        VARCHAR(100) NOT NULL UNIQUE,  -- e.g. "users.get_all"
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```sql
-- 000003_create_group_permissions_table.up.sql
CREATE TABLE IF NOT EXISTS group_permissions (
    id             SERIAL      PRIMARY KEY,
    name           VARCHAR(100) NOT NULL UNIQUE,  -- e.g. "admin", "moderator"
    permission_ids INTEGER[]   NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Link users to a group permission
ALTER TABLE users ADD COLUMN IF NOT EXISTS group_permission_id INTEGER REFERENCES group_permissions(id);
```

```sql
-- 000004_create_system_logs_table.up.sql
CREATE TABLE IF NOT EXISTS system_logs (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    level      VARCHAR(10) NOT NULL,  -- debug | info | warn | error
    type       VARCHAR(20) NOT NULL,  -- request | response | system | audit
    message    TEXT        NOT NULL,
    context    JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_system_logs_level   ON system_logs(level);
CREATE INDEX IF NOT EXISTS idx_system_logs_type    ON system_logs(type);
CREATE INDEX IF NOT EXISTS idx_system_logs_created ON system_logs(created_at DESC);
```

### Seed

`cmd/seed/main.go` runs after migrations and is idempotent (safe to re-run):

1. Upserts all permissions from the seed list into `permissions`.
2. Creates default group_permissions (`admin`, `moderator`, `read_only`) if they don't exist.
3. Assigns `admin` group to the first user (or a user whose email matches `ADMIN_EMAIL` in config).

Run via:

```bash
make seed
# or directly:
docker compose exec app ./server seed
```