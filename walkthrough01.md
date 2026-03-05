# 🏥 Code Review — Medical-Equipment-LineChatBot

## สรุปโครงสร้างโปรเจกต์

```
Medical-Equipment-LineChatBot/
├── cmd/app/main.go                    # Entry point + Graceful shutdown
├── internal/
│   ├── config/config.go               # Environment configuration
│   ├── application/                   # Application Layer
│   │   ├── dto/        (8 files)      # Data Transfer Objects
│   │   ├── mapper/     (1 file)       # Entity ↔ DTO mappers
│   │   ├── service/    (7 files)      # Application services
│   │   └── usecase/    (11 files)     # Use cases (business orchestration)
│   ├── domain/                        # Domain Layer
│   │   ├── constants/  (2 files)      # Action constants
│   │   └── line/
│   │       ├── entity/    (14 files)  # Domain entities
│   │       ├── model/     (2 files)   # Value objects
│   │       └── repository/(13 files)  # Repository interfaces
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── bootstrap/app.go           # DI Container + App init
│   │   ├── client/    (3 files)       # External API clients
│   │   ├── database/  (2 files)       # DB connection + seeder
│   │   ├── line/templates/ (16 files) # LINE Flex Message templates
│   │   ├── persistence/ (13 files)    # Repository implementations
│   │   └── session/   (1 file)       # In-memory session store
│   ├── interfaces/http/               # Interface Layer
│   │   ├── handlers/  (9 files)       # HTTP handlers
│   │   ├── middleware/(3 files)       # Auth, CORS, Logger
│   │   └── routes/    (10 files)      # Route registration
│   └── utils/                         # Utilities
│       ├── errors/    (3 files)       # Error types + HTTP responses
│       ├── ptr/       (1 file)        # Pointer helpers
│       ├── scheduler/ (1 file)        # Notification scheduler
│       └── token/     (1 file)        # Token generator
├── .env                               # ⚠️ Contains REAL secrets
├── .gitignore
├── go.mod / go.sum
└── .air.toml                          # Hot reload config
```

---

## ✅ สิ่งที่ทำได้ดี

| ด้าน | รายละเอียด |
|------|-----------|
| **Clean Architecture** | แบ่ง layer ชัดเจน: `domain` → `application` → `infrastructure` → `interfaces` |
| **Dependency Injection** | ใช้ constructor injection ผ่าน [app.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/bootstrap/app.go) ไม่มี hidden dependency |
| **Repository Pattern** | Domain กำหนด interface, Infrastructure implement — แยก concern ได้ดี |
| **Error Handling** | มี centralized error mapping [response.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/utils/errors/response.go) — map domain error → HTTP status |
| **Password Security** | ใช้ `bcrypt` hashing ใน [admin_service.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/application/service/admin_service.go) — best practice |
| **Graceful Shutdown** | [main.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/cmd/app/main.go) จัดการ signal + cleanup resources ครบ |
| **Session Cleanup** | [session_store.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/session/session_store.go) มี goroutine cleanup + graceful close |
| **Webhook Signature** | ตรวจสอบ LINE signature ผ่าน `webhook.ParseRequest` ก่อนประมวลผล |

---

## 🔴 ช่องโหว่ด้านความปลอดภัย (Security Vulnerabilities)

### 1. 🚨 Secrets ถูก Expose ใน [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) (CRITICAL)

> [!CAUTION]
> ไฟล์ [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) มี LINE Channel Token, Channel Secret, และ DB password จริง! ถ้า [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) หลุดเข้า Git history จะทำให้ credentials ทั้งหมดถูก compromise

**ไฟล์**: [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env)

```
LINE_CHANNEL_TOKEN=G2+t5K7oCV6uxnGMCb...  ← Token จริง!
LINE_CHANNEL_SECRET=4e28ea2ad3e766b644...  ← Secret จริง!
DB_PASSWORD=4321                           ← รหัสผ่าน DB จริง!
```

**แนวทางแก้ไข:**
- ลบ [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) ออกจาก repository ทันที (`git rm --cached .env`)
- ใช้ `.env.example` แทน (มีแค่ key ไม่มี value)
- ถ้าเคย commit [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) แล้ว → **ต้อง rotate ทุก credentials**
- พิจารณาใช้ secret manager (Vault, AWS Secrets Manager)

