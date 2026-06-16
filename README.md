# dauth

Standalone authentication microservice providing centralized identity management for multiple products. One global user identity spans all products. Other services verify JWTs locally using dauth's public key — no per-request round-trip required.

**Stack:** Go · PostgreSQL · chi · RS256 JWT

---

## Architecture

Layered service with strict import boundaries:

```
dauth/
├── cmd/server/          # main.go — wires everything, starts HTTP server
├── internal/
│   ├── core/            # domain models + error types (no framework deps)
│   ├── repository/      # Postgres queries via pgx/v5
│   ├── service/         # business logic: register, login, refresh, revoke, invite
│   └── transport/
│       └── http/        # chi-based REST handlers, middleware, routes
├── pkg/
│   └── jwtutil/         # RS256 sign/verify + JWKS — importable by other services
├── migrations/          # SQL migration files (golang-migrate)
├── config/              # env-based config loading
├── keys/                # RSA key pair (PEM)
├── Dockerfile
└── docker-compose.yml
```

**Import rules:** `core` → nothing internal. `repository` → `core`. `service` → `core` + `repository`. `transport/http` → `service` only.

---

## API

### Auth
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/register` | Create user + assign to product, return tokens |
| `POST` | `/api/v1/auth/login` | email + password → access token + refresh token |
| `POST` | `/api/v1/auth/refresh` | Rotate refresh token → new token pair |
| `POST` | `/api/v1/auth/logout` | Revoke refresh token |
| `GET`  | `/api/v1/auth/me` | Current user info (requires Bearer token) |

### Invitations (admin only)
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/invitations` | Generate invite token, return token + expiry |
| `GET`  | `/api/v1/invitations/:token` | Validate token, return email + product |
| `POST` | `/api/v1/invitations/:token/accept` | Complete registration via invite |

### Public Key
| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/api/v1/.well-known/jwks.json` | RS256 public key for JWT verification |

### Health
| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/health` | Liveness check |

---

## Token Design

**Access token:** RS256 JWT, 15 min TTL. Claims: `sub` (user_id), `product_id`, `role`, `exp`, `iat`.

**Refresh token:** `crypto/rand` 32 bytes (opaque), SHA-256 hashed in DB, 30-day TTL. Rotated on every use.

**Invitation token:** `crypto/rand` 32 bytes, SHA-256 hashed in DB, 7-day TTL. Single-use.

---

## Configuration

| Env var | Default | Description |
|---------|---------|-------------|
| `DATABASE_URL` | _(required)_ | Postgres connection string |
| `PRIVATE_KEY_PATH` | _(required)_ | Path to RSA private key PEM |
| `PUBLIC_KEY_PATH` | _(required)_ | Path to RSA public key PEM |
| `PORT` | `8080` | HTTP listen port |
| `ACCESS_TOKEN_TTL` | `15m` | Access token lifetime |
| `REFRESH_TOKEN_TTL` | `720h` | Refresh token lifetime |
| `INVITATION_TTL` | `168h` | Invitation token lifetime |

---

## Development

```sh
make build          # compile → ./bin/dauth
make run            # go run ./cmd/server
make test           # go test ./... -race
make test/cover     # coverage report (HTML)
make lint           # golangci-lint
make migrate-up     # apply migrations (requires DATABASE_URL)
make migrate-down   # rollback migrations
```

---

## Consuming dauth in Other Services

**Go services:** import `pkg/jwtutil` directly as a Go module and call `Manager.Verify()`.

**Any language:** fetch the public key from `GET /api/v1/.well-known/jwks.json` and verify JWTs locally using any RS256-capable JWT library.

No per-request call to dauth needed for token verification.

---

## Implementation Status

| Layer | Status |
|-------|--------|
| `internal/core` — models + errors | ✅ Done |
| `pkg/jwtutil` — RS256 sign/verify/JWKS | ✅ Done |
| `config` — env loading | ✅ Done |
| `internal/repository` — Postgres queries | 🔲 Planned |
| `internal/service` — business logic | 🔲 Planned |
| `internal/transport/http` — REST handlers | 🔲 Planned |
| `cmd/server` — entry point | 🔲 Planned |
| `migrations` — SQL schema | 🔲 Planned |
| `Dockerfile` + `docker-compose.yml` | 🔲 Planned |
