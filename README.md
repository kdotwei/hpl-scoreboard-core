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
      - [GET /api/v1/scores/paginated](#get-apiv1scorespaginated)
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
- **Pagination Support**: Efficient cursor-based and offset-based pagination
- **CORS Support**: Configured for frontend integration (localhost:5173, localhost:3000)
- **Database Migration**: Automated schema migrations with golang-migrate
- **Type-safe Database Queries**: Using sqlc for compile-time SQL validation
- **Comprehensive Testing**: Unit tests with testcontainers for integration testing
- **Clean Architecture**: Separation of concerns with handlers, services, and data layers
- **Leaderboard Support**: Retrieve top-performing HPL scores ordered by GFLOPS

## 🏗️ Architecture

The project follows a clean architecture pattern:

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── db/           # Database layer (sqlc generated)
│   ├── handler/      # HTTP handlers (controllers)
│   ├── middleware/   # HTTP middleware (auth)
│   ├── service/      # Business logic layer
│   └── token/        # JWT token management
└── migrations/       # Database migrations
```

**Note**: The config directory is not currently used; configuration is handled directly via environment variables in [main.go](cmd/api/main.go).

## 📋 Prerequisites

- Go 1.24 or higher
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)
- [sqlc](https://sqlc.dev/) (for code generation from SQL)
- Docker (optional, for testcontainers
- Docker (optional, for containerized deployment)

## 🚀 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/kdotwei/hpl-scoreboard-core.git
   cd hpl-scoreboard-core
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run database migrations**
   ```bash
   # Using golang-migrate (install if needed)
   migrate -path migrations -database "postgresql://user:password@localhost:5432/hpl_scoreboard?sslmode=disable" up
   ```

  > [!WARNING]
  > Please run a PostgreSQL service before this work.

4. **Build and run the application**
   ```bash
   go build -o hpl-scoreboard-core cmd/api/main.go
   ./hpl-scoreboard-core
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
|----------|-------------|---------|` |
| `JWT_SECRET_KEY` | JWT signing key (32 characters minimum) | `12345678901234567890123456789012` (development only)alhost:5432/hpl_scoreboard?sslmode=disable` |
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

#### GET /api/v1/scores
Retrieve a list of scores with offset-based pagination (public endpoint).

**Query Parameters:**
- `limit` (optional): Maximum number of scores to return (default: 10)
- `offset` (optional): Number of scores to skip (default: 0)

**Example:**
```
GET /api/v1/scores?limit=20&offset=0
```

**Response:**
```json
[
  {
    "id": "uuid-here",
    "user_id": "username",
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
]
```

#### GET /api/v1/scores/paginated
Retrieve scores with cursor-based pagination for better performance (public endpoint).

**Query Parameters:**
- `limit` (optional): Maximum number of scores to return (1-100, default: 10)
- `offset` (optional): Number of scores to skip (default: 0)

**Example:**
```
GET /api/v1/scores/paginated?limit=50&offset=0
```

**Response:**
```json
{
  "scores": [
    {
      "id": "uuid-here",
      "user_id": "username",
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
  ],
  "total": 1000,
  "limit": 50,
  "offset": 0
}
```
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
 with routes and CORS
├── internal/                   # Private application code
│   ├── db/                     # Database layer (sqlc generated)
│   │   ├── db.go              # Database connection and queries
│   │   ├── models.go          # Generated models
│   │   ├── querier.go         # Generated query interface
│   │   ├── score.sql.go       # Generated score queries
│   │   ├── score_test.go      # Score integration tests
│   │   ├── main_test.go       # Test setup with testcontainers
│   │   └── query/
│   │       └── score.sql      # SQL query definitions
│   ├── handler/               # HTTP request handlers
│   │   ├── handler.go         # Handler struct and constructor
│   │   ├── login.go           # Login endpoint
│   │   ├── login_test.go      # Login handler tests
│   │   ├── score.go           # Score endpoints (Create, List, Paginated)
│   │   ├── score_test.go      # Score handler tests
│   │   └── auth_test.go       # Auth helper tests
│   ├── middleware/            # HTTP middleware
│   │   ├── auth.go            # JWT authentication middleware
│   │   └── auth_test.go       # Middleware tests
│   ├── service/               # Business logic layer
│   │   ├── service.go         # Service interface and constructor
│   │   ├── score.go           # Score business logic
│   │   └── mocks/             # Generated service mocks
│   │       └── Service.go     # Mockery-generated service mock
│   └── token/                 # JWT token management
│       ├── jwt_maker.go       # JWT implementation
│       ├── jwt_maker_test.go  # JWT maker tests
│       ├── maker.go           # Token maker interface
│       ├── payload.go         # JWT payload structure
│       └── mocks/             # Generated token mocks
│           └── Maker.go       # Mockery-generated token maker mock
├── migrations/                # Database migration files
│   ├── 000001_init_schema.up.sql      # Initial schema
│   ├── 000001_init_schema.down.sql
│   ├── 000002_add_hpl_metrics.up.sql  # Add HPL metrics columns
│   └── 000002_add_hpl_metrics.down.sql
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── sqlc.yaml                  # sqlc configuration
├── main_test.go               # Root level test file
├── test_integration.sh        # Integration test script
├── test_pagination.md         # Pagination testing guide
├── PAGINATION_IMPLEMENTATION.md  # Pagination implementation details

```bash
migrate create -ext sql -dir migrations -seq your_migration_name
```

Apply migrations:

```bash
migrate -path migrations -database "your-db-url" up
```

Rollback migrations:

```x] Add public leaderboard endpoints
- [x] Implement pagination (offset-based and cursor-based)
- [x] CORS support for frontend integration
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