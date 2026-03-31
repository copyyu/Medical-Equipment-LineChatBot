# ── Stage 1: Build ──────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/server ./cmd/app

# ── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /build/server .

EXPOSE 3000

CMD ["/app/server"]
