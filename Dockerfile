# ── Stage 1: Build ──────────────────────────────────────────
FROM golang:1.24-alpine AS builder

ARG VERSION=dev

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /build/server ./cmd/app

# ── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /build/server .

# Run as non-root user
USER appuser

EXPOSE 3000

# Health check using the liveness endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:3000/health/live || exit 1

CMD ["/app/server"]
