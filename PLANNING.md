# Implementation Plan: Finance Data Processing & Access Control Backend

## 1. Technology Stack
* **Language:** Go (Golang)
* **Web Framework:** Gin (`github.com/gin-gonic/gin`)
* **ORM:** GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`)
* **Database:** PostgreSQL (Primary relational data)
* **Caching/Sessions:** Redis (`github.com/redis/go-redis/v9`) - *Used for rate limiting and dashboard summary caching.*
* **Search (Phase 2):** Elasticsearch - *For advanced full-text search on record notes.*
* **Validation:** `github.com/go-playground/validator/v10`

---

## 2. Project Directory Structure
This project follows a **Vertical Slice Architecture** within the `api` layer to keep logic isolated and modular.

```text
finance-backend/
├── cmd/
│   └── api/
│       └── main.go                 # App entry point & router setup
├── internal/
│   ├── config/                     # Configuration and environment loading
│   ├── db/                         # Connection pools (Postgres, Redis)
│   ├── models/                     # GORM schemas (Database source of truth)
│   ├── middleware/                 # Auth, RBAC, and Global Error handlers
│   ├── errors/                     # Custom app errors & status codes
│   └── api/                        # Feature-sliced API endpoints
│       ├── auth_login/
│       │   ├── handler.go
│       │   ├── request.go
│       │   └── response.go
│       ├── record_create/          # Example of a slice
│       │   ├── handler.go          # Business logic & Gin handler
│       │   ├── request.go          # Request DTOs with validation tags
│       │   └── response.go         # Response DTOs
│       ├── record_list/
│       └── dashboard_summary/
├── pkg/
│   └── utils/                      # Password hashing, JWT utils
├── go.mod
└── README.md
```

---

## 3. Database Schema Design (GORM)
Located in `internal/models/`. Uses standard GORM conventions with soft deletes.

### `models/user.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleAnalyst Role = "analyst"
	RoleViewer  Role = "viewer"
)

type User struct {
	ID           uint           `gorm:"primaryKey"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	Role         Role           `gorm:"type:varchar(20);not null;default:'viewer'"`
	IsActive     bool           `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### `models/record.go`
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type RecordType string

const (
	TypeIncome  RecordType = "income"
	TypeExpense RecordType = "expense"
)

type Record struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null;index"`
	Amount    float64        `gorm:"type:decimal(12,2);not null"`
	Type      RecordType     `gorm:"type:varchar(10);not null"`
	Category  string         `gorm:"type:varchar(50);not null;index"`
	Date      time.Time      `gorm:"not null;index"`
	Notes     string         `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	
	User      User           `gorm:"foreignKey:UserID"`
}
```

---

## 4. Internal Error Handling
A centralized error package ensures consistent API responses and prevents internal stack traces from leaking.

### `errors/errors.go`
```go
package errors

type AppError struct {
	Code    int      `json:"-"` // HTTP Status Code
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

func NewBadRequest(msg string, details ...string) *AppError {
	return &AppError{Code: 400, Message: msg, Details: details}
}
// Additional constructors for Unauthorized (401), Forbidden (403), etc.
```

---

## 5. API Slice Example: `record_create`

### `request.go`
```go
package record_create

import "time"

type CreateRecordRequest struct {
	Amount   float64   `json:"amount" binding:"required,gt=0"`
	Type     string    `json:"type" binding:"required,oneof=income expense"`
	Category string    `json:"category" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
	Notes    string    `json:"notes"`
}
```

### `handler.go`
```go
package record_create

import (
	"net/http"
	"[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)"
	"finance-backend/internal/models"
	"finance-backend/internal/errors"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) Handle(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequest("Validation failed", err.Error()))
		return
	}

	userID := c.GetUint("user_id")
	record := models.Record{
		UserID: userID,
		Amount: req.Amount,
		Type:   models.RecordType(req.Type),
		Category: req.Category,
		Date: req.Date,
		Notes: req.Notes,
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.Error(err) // Handled by global middleware
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": record.ID, "status": "success"})
}
```

---

## 6. Access Control Logic (RBAC)
Middleware enforces role checks before the request reaches the handler.

### `middleware/role.go`
```go
package middleware

import (
	"[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)"
	"finance-backend/internal/models"
)

func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := models.Role(c.GetString("user_role"))
        for _, r := range allowedRoles {
            if r == userRole {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(403, gin.H{"error": "Access denied"})
    }
}
```

---

## 7. Performance & Features
* **Redis:** Caches dashboard summaries (Total Income/Expense) for 5 minutes. Cache is invalidated on new record creation.
* **GORM Hooks:** Use `AfterSave` hooks on the Record model to automatically clear Redis cache.
* **Filtering:** `record_list` will support GORM scopes for filtering by `date_range`, `category`, and `type`.

---

## 8. Implementation Steps
1.  **Project Init:** Set up `main.go` and Dockerized Postgres/Redis.
2.  **Auth Layer:** Implement JWT login and RBAC middleware.
3.  **Core CRUD:** Build record management slices (Create, Read, Update, Delete).
4.  **Analytics:** Build the `dashboard_summary` with GORM aggregations and Redis.
5.  **Validation:** Add comprehensive input check and custom error responses.