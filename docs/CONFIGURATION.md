# Configuration

All configuration is read from environment variables at startup (a local `.env`
file is loaded if present). Required variables are validated on boot — the app
**fails fast** with a clear error if any are missing (see `config.Validate`).

See [`.env.example`](../.env.example) for a ready-to-copy template.

## Required

| Variable | Description |
|---|---|
| `LINE_CHANNEL_TOKEN` | LINE Messaging API channel access token |
| `LINE_CHANNEL_SECRET` | LINE channel secret (used to verify webhook signatures — an empty value would make signatures forgeable, hence it is required) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_USER` | PostgreSQL user |
| `DB_NAME` | PostgreSQL database name |

## Server & runtime

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | HTTP listen port |
| `BASE_URL` | — | Public base URL (used to build asset/callback links) |
| `APP_ENV` | `development` | `production` emits JSON logs; anything else emits text logs |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |
| `ALLOWED_ORIGINS` | `*` | CORS allow-list (comma-separated). **Set specific origins in production.** |

## Timeouts (seconds)

| Variable | Default | Description |
|---|---|---|
| `HTTP_READ_TIMEOUT_SEC` | `15` | Max time to read a request |
| `HTTP_WRITE_TIMEOUT_SEC` | `120` | Max time to write a response (generous, to accommodate slow webhook/OCR work) |
| `HTTP_IDLE_TIMEOUT_SEC` | `120` | Keep-alive idle timeout |
| `SHUTDOWN_TIMEOUT_SEC` | `15` | Graceful shutdown budget |
| `LINE_API_TIMEOUT_SEC` | `10` | HTTP timeout for outbound LINE API calls |
| `OCR_API_TIMEOUT_SEC` | `90` | HTTP timeout for OCR calls |

## Database pool

| Variable | Default |
|---|---|
| `DB_PASSWORD` | — |
| `DB_MAX_OPEN_CONNS` | `25` |
| `DB_MAX_IDLE_CONNS` | `10` |
| `DB_CONN_MAX_LIFETIME_MIN` | `30` |

## Redis & OCR

| Variable | Default | Description |
|---|---|---|
| `REDIS_URL` | `redis://localhost:6379` | Event bus + webhook idempotency store. If unavailable, the app runs in a degraded mode (real-time events and idempotency disabled). |
| `OCR_API_URL` | — | OCR service base URL. If empty, OCR features are disabled. |

## Contact info (LINE replies)

`CONTACT_CENTER_NAME`, `CONTACT_PHONE`, `CONTACT_EMAIL`, `CONTACT_EMERGENCY_PHONE`,
`CONTACT_WORKING_HOURS` — surfaced in "contact staff" LINE messages.