---

### 2. 🚨 Register Endpoint เปิด Public (CRITICAL)

> [!CAUTION]
> ใครก็สามารถสร้าง Admin account ได้โดยไม่ต้อง login! ไม่มี Authorization

**ไฟล์**: [admin_routes.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/interfaces/http/routes/admin_routes.go#L20)

```go
// Public routes - ไม่ต้อง auth
admin.Post("/login", adminHandler.Login)
admin.Post("/register", adminHandler.Register)  // ← ❌ ใครก็ register ได้!
```

**แนวทางแก้ไข:**
- ย้าย `/register` ไปอยู่ใน protected group (ต้อง login ก่อน)
- หรือใช้ invite-only system (ส่ง invite code)
- เก็บ role check ให้เฉพาะ super admin สร้าง account ได้

---

### 3. ⚠️ CORS เปิดทุก Origin (HIGH)

> [!WARNING]
> `AllowOrigins: "*"` อนุญาตให้ทุกเว็บไซต์เรียก API ได้ — เสี่ยงต่อ CSRF attack

**ไฟล์**: [fiber.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/interfaces/http/middleware/fiber.go#L11)

```go
cors.New(cors.Config{
    AllowOrigins: "*",  // ← ❌ อันตรายใน production!
})
```

**แนวทางแก้ไข:**
- กำหนด origins เฉพาะ frontend domain: `AllowOrigins: "https://your-frontend.com"`
- ย้ายค่านี้ไป env variable เพื่อ config ต่าง environment ได้

---

### 4. ⚠️ ไม่มี Rate Limiting (MEDIUM)

ไม่มี rate limiter สำหรับ login endpoint → เสี่ยงต่อ brute-force attack

**แนวทางแก้ไข:**
```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use("/api/admin/login", limiter.New(limiter.Config{
    Max:        5,
    Expiration: 15 * time.Minute,
}))
```

---

### 5. ⚠️ ไม่มี SSL/TLS สำหรับ Postgres Connection (MEDIUM)

**ไฟล์**: [db_connect.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/database/db_connect.go#L27)

```go
dsn := fmt.Sprintf(
    "host=%s ... sslmode=disable",  // ← ❌ ข้อมูลส่งแบบ plaintext
)
```

ถ้า DB อยู่คนละ server → ข้อมูลจะถูกส่งแบบไม่เข้ารหัส

---

### 6. ⚠️ No Input Validation on Registration (MEDIUM)

ไม่มีการ validate ข้อมูลที่ register เช่น ความยาว password, format ของ email

**ไฟล์**: [admin_service.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/application/service/admin_service.go#L48)

---

## 🟡 ปัญหาด้านสถาปัตยกรรม (Architecture Issues)

### 1. Entity ผูกกับ GORM (Domain Layer Leak)

> [!IMPORTANT]
> Domain entities ไม่ควรพึ่ง infrastructure library

**ไฟล์**: [entity/equipment.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/domain/line/entity/equipment.go#L6-L7)

```go
import "gorm.io/gorm"  // ❌ GORM อยู่ใน domain entity

type Equipment struct {
    DeletedAt gorm.DeletedAt  // ❌ infra concern อยู่ใน domain
    // gorm tags everywhere: `gorm:"primaryKey"`, `gorm:"size:100"`
}
```

**แนวทางแก้ไข:**
- สร้าง domain entity แยก (pure struct ไม่มี gorm tag)
- ให้ infrastructure layer มี GORM model ของตัวเอง + mapper

---

### 2. Global Database Variables

**ไฟล์**: [db_connect.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/database/db_connect.go#L15-L18)

```go
var (
    DB    *gorm.DB   // ← ❌ Global state
    SqlDB *sql.DB
)
```

ทำให้ test ยาก, ไม่ thread-safe ในบาง scenario, ขัดกับ DI pattern ที่ใช้อยู่

**แนวทางแก้ไข:**
- Return `*gorm.DB` จาก [Connect()](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/database/db_connect.go#20-71) แล้ว inject เข้า repos ผ่าน constructor

---

### 3. Config Validation ถูก Comment Out

**ไฟล์**: [db_connect.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/infrastructure/database/db_connect.go#L22-L25)

```go
// if cfg.DB.Host == "" || cfg.DB.User == "" ...  ← ❌ ถูก comment ไว้!
```

ถ้าค่าใน [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) ขาดไป → app จะ crash แบบไม่มี error message ที่ชัดเจน

---

### 4. Repositories ไม่ Consistent เรื่อง `context.Context`

บาง methods รับ `ctx context.Context` บางอันไม่รับ:

```go
// EquipmentRepository interface
FindByIDCode(idCode string) (*entity.Equipment, error)     // ← ❌ ไม่มี ctx
Create(ctx context.Context, equipment *entity.Equipment) error  // ← ✅ มี ctx
```

ควรใช้ `context.Context` ทุก method เพื่อรองรับ timeout, cancellation, tracing

---

### 5. Duplicated Token Generation Logic

มี 2 ที่ที่ generate token แบบเดียวกัน:
- [token/generator.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/utils/token/generator.go)
- [admin_service.go](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/application/service/admin_service.go#L209-L215) (private [generateToken()](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/application/service/admin_service.go#209-216))

ควรใช้ `token.GenerateToken()` จากที่เดียว

---

### 6. Domain Naming ไม่สื่อ

```
internal/domain/line/  ← ❌ "line" ไม่ใช่ชื่อ domain
```

ชื่อ `line` เป็น infrastructure concern (LINE messaging platform) ไม่ใช่ business domain. Entity อย่าง [Equipment](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/domain/line/entity/equipment.go#66-100), `Ticket`, [Admin](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/application/service/admin_service.go#205-208) ควรอยู่ใน domain ที่ชื่อสื่อความหมายกว่า เช่น `medical` หรือ `equipment`

---

### 7. [EquipmentRepository](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/internal/domain/line/repository/equipment_repository.go#8-48) Interface ใหญ่เกิน

Interface มี **20+ methods** — ละเมิดหลัก Interface Segregation Principle (ISP)

ควรแยกเป็น interfaces ย่อย เช่น:
- `EquipmentReader` — Find, Get, Count
- `EquipmentWriter` — Create, Update, Delete
- `EquipmentDashboard` — CountExpired, CountByStatus

---

## 🔵 สิ่งที่ขาดหายไป (Missing)

| รายการ | สถานะ |
|--------|--------|
| **Unit Tests** | ❌ ไม่มีแม้แต่ไฟล์เดียว (`*_test.go` = 0) |
| **API Documentation** | ❌ ไม่มี Swagger/OpenAPI spec |
| **Structured Logging** | ⚠️ ใช้ `log.Println` — ควรใช้ structured logger (zerolog, zap) |
| **Request Validation** | ⚠️ ไม่มี struct tag validation (เช่น `validate:"required,email"`) |
| **Database Migrations** | ⚠️ ใช้ `AutoMigrate` — ไม่แนะนำใน production (ใช้ `golang-migrate` แทน) |
| **Health Check** | ✅ มี แต่ Logger middleware ถูก define ไว้ไม่ได้ใช้ |
| **Dockerfile / CI/CD** | ❌ ไม่มี |
| **README.md** | ❌ ไม่มี documentation |

---

## 📊 สรุประดับความรุนแรง

| ระดับ | จำนวน | รายการ |
|-------|--------|--------|
| 🔴 CRITICAL | 2 | Secrets exposed, Open registration |
| 🟠 HIGH | 1 | CORS wildcard |
| 🟡 MEDIUM | 3 | No rate limiting, No SSL for DB, No input validation |
| 🔵 ARCHITECTURE | 7 | GORM in domain, Global DB, etc. |
| ⚪ MISSING | 8 | Tests, docs, CI/CD, etc. |

---

## 🎯 แนะนำลำดับในการแก้ไข

1. **ทันที**: ลบ [.env](file:///c:/Users/copyyu/Medical-Equipment-LineChatBot/.env) ออก, rotate credentials, protect `/register`
2. **สัปดาห์นี้**: Fix CORS, เพิ่ม rate limiting, เพิ่ม input validation
3. **เดือนนี้**: เขียน unit tests, แยก GORM model ออกจาก domain entity
4. **ระยะยาว**: เพิ่ม CI/CD, structured logging, API documentation
