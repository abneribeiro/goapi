# Equipment Rental API

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1-6BA539?style=flat&logo=openapi-initiative)](https://www.openapis.org/)


> A production-ready RESTful API for equipment rental management built with Go's standard library.

<div align="center">

   https://github.com/user-attachments/assets/14894f24-285b-46a4-9c32-a9868e184730

</div>

This API enables users to list equipment for rent, make reservations, and manage the complete rental lifecycle.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [API Documentation](#api-documentation)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Database Schema](#database-schema)
- [Configuration](#configuration)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Management**: Registration and authentication with JWT tokens
- **Equipment Catalog**: Full CRUD operations with photo uploads, search, and filtering
- **Reservation System**: Complete workflow with approval, rejection, cancellation, and completion
- **Notification System**: Real-time notifications for reservation updates
- **API Documentation**: Interactive Scalar UI with OpenAPI 3.1 specification
- **Pagination**: Built-in pagination support for list endpoints
- **CORS Support**: Configurable cross-origin resource sharing
- **Structured Logging**: JSON-formatted logs with request tracing
- **Graceful Shutdown**: Proper handling of server shutdown signals

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer                              │
│  ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌─────────────┐  │
│  │  CORS   │→ │  Logger  │→ │ Recovery │→ │    Auth     │  │
│  └─────────┘  └──────────┘  └──────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Handler Layer                            │
│  ┌────────┐ ┌──────┐ ┌───────────┐ ┌─────────────┐ ┌─────┐ │
│  │  Auth  │ │ User │ │ Equipment │ │ Reservation │ │Notif│ │
│  └────────┘ └──────┘ └───────────┘ └─────────────┘ └─────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer                            │
│         (Business Logic, Validation, Notifications)          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                           │
│              (Data Access, SQL Queries)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     PostgreSQL                               │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.24** | Primary language using standard library `net/http` |
| **PostgreSQL 16** | Relational database for data persistence |
| **JWT (HS256)** | Stateless authentication |
| **Docker** | Containerization and deployment |
| **Scalar** | Interactive API documentation UI |
| **OpenAPI 3.1** | API specification format |

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (recommended)
- PostgreSQL 16 (if running locally without Docker)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/abneribeiro/goapi.git
cd goapi
```

2. Install dependencies:
```bash
go mod download
```

### Running with Docker

The easiest way to run the application is using Docker Compose:

```bash
# Start all services (API + PostgreSQL)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

The API will be available at `http://localhost:8080`.

### Running Locally

1. Start PostgreSQL (using Docker):
```bash
docker-compose up -d postgres
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=equipment_rental
export JWT_SECRET=your-secure-secret-key
```

3. Run the application from the project root:
```bash
# Make sure you're in the project root directory
cd /path/to/goapi
go run ./cmd/api
```

**Important**: Always run the server from the project root directory, not from `cmd/api`, to ensure proper path resolution for static files (docs, uploads).

4. (Optional) Seed the database with sample data:
```bash
go run ./scripts/seed/main.go
```

This creates test users:
| Email | Password | Role |
|-------|----------|------|
| owner@example.com | Password123 | owner |
| renter@example.com | Password123 | renter |
| owner2@example.com | Password123 | owner |

## API Documentation

Interactive API documentation is available at:

- **Scalar UI**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/docs/openapi.yaml

The documentation includes:
- Complete endpoint descriptions
- Request/response examples
- Authentication details
- Try-it-out functionality

## API Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | API health status |

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/users/me` | Required | Get current user profile |
| PUT | `/api/v1/users/me` | Required | Update current user profile |

### Equipment

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/equipment` | - | List equipment (paginated, filterable) |
| GET | `/api/v1/equipment/search` | - | Search equipment |
| GET | `/api/v1/equipment/categories` | - | Get available categories |
| GET | `/api/v1/equipment/{id}` | - | Get equipment by ID |
| GET | `/api/v1/equipment/{id}/availability` | - | Get availability calendar |
| POST | `/api/v1/equipment` | Required | Create equipment |
| PUT | `/api/v1/equipment/{id}` | Required | Update equipment |
| DELETE | `/api/v1/equipment/{id}` | Required | Delete equipment |
| POST | `/api/v1/equipment/{id}/photos` | Required | Upload equipment photo |

### Reservations

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/reservations` | Required | List my reservations |
| GET | `/api/v1/reservations/owner` | Required | List reservations for my equipment |
| GET | `/api/v1/reservations/{id}` | Required | Get reservation by ID |
| POST | `/api/v1/reservations` | Required | Create reservation |
| PUT | `/api/v1/reservations/{id}/approve` | Required | Approve reservation (owner) |
| PUT | `/api/v1/reservations/{id}/reject` | Required | Reject reservation (owner) |
| PUT | `/api/v1/reservations/{id}/cancel` | Required | Cancel reservation |
| PUT | `/api/v1/reservations/{id}/complete` | Required | Complete reservation (owner) |

### Notifications

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/notifications` | Required | List notifications |
| GET | `/api/v1/notifications/unread-count` | Required | Get unread count |
| PUT | `/api/v1/notifications/{id}/read` | Required | Mark as read |
| PUT | `/api/v1/notifications/read-all` | Required | Mark all as read |
| DELETE | `/api/v1/notifications/{id}` | Required | Delete notification |

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. After successful login, include the token in the `Authorization` header:

```
Authorization: Bearer <your_token>
```

**Token expiration**: 24 hours (configurable via `JWT_EXPIRATION_HOURS`)

### Example: Login and Use Token

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "owner@example.com", "password": "Password123"}' \
  | jq -r '.data.token')

# Use token for authenticated requests
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

## Database Schema

```
┌─────────────────┐     ┌─────────────────────┐
│     users       │     │     equipment       │
├─────────────────┤     ├─────────────────────┤
│ id (PK)         │◄────│ owner_id (FK)       │
│ email (UNIQUE)  │     │ id (PK)             │
│ password_hash   │     │ name                │
│ name            │     │ description         │
│ phone           │     │ category            │
│ role            │     │ price_per_hour      │
│ verified        │     │ price_per_day       │
│ created_at      │     │ price_per_week      │
│ updated_at      │     │ location            │
└─────────────────┘     │ latitude/longitude  │
        │               │ available           │
        │               │ auto_approve        │
        │               └─────────────────────┘
        │                        │
        ▼                        ▼
┌─────────────────────────────────────────────┐
│              reservations                    │
├─────────────────────────────────────────────┤
│ id (PK)                                      │
│ equipment_id (FK) ─────────────────────────►│
│ renter_id (FK) ─────────────────────────────│
│ start_date, end_date                         │
│ status (pending/approved/rejected/...)       │
│ total_price                                  │
│ cancellation_reason                          │
└─────────────────────────────────────────────┘
```

## Configuration

All configuration is done via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host address | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `equipment_rental` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_EXPIRATION_HOURS` | Token expiration (hours) | `24` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `debug` |
| `UPLOAD_PATH` | File upload directory | `./uploads` |
| `DOCS_PATH` | Documentation files directory | `./docs` |
| `MAX_FILE_SIZE` | Max upload size (bytes) | `10485760` |

You can also create a `.env` file in the project root for local development.

## Testing

Run all tests:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Run specific package tests:
```bash
go test ./internal/handler/... -v
go test ./internal/middleware/... -v
```

## Project Structure

```
goapi/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── database/                # Database connection & migrations
│   │   ├── postgres.go
│   │   └── migrations.go
│   ├── handler/                 # HTTP handlers (controllers)
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── equipment.go
│   │   ├── reservation.go
│   │   ├── notification.go
│   │   └── docs.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── model/                   # Data models & DTOs
│   │   ├── user.go
│   │   ├── equipment.go
│   │   ├── reservation.go
│   │   ├── notification.go
│   │   └── response.go
│   ├── pkg/                     # Internal packages
│   │   ├── jwt/                 # JWT utilities
│   │   ├── logger/              # Structured logging
│   │   ├── pagination/          # Pagination helpers
│   │   └── validator/           # Input validation
│   ├── repository/              # Data access layer
│   │   ├── user.go
│   │   ├── equipment.go
│   │   ├── reservation.go
│   │   └── notification.go
│   ├── router/                  # Route configuration
│   │   └── router.go
│   └── service/                 # Business logic layer
│       ├── auth.go
│       ├── user.go
│       ├── equipment.go
│       ├── reservation.go
│       └── notification.go
├── docs/
│   └── openapi.yaml             # OpenAPI 3.1 specification
├── scripts/
│   └── seed/                    # Database seeding
│       └── main.go
├── requests/                    # HTTP request test files
├── uploads/                     # File upload directory
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Reservation Workflow

```
                    ┌─────────────┐
                    │   PENDING   │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
       ┌──────────┐  ┌──────────┐  ┌──────────┐
       │ APPROVED │  │ REJECTED │  │CANCELLED │
       └────┬─────┘  └──────────┘  └──────────┘
            │
            ▼
      ┌───────────┐
      │ COMPLETED │
      └───────────┘
```

## API Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Paginated Response:**
```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Built with Go by [Abner Ribeiro](https://github.com/abneribeiro)
