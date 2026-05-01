# Quote Schedule System

REST API built with Go + Gin + GORM + MySQL, with a React/MUI frontend.

## Requirements

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)

## Setup

**1. Start MySQL**

```bash
docker compose up -d
```

**2. Configure environment**

```bash
cp br-api/.env.example br-api/.env.local
```

**3. Run the API**

```bash
cd br-api
go run ./cmd/api
```

API will be available at `http://localhost:8080`.

## Health check

```bash
curl http://localhost:8080/api/v1/health
```
