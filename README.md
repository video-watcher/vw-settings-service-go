# Settings Service Go

Settings Service implementation in Go - migrated from Python FastAPI version.

## Features

- User settings management (theme, language, notifications)
- System settings management (key-value JSONB storage)
- JWT authentication
- Clean Architecture
- PostgreSQL database
- High performance (50%+ faster than Python version)

## Tech Stack

- **Framework:** Fiber v2
- **Database:** PostgreSQL (pgx driver)
- **Auth:** JWT
- **Logging:** zerolog
- **Testing:** testify

## Project Structure

```
settings-service-go/
├── cmd/server/          # Application entry point
├── internal/
│   ├── domain/          # Business entities and interfaces
│   ├── application/     # Use cases and DTOs
│   ├── infrastructure/  # External implementations
│   └── config/          # Configuration
├── tests/               # Tests
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 14+
- Make

### Installation

```bash
# Clone repository
git clone https://github.com/video-watcher/vw-settings-service-go.git
cd vw-settings-service-go

# Install dependencies
go mod download

# Run tests
make test

# Build
make build

# Run
make run
```

### Environment Variables

```bash
PORT=8005
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
JWT_SECRET=your-secret-key
JWT_ALGORITHM=HS256
LOG_LEVEL=info
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/settings/me` | User | Get user settings |
| PATCH | `/api/settings/me` | User | Update user settings |
| GET | `/api/settings/system` | Admin | List system settings |
| GET | `/api/settings/system/{key}` | Admin | Get system setting |
| PUT | `/api/settings/system/{key}` | Admin | Set system setting |
| DELETE | `/api/settings/system/{key}` | Admin | Delete system setting |
| GET | `/health` | None | Health check |

## Development

```bash
# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt

# Run locally
make dev
```

## Migration from Python

This service is a direct replacement for the Python FastAPI version. API compatibility is maintained for seamless migration.

See [Migration Guide](docs/migration.md) for details.

## License

MIT
