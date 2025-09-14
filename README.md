# grpc-balance-processor
Short guide to run the project, verify gRPC endpoints, and hack locally.

## Prerequisites
- Docker & Docker Compose
- `grpcurl` (for local testing)
- (Optional) Go 1.24+ if you want to run outside Docker

## 1) Clone & prepare env
```bash
git clone https://github.com/MaksimPozharskiy/grpc-balance-processor.git
cd grpc-balance-processor
cp .env.example .env
```

## 2) Run with Docker
```bash
docker compose up -d --build
docker compose ps
```

## Local development (optional)
Build locally:
```bash
go build ./cmd/app
```

Run tests:
```bash
go test ./...
```

Lint (golangci-lint):
```bash
golangci-lint run ./...
```
