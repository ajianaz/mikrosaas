# MikroSaaS — SaaS Core Engine

[![Go Reference](https://pkg.go.dev/badge/github.com/ajianaz/mikrosaas.svg)](https://pkg.go.dev/github.com/ajianaz/mikrosaas)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajianaz/mikrosaas)](https://goreportcard.com/report/github.com/ajianaz/mikrosaas)

MikroSaaS is a Go module for building multi-tenant SaaS products from scratch. It provides the core building blocks: authentication, user management, tenant isolation, and dynamic RBAC.

## Features

- **Multi-tenant** — Tenant isolation via `tenant_id` column
- **Authentication** — JWT (access + refresh tokens)
- **User Management** — Full CRUD with bcrypt password hashing
- **Dynamic RBAC** — Roles, permissions, and user-role assignment
- **Config-driven** — stdlib only, no Viper, no external config libs
- **PostgreSQL** — UUID v7 PKs, sqlx, hand-written SQL

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.24+ |
| Router | [chi v5](https://github.com/go-chi/chi) |
| Database | PostgreSQL 17/18 |
| DB Driver | sqlx + pgx (stdlib) |
| Auth | JWT (golang-jwt/v5) |
| Config | stdlib (os.Getenv) |

## Quick Start

```go
package main

import (
    "github.com/ajianaz/mikrosaas/pkg"
)

func main() {
    app, err := pkg.NewApp(pkg.Config{
        DatabaseURL: "postgres://user:pass@localhost:5432/mikrosaas?sslmode=disable",
        JWTSecret:   "your-secret-key",
    })
    if err != nil {
        panic(err)
    }
    app.Start()
}
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | — | Secret key for JWT signing |
| `SERVER_HOST` | ❌ | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | ❌ | `8080` | Server bind port |
| `APP_ENV` | ❌ | `development` | `development`, `staging`, `production` |
| `BCRYPT_COST` | ❌ | `12` | Bcrypt hashing cost (4-31) |

## Project Structure

```
mikrosaas/
├── pkg/                    → PUBLIC API
│   ├── mikrosaas.go        → Package declaration + version
│   ├── app.go              → App struct, NewApp(), Start(), Shutdown()
│   └── config.go           → Config struct with defaults
├── internal/               → PRIVATE
│   ├── config/             → Environment config loader
│   ├── auth/               → JWT, login, refresh
│   ├── user/               → User CRUD
│   ├── tenant/             → Tenant management
│   ├── rbac/               → Dynamic RBAC
│   └── database/           → PG connection pool
├── migrations/core/        → Core engine migrations
├── cmd/server/main.go      → Standalone server entry point
├── Makefile
└── .golangci.yml
```

## Development

```bash
# Install dependencies
make tidy

# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Run server (requires DATABASE_URL and JWT_SECRET)
make run
```

## 2-Repo Strategy

MikroSaaS is the engine. Products are built on top:

```
github.com/ajianaz/mikrosaas    ← This repo: core engine
github.com/ajianaz/mikrodius    ← ISP product (rewrite on mikrosaas)
```

## License

Private — All rights reserved.
