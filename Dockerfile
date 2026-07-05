# ── Stage 1: Build ──────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile a small, static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/server ./cmd/app

# ── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:3.20

# ca-certificates for TLS (LINE/OCR/Redis), tzdata for correct local time.
# Create an unprivileged user to run the service.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=builder /build/server .

# Run as non-root
USER appuser

EXPOSE 3000

# Liveness probe via the /livez endpoint (busybox wget ships with alpine).
# Adjust the port here if PORT is overridden at runtime.
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:3000/livez || exit 1

CMD ["/app/server"]
