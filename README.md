# Medical Equipment LINE ChatBot

> ระบบจัดการเครื่องมือแพทย์ผ่าน LINE Official Account — Backend API พัฒนาด้วย Go (Fiber Framework) ออกแบบตาม Clean Architecture

[![CI](https://github.com/copyyu/Medical-Equipment-LineChatBot/actions/workflows/ci.yml/badge.svg)](https://github.com/copyyu/Medical-Equipment-LineChatBot/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)

---

## สารบัญ

- [ภาพรวมโครงการ](#ภาพรวมโครงการ)
- [ฟีเจอร์หลัก](#ฟีเจอร์หลัก)
- [สถาปัตยกรรมระบบ](#สถาปัตยกรรมระบบ)
- [เทคโนโลยีที่ใช้](#เทคโนโลยีที่ใช้)
- [โครงสร้างโปรเจกต์](#โครงสร้างโปรเจกต์)
- [การติดตั้งและรัน](#การติดตั้งและรัน)
- [ตัวแปรสิ่งแวดล้อม](#ตัวแปรสิ่งแวดล้อม)
- [API Endpoints](#api-endpoints)
- [CI/CD](#cicd)
- [Entity Relationships](#entity-relationships)

---

## ภาพรวมโครงการ

**Medical Equipment LINE ChatBot** คือระบบ Backend สำหรับบริหารจัดการเครื่องมือแพทย์ในโรงพยาบาล โดยรองรับการใช้งานผ่านช่องทางต่าง ๆ ดังนี้:

- สอบถามข้อมูลเครื่องมือแพทย์ผ่าน **LINE Official Account**
- ถ่ายภาพรหัสเครื่องมือเพื่อให้ระบบ **AI OCR** อ่านรหัสและค้นหาข้อมูลอัตโนมัติ
- แจ้งซ่อมและติดตามสถานะ Ticket ผ่าน LINE ได้ทันที
- ผู้ดูแลระบบสามารถจัดการข้อมูลผ่าน **REST API** และรับ **Real-time Events** ผ่าน Server-Sent Events (SSE)

---

## ฟีเจอร์หลัก

### LINE ChatBot

| ฟีเจอร์ | รายละเอียด |
|---------|-----------|
| ค้นหาเครื่องมือ | พิมพ์รหัสเครื่อง หรือชื่อเครื่อง เพื่อค้นหาข้อมูล |
| ถ่ายรูป OCR | ส่งรูปรหัสเครื่อง แล้วระบบจะใช้ AI อ่านรหัสและค้นหาเครื่องมือให้อัตโนมัติ |
| แจ้งซ่อม | เปิด Ticket แจ้งซ่อมเครื่องมือผ่าน LINE โดยตรง |
| ติดตามสถานะ | ดูสถานะ Ticket แจ้งซ่อมแบบ real-time |
| Postback Actions | เมนูโต้ตอบแบบ Rich Menu, Quick Reply และ Carousel |
| ข้อความต้อนรับ | ส่งข้อความต้อนรับอัตโนมัติเมื่อผู้ใช้ Follow บัญชี |

### ระบบจัดการสำหรับผู้ดูแล (Admin)

| ฟีเจอร์ | รายละเอียด |
|---------|-----------|
| Dashboard | สรุปภาพรวมเครื่องมือ, Ticket และสถิติต่าง ๆ |
| จัดการเครื่องมือ | CRUD, ค้นหา, กรอง และ Pagination |
| นำเข้า Excel | Import ข้อมูลเครื่องมือจากไฟล์ Excel (.xlsx) |
| ระบบ Ticket | จัดการ Ticket แจ้งซ่อม พร้อมบันทึกประวัติ (History) |
| การแจ้งเตือน | ตั้ง Schedule แจ้งเตือนบำรุงรักษา (PM/Calibration) |
| Activity Log | บันทึกกิจกรรมทุกการเปลี่ยนแปลงในระบบ |
| Real-time SSE | Event Stream สำหรับอัปเดต UI แบบ real-time |
| Authentication | ระบบล็อกอิน Admin พร้อม Session-based Authentication |

### AI OCR

| ฟีเจอร์ | รายละเอียด |
|---------|-----------|
| อ่านรหัสเครื่อง | OCR จากรูปถ่ายรหัสเครื่องมือแพทย์ |
| Fuzzy Matching | ค้นหาแบบใกล้เคียงด้วย Levenshtein Distance |
| Confidence Score | แสดงค่าความเชื่อมั่นของผลลัพธ์ OCR |

---

## สถาปัตยกรรมระบบ

โปรเจกต์ออกแบบตาม **Clean Architecture** (Layered Architecture) แบ่งออกเป็น 4 ชั้นหลัก:

```
+--------------------------------------------------+
|                 Interfaces Layer                 |
|        (HTTP Handlers, Routes, Middleware)        |
+--------------------------------------------------+
|               Application Layer                  |
|         (Use Cases, Services, DTOs, Mappers)     |
+--------------------------------------------------+
|                 Domain Layer                     |
|      (Entities, Repository Interfaces, Events)   |
+--------------------------------------------------+
|             Infrastructure Layer                 |
|    (Database, Redis, LINE Client, OCR Client)    |
+--------------------------------------------------+
```

### หลักการออกแบบ

- **Dependency Inversion** — ชั้นบนไม่ขึ้นตรงกับชั้นล่าง โดยใช้ Interface เป็นสัญญา (Contract)
- **Separation of Concerns** — แยกความรับผิดชอบอย่างชัดเจนในแต่ละ Layer
- **Event-Driven** — ใช้ Redis Pub/Sub เป็น Event Bus สำหรับ real-time events (SSE)

---

## เทคโนโลยีที่ใช้

| เทคโนโลยี | รายละเอียด |
|-----------|-----------|
| Go 1.24 | ภาษาหลักในการพัฒนา |
| Fiber v2 | Web Framework ประสิทธิภาพสูง (สร้างบน fasthttp) |
| PostgreSQL 16 | ฐานข้อมูลหลัก |
| GORM | ORM สำหรับจัดการฐานข้อมูล |
| Redis 7 | Event Bus (Pub/Sub) สำหรับ real-time SSE |
| LINE Bot SDK v8 | เชื่อมต่อ LINE Messaging API |
| Docker / Docker Compose | Containerization และ Orchestration |
| GitHub Actions | CI/CD Pipeline |

---

## โครงสร้างโปรเจกต์

```
Medical-Equipment-LineChatBot/
├── cmd/
│   └── app/
│       └── main.go                    # Entry point ของแอปพลิเคชัน
├── internal/
│   ├── config/
│   │   └── config.go                  # โหลด Environment Variables
│   ├── domain/                        # Domain Layer
│   │   ├── constants/                 # ค่าคงที่ (Actions, Messages)
│   │   ├── event/                     # Event Bus interface และ Event model
│   │   └── line/
│   │       ├── entity/                # Entities (Equipment, Ticket, Admin, ...)
│   │       ├── model/                 # LINE Message models
│   │       └── repository/            # Repository Interfaces
│   ├── application/                   # Application Layer
│   │   ├── dto/                       # Data Transfer Objects
│   │   ├── mapper/                    # Entity <-> DTO Mappers
│   │   ├── service/                   # Business Services
│   │   └── usecase/                   # Use Cases (Business Logic)
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── bootstrap/                 # Dependency Injection และ App Initialization
│   │   ├── client/                    # External Clients (LINE, OCR)
│   │   ├── database/                  # PostgreSQL Connection
│   │   ├── line/                      # LINE Message Templates
│   │   ├── persistence/               # Repository Implementations (GORM)
│   │   ├── redis/                     # Redis Connection และ Event Bus Implementation
│   │   └── session/                   # Session Store (OCR Confirmations)
│   ├── interfaces/                    # Interfaces Layer
│   │   └── http/
│   │       ├── handlers/              # HTTP Handlers (Webhook, Equipment, ...)
│   │       ├── middleware/            # CORS, Logger, Auth Middleware
│   │       └── routes/                # Route Definitions
│   └── utils/                         # Utilities
│       ├── errors/                    # Custom Error Types
│       ├── ptr/                       # Pointer Helpers
│       ├── scheduler/                 # Cron-based Notification Scheduler
│       └── token/                     # Token Generation
├── .github/
│   └── workflows/
│       └── ci.yml                     # GitHub Actions CI Pipeline
├── Dockerfile                         # Multi-stage Docker Build
├── docker-compose.yml                 # PostgreSQL + Redis + App
├── .env.example                       # ตัวอย่างตัวแปรสิ่งแวดล้อม
├── go.mod                             # Go Module Dependencies
└── go.sum                             # Dependency Checksums
```

---

## การติดตั้งและรัน

### ข้อกำหนดเบื้องต้น

- [Go 1.24+](https://golang.org/dl/)
- [Docker และ Docker Compose](https://docs.docker.com/get-docker/)
- [LINE Developers Account](https://developers.line.biz/) (สำหรับ Channel Token และ Secret)

### วิธีที่ 1: Docker Compose (แนะนำ)

```bash
# 1. Clone โปรเจกต์
git clone https://github.com/copyyu/Medical-Equipment-LineChatBot.git
cd Medical-Equipment-LineChatBot

# 2. สร้างไฟล์ .env จากตัวอย่าง
cp .env.example .env
# แก้ไขค่าใน .env ให้ตรงกับ LINE Channel ของคุณ

# 3. รันทุก Service ด้วย Docker Compose
docker compose up -d

# 4. ดู Logs
docker compose logs -f app
```

ระบบจะรัน 3 Services:

| Service | Port | รายละเอียด |
|---------|------|-----------|
| app | 3000 | Go Fiber API Server |
| db | 5432 | PostgreSQL 16 |
| redis | 6379 | Redis 7 |

### วิธีที่ 2: รันแบบ Development (Local)

```bash
# 1. รัน PostgreSQL และ Redis ด้วย Docker
docker compose up -d db redis

# 2. ติดตั้ง Dependencies
go mod download

# 3. สร้างไฟล์ .env
cp .env.example .env
# แก้ DB_HOST=localhost, REDIS_URL=redis://localhost:6379

# 4. รัน Server
go run ./cmd/app

# 5. (ทางเลือก) ใช้ Air สำหรับ Hot Reload
# ติดตั้ง: go install github.com/air-verse/air@latest
air
```

### ตรวจสอบว่าระบบทำงาน

```bash
# Health Check
curl http://localhost:3000/health

# Root Endpoint
curl http://localhost:3000/
# Response: {"message":"Medical Equipment Webhook Server","status":"running","version":"1.0.0"}
```

---

## ตัวแปรสิ่งแวดล้อม

สร้างไฟล์ `.env` จาก `.env.example`:

| ตัวแปร | ค่าตัวอย่าง | คำอธิบาย |
|--------|------------|---------|
| `LINE_CHANNEL_TOKEN` | `your_token` | LINE Messaging API Channel Access Token |
| `LINE_CHANNEL_SECRET` | `your_secret` | LINE Messaging API Channel Secret |
| `PORT` | `3000` | พอร์ตที่ Server ทำงาน |
| `DB_HOST` | `db` | PostgreSQL Host (`db` สำหรับ Docker, `localhost` สำหรับ local) |
| `DB_PORT` | `5432` | PostgreSQL Port |
| `DB_USER` | `postgres` | PostgreSQL Username |
| `DB_PASSWORD` | `postgres` | PostgreSQL Password |
| `DB_NAME` | `medical_equipment` | ชื่อฐานข้อมูล |
| `REDIS_URL` | `redis://redis:6379` | Redis Connection URL |
| `OCR_API_URL` | `http://ocr-service:8000` | URL ของ OCR API (ไม่บังคับ) |
| `BASE_URL` | `https://your-domain.com` | Base URL ของ Server |
| `CONTACT_CENTER_NAME` | `ศูนย์เครื่องมือแพทย์` | ชื่อศูนย์สำหรับแสดงใน LINE |
| `CONTACT_PHONE` | `02-xxx-xxxx` | เบอร์โทรติดต่อ |
| `CONTACT_EMAIL` | `contact@hospital.com` | อีเมลติดต่อ |
| `CONTACT_EMERGENCY_PHONE` | `02-xxx-xxxx` | เบอร์ฉุกเฉิน |
| `CONTACT_WORKING_HOURS` | `จ-ศ 08:00-17:00` | เวลาทำการ |

---

## API Endpoints

### สาธารณะ (Public)

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/` | Root — แสดงสถานะ Server |
| `GET` | `/health` | Health Check |
| `POST` | `/webhook` | LINE Webhook Callback |
| `GET` | `/api/events/stream` | SSE Event Stream (query: `?types=equipment.updated,ticket.created`) |

### Authentication

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `POST` | `/api/admin/login` | เข้าสู่ระบบ Admin |

### Protected (ต้องผ่านการยืนยันตัวตน)

#### Dashboard

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/api/dashboard` | ข้อมูลสรุป Dashboard |

#### Equipment (เครื่องมือแพทย์)

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/api/equipment` | รายการเครื่องมือ (รองรับ filter, pagination) |
| `GET` | `/api/equipment/:id` | รายละเอียดเครื่องมือ |
| `PUT` | `/api/equipment/:id` | อัปเดตข้อมูลเครื่องมือ |
| `DELETE` | `/api/equipment/:id` | ลบเครื่องมือ |

#### Import (นำเข้าข้อมูล)

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `POST` | `/api/import/upload` | อัปโหลดไฟล์ Excel เพื่อนำเข้าข้อมูล |

#### Ticket (แจ้งซ่อม)

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/api/tickets` | รายการ Ticket |
| `GET` | `/api/tickets/:id` | รายละเอียด Ticket |
| `PUT` | `/api/tickets/:id` | อัปเดตสถานะ Ticket |

#### Notification (การแจ้งเตือน)

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/api/notifications` | รายการการแจ้งเตือน |
| `POST` | `/api/notifications` | สร้างการแจ้งเตือนใหม่ |

#### Activity Log

| Method | Path | คำอธิบาย |
|--------|------|---------|
| `GET` | `/api/activity-logs` | ประวัติกิจกรรมทั้งหมด |

---

## CI/CD

โปรเจกต์ใช้ **GitHub Actions** สำหรับ Continuous Integration:

```yaml
# .github/workflows/ci.yml
Trigger: push / pull_request -> main, master

Jobs:
  - Checkout code
  - Setup Go 1.24
  - Download dependencies
  - Build
  - Run tests
```

---

## Entity Relationships

### Equipment (เครื่องมือแพทย์)

- มีความสัมพันธ์กับ **EquipmentModel** (ยี่ห้อ/รุ่น), **Department** (แผนก) และ **MaintenanceRecord** (ประวัติบำรุงรักษา)
- สถานะเครื่องมือ: `active`, `defective`, `wait_decom`, `decommission`, `active_ready_to_sell`, `missing`, `plan_to_replace`

### Ticket (ใบแจ้งซ่อม)

- รูปแบบเลข Ticket: `REQ-YYYY-XXXXX`
- สถานะ: `in_process` (กำลังดำเนินการ), `return_equipment_back` (ส่งคืนเครื่องแล้ว), `send_to_outsource` (ส่งซ่อมภายนอก)
- ระดับความเร่งด่วน: `low`, `medium`, `high`, `urgent`

---

