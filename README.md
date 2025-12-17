# HPL Scoreboard Core

A high-performance computing (HPC) scoreboard API service for tracking and managing High-Performance Linpack (HPL) benchmark results. This service provides RESTful APIs for submitting HPL benchmark scores and retrieving leaderboards.

## 📋 Table of Contents

- [HPL Scoreboard Core](#hpl-scoreboard-core)
  - [📋 Table of Contents](#-table-of-contents)
  - [🔍 Overview](#-overview)
  - [✨ Features](#-features)
  - [🏗️ Architecture](#️-architecture)
  - [📋 Prerequisites](#-prerequisites)
  - [🚀 Installation](#-installation)
  - [⚙️ Configuration](#️-configuration)
    - [Environment Variables](#environment-variables)
  - [🔌 API Endpoints](#-api-endpoints)
    - [Authentication](#authentication)
      - [POST /api/v1/login](#post-apiv1login)
    - [Scores](#scores)
      - [POST /api/v1/scores](#post-apiv1scores)
  - [🗄️ Database Schema](#️-database-schema)
    - [Scores Table](#scores-table)
  - [🛠️ Development](#️-development)
    - [Code Generation](#code-generation)
    - [Running Tests](#running-tests)
    - [Linting](#linting)
    - [Database Migrations](#database-migrations)
  - [🧪 Testing](#-testing)
  - [📁 Project Structure](#-project-structure)
  - [🤝 Contributing](#-contributing)
    - [Development Guidelines](#development-guidelines)
  - [📄 License](#-license)
  - [🚧 Roadmap](#-roadmap)
  - [📧 Support](#-support)

## 🔍 Overview

HPL Scoreboard Core is a Go-based REST API service designed to collect, store, and display HPL (High-Performance Linpack) benchmark results. It provides authenticated endpoints for submitting benchmark scores and public endpoints for viewing leaderboards.

The service is built with modern Go practices, using PostgreSQL for data persistence and JWT for authentication.

## ✨ Features

- **JWT-based Authentication**: Secure API access with JSON Web Tokens
- **HPL Score Management**: Submit and retrieve HPL benchmark results
- **Database Migration**: Automated schema migrations with golang-migrate
- **Type-safe Database Queries**: Using sqlc for compile-time SQL validation
- **Comprehensive Testing**: Unit tests with testcontainers for integration testing
- **Clean Architecture**: Separation of concerns with handlers, services, and data layers
- **Docker Support**: Easy deployment with containerization
- **Leaderboard Support**: Retrieve top-performing HPL scores

## 🏗️ Architecture

The project follows a clean architecture pattern:

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── db/           # Database layer (sqlc generated)
│   ├── handler/      # HTTP handlers (controllers)
│   ├── middleware/   # HTTP middleware (auth, logging)
│   ├── service/      # Business logic layer
│   └── token/        # JWT token management
└── migrations/       # Database migrations
```

## 📋 Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make (optional, for running Makefile commands)
- Docker (optional, for containerized deployment)

## 🚀 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/kdotwei/hpl-scoreboard.git
   cd hpl-scoreboard-core
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb hpl_scoreboard
   ```

4. **Run database migrations**
   ```bash
   # Using golang-migrate (install if needed)
   migrate -path migrations -database "postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable" up
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

## ⚙️ Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory:

```env
# Database Configuration
DB_SOURCE=postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable

# Server Configuration
SERVER_ADDRESS=:8080

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-here-32-chars
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_SOURCE` | PostgreSQL connection string | `postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable` |
| `SERVER_ADDRESS` | Server listen address | `:8080` |
| `JWT_SECRET_KEY` | JWT signing key (32 characters minimum) | Development key |

## 🔌 API Endpoints

### Authentication

#### POST /api/v1/login
Login and receive JWT token for authenticated endpoints.

**Request:**
```json
{
  "username": "your-username"
}
```

**Response:**
```json
{
  "access_token": "jwt-token-here",
  "user": {
    "username": "your-username"
  }
}
```

### Scores

#### POST /api/v1/scores
Submit a new HPL benchmark score (requires authentication).

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request:**
```json
{
  "gflops": 1234.56,
  "problem_size_n": 50000,
  "block_size_nb": 256,
  "linux_username": "hpc-user",
  "n": 50000,
  "nb": 256,
  "p": 4,
  "q": 4,
  "execution_time": 1800.5
}
```

**Response:**
```json
{
  "id": "uuid-here",
  "user_id": "your-username",
  "gflops": 1234.56,
  "problem_size_n": 50000,
  "block_size_nb": 256,
  "linux_username": "hpc-user",
  "n": 50000,
  "nb": 256,
  "p": 4,
  "q": 4,
  "execution_time": 1800.5,
  "submitted_at": "2024-12-18T10:00:00Z"
}
```

## 🗄️ Database Schema

### Scores Table

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key (auto-generated) |
| `user_id` | VARCHAR | User identifier |
| `gflops` | DOUBLE PRECISION | Performance in GFLOPS |
| `problem_size_n` | INT | Problem size N |
| `block_size_nb` | INT | Block size NB |
| `linux_username` | VARCHAR | System username |
| `n` | INT | Matrix dimension N |
| `nb` | INT | Block size |
| `p` | INT | Process grid P dimension |
| `q` | INT | Process grid Q dimension |
| `execution_time` | DOUBLE PRECISION | Execution time in seconds |
| `submitted_at` | TIMESTAMPTZ | Submission timestamp |

## 🛠️ Development

### Code Generation

The project uses sqlc for type-safe database queries:

```bash
# Generate database code from SQL queries
sqlc generate
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests with testcontainers
go test ./internal/db/...
```

### Linting

```bash
# Run golangci-lint
golangci-lint run
```

### Database Migrations

Create a new migration:

```bash
migrate create -ext sql -dir migrations -seq your_migration_name
```

Apply migrations:

```bash
migrate -path migrations -database "your-db-url" up
```

Rollback migrations:

```bash
migrate -path migrations -database "your-db-url" down
```

## 🧪 Testing

The project includes comprehensive testing:

- **Unit Tests**: For handlers, services, and utilities
- **Integration Tests**: Using testcontainers for database testing
- **Mocks**: Generated mocks for service interfaces

Key testing features:
- PostgreSQL integration tests with testcontainers
- JWT token testing
- HTTP handler testing
- Service layer testing

## 📁 Project Structure

```
.
├── cmd/api/                    # Application entry point
│   └── main.go                 # Main application file
├── internal/                   # Private application code
│   ├── config/                 # Configuration management
│   ├── db/                     # Database layer
│   │   ├── db.go              # Database connection and queries
│   │   ├── models.go          # Generated models
│   │   ├── querier.go         # Generated query interface
│   │   ├── score.sql.go       # Generated score queries
│   │   └── query/
│   │       └── score.sql      # SQL query definitions
│   ├── handler/               # HTTP request handlers
│   │   ├── handler.go         # Handler struct and constructor
│   │   ├── login.go           # Login endpoint
│   │   └── score.go           # Score endpoints
│   ├── middleware/            # HTTP middleware
│   │   └── auth.go            # JWT authentication middleware
│   ├── service/               # Business logic layer
│   │   ├── service.go         # Service interface and constructor
│   │   ├── score.go           # Score business logic
│   │   └── mocks/             # Generated service mocks
│   └── token/                 # JWT token management
│       ├── jwt_maker.go       # JWT implementation
│       ├── maker.go           # Token maker interface
│       ├── payload.go         # JWT payload structure
│       └── mocks/             # Generated token mocks
├── migrations/                # Database migration files
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   ├── 000002_add_hpl_metrics.up.sql
│   └── 000002_add_hpl_metrics.down.sql
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── sqlc.yaml                  # sqlc configuration
└── README.md                  # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation as needed
- Run linting and tests before submitting PRs
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🚧 Roadmap

- [ ] Add public leaderboard endpoints
- [ ] Implement score filtering and sorting
- [ ] Add metrics and monitoring
- [ ] Docker compose setup
- [ ] API rate limiting
- [ ] User management system
- [ ] Score validation and verification
- [ ] Performance benchmarking dashboard

## 📧 Support

For questions or support, please open an issue in the GitHub repository.