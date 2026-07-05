# Operations

## Health probes

| Endpoint | Purpose | Semantics |
|---|---|---|
| `GET /health` | Basic status | Always 200 if the process serves requests |
| `GET /livez` | Liveness | 200 if the process is up (no dependency checks) — use for restart decisions |
| `GET /readyz` | Readiness | 200 when the **database** is reachable, else **503**. Redis is reported (`ok`/`unavailable`/`disabled`) but does not fail readiness, since the app runs degraded without it. Use for load-balancer/K8s traffic gating. |

The Docker image ships a `HEALTHCHECK` that polls `/livez`.

## Logging

Structured logs via `log/slog`:
- `APP_ENV=production` → JSON; otherwise human-readable text.
- `LOG_LEVEL` controls verbosity (`debug|info|warn|error`).
- Each HTTP request logs one `http_request` line with `request_id`, `method`,
  `path`, `status`, `latency_ms`, `ip`. Level scales with status (5xx=error,
  4xx=warn).
- The standard `log` package is bridged into slog, so legacy `log.Printf` calls
  are captured too.

## Request correlation

Every request gets an `X-Request-ID` (reused from an inbound header if present,
otherwise generated). It is returned in the response header and included in
error response bodies (`request_id`) and access logs — quote it in bug reports.

## Error responses

All API errors share:

```json
{ "success": false, "error": "Resource not found", "code": "NOT_FOUND", "request_id": "…" }
```

`code` is stable/machine-readable (see `openapi.yaml` for the full list). 5xx
responses use a generic message and never leak internal details.

## Rate limiting

`/api/admin/login` and `/api/admin/register` are limited to **10 requests/minute
per IP**; exceeding it returns `429` with `code: RATE_LIMITED`.

## Graceful shutdown

On `SIGINT`/`SIGTERM` the server stops accepting connections and drains in-flight
requests up to `SHUTDOWN_TIMEOUT_SEC`, then closes the scheduler, event bus,
Redis, session store and database.

## Timeouts

Outbound LINE/OCR calls and the HTTP server have explicit timeouts (see
[CONFIGURATION.md](CONFIGURATION.md)) so a stuck dependency can't hang requests
indefinitely.

## Database migrations

```bash
make migrate-up      # apply
make migrate-down    # roll back one
make migrate-version # current version
```
