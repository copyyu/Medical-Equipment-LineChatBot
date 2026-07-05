# Architecture

The service follows **Clean Architecture** — dependencies point inward, and the
application layer depends on abstractions (ports), not concrete infrastructure.

```
cmd/                        entry points (app, migrate)
internal/
  domain/                   entities, repository interfaces, ports (no framework deps)
    line/entity             GORM entities + enums
    line/repository         repository interfaces + ErrDuplicate sentinel
    event                   EventBus port
    port                    TxManager port
  application/              use cases + services (orchestration, business rules)
    usecase, service, dto, mapper
  infrastructure/           concrete implementations (GORM, Redis, LINE, OCR)
    persistence             repository implementations
    database                connection + TxManager implementation
    redis                   event bus + idempotency store
    client                  LINE / OCR HTTP clients
    logger                  slog setup
    bootstrap               dependency wiring (composition root)
  interfaces/http           handlers, routes, middleware
```

## Key design decisions

### Transactions (unit of work)
`port.TxManager` exposes `WithTransaction(ctx, func(ctx) error)`. The GORM-backed
implementation stores the active `*gorm.DB` in the context; repositories resolve
their handle via `database.DBFromContext(ctx, fallback)`, so any repository call
made with that context transparently joins the transaction. The port carries **no
ORM types**, keeping the application layer decoupled from GORM. It delegates to
`db.Transaction`, so commit/rollback and panic recovery are handled correctly.

Multi-step writes (ticket + history) run inside one transaction so a partial
failure can never leave a ticket without its audit history.

### Webhook idempotency
LINE delivers webhooks at-least-once. Each event is de-duplicated by its
`webhookEventId` using an atomic Redis `SET NX` (`redis.IdempotencyStore`). The
store fails open (process the event) when Redis is unavailable, rather than
silently dropping events.

### Concurrency safety
- Ticket numbers are generated from the latest row; the `ticket_no` unique index
  is the source of truth, and creation retries on a collision.
- Fire-and-forget goroutines (notifications, event publishing) run through
  `goSafe`, which recovers panics so a background failure can't crash the process.

### Error handling & observability
- A single response envelope (`success`, `error`, `code`, `request_id`) is used
  across the API; 5xx responses never leak internal error strings.
- Structured logging via `log/slog` (JSON in production). `slog.SetDefault` also
  bridges the standard `log` package. Every request carries an `X-Request-ID`.

### Data & schema
Schema is managed with `golang-migrate` (`migrations/`, `cmd/migrate`).
`AutoMigrate` is **not** used in production.
