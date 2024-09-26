# Binance Candle Crawler Service

This project is a Go-based service for crawling Binance candle data.

## Prerequisites

- Go (version not specified, use the latest stable version)
- Docker
- PostgreSQL
- `migrate` CLI tool for database migrations
- `swag` for Swagger documentation generation

## Setup

1. Environment Variables:
   Ensure you have an `app.env` file in the project root with the following variables set:
   - POSTGRES_USER
   - POSTGRES_PASSWORD
   - POSTGRES_HOST
   - POSTGRES_PORT
   - POSTGRES_DB
   - SSL_MODE
   - DB_TIMEZONE

2. Install dependencies:
   ```bash
   make install
   make deps
   ```

3. Set up the database:
   ```bash
   make network
   make postgres
   make startdb
   make migration-up
   ```

4. Generate Swagger documentation:
   ```bash
   make docs
   ```

## Development

Run the service locally:
```bash
make dev
```

## Building

Build the project:
```bash
make build
```

Build for multiple platforms:
```bash
make build-all
```

## Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

## Linting

Run the linter:
```bash
make lint
```

## Database Migrations

Create a new migration:
```bash
make migration-new name=your_migration_name
```

Apply migrations:
```bash
make migration-up
```

Revert migrations:
```bash
make migration-down
```

## Cleaning

Clean build artifacts:
```bash
make clean
```

## Dependency Management

Update dependencies:
```bash
make update-deps
```

## Docker

Create a Docker network:
```bash
make network
```

Create and run a PostgreSQL container:
```bash
make postgres
```

Start the PostgreSQL container:
```bash
make startdb
```

## Additional Commands

- `make run`: Build and run the project
- `make lint`: Run the linter
- `make docs`: Generate Swagger documentation

## Notes

- The main application is located in `./cmd/server`
- Migrations are stored in `internal/initializers/migrations`
- The project uses a custom Docker network named `beego_network`

For more details on each command, refer to the Makefile in the project root.