.PHONY: run test lint build migrate-up migrate-down migrate-version migrate-create \
       docker-up docker-down generate-mocks coverage

# ── Development ──────────────────────────────────────────────
run:
	go run ./cmd/app

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/app

# ── Database Migration ───────────────────────────────────────
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-down-all:
	go run ./cmd/migrate down --all

migrate-version:
	go run ./cmd/migrate version

migrate-force:
	go run ./cmd/migrate force $(version)

# ── Testing ──────────────────────────────────────────────────
test:
	go test -race -v -count=1 ./...

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML report: go tool cover -html=coverage.out"

# ── Linting & Security ──────────────────────────────────────
lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Run 'gofmt -w .' to fix formatting" && gofmt -l . && exit 1)

vuln:
	govulncheck ./...

# ── Code Generation ─────────────────────────────────────────
generate-mocks:
	go run github.com/vektra/mockery/v2@latest

# ── Docker ───────────────────────────────────────────────────
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f app

# ── All Checks (CI equivalent) ──────────────────────────────
check: fmt-check vet lint test
	@echo "✅ All checks passed"
