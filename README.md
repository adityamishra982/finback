# Finback

Finback is a robust RESTful API for managing financial records and generating summaries. Built with Go, the Gin web framework, and PostgreSQL via GORM, it features secure JWT-based authentication and Role-Based Access Control (RBAC).

## Features

- **User Authentication**: Secure JWT-based login and registration.
- **Role-Based Access Control**: Three built-in user roles defining access permissions:
  - `admin`: Can create, update, and delete financial records, as well as manage users.
  - `analyst`: Can view financial records and dashboard summaries.
  - `viewer`: Can view records and dashboard summaries (read-only access).
- **Financial Records Management**: Track income and expenses categorized by types and dates.
- **Dashboard Summary**: Aggregate viewing for quick financial overviews.
- **PostgreSQL Database integration**: Uses GORM for efficient object-relational mapping and database migrations.

## Tech Stack

- [Go](https://go.dev/) (1.26.1+)
- [Gin Web Framework](https://github.com/gin-gonic/gin) for routing
- [GORM](https://gorm.io/) with PostgreSQL driver
- [Docker & Docker Compose](https://www.docker.com/) for local development databases (PostgreSQL, Redis)
- [golang-jwt](https://github.com/golang-jwt/jwt) for JWT parsing and generation

## Project Structure

```text
├── cmd
│   └── api               # Application entrypoint (main.go)
├── internal
│   ├── api               # Route handlers (auth, admin, dashboard, records)
│   ├── config            # Environment configuration 
│   ├── db                # Database connection setup
│   ├── errors            # Custom error responses
│   ├── middleware        # Gin middlewares (Auth, Role check, Error handler)
│   └── models            # GORM Data models (User, Record)
├── pkg
│   └── utils             # Shared utilities (JWT, password hashing)
├── docker-compose.yml    # Database infrastructure setup
├── go.mod                # Go module dependencies
└── build.cmd / build.sh  # Build scripts build automation
```

## Getting Started

### Prerequisites

- Go 1.26 or higher
- Docker and Docker Compose (for running PostgreSQL and Redis locally)

### Installation & Setup

1. **Clone the repository** (if you haven't already navigated to it)
2. **Start the Database**
   Using Docker Compose, start the PostgreSQL and Redis containers in the background:
   ```bash
   docker-compose up -d
   ```
3. **Configure Environment Variables**
   Create a `.env` file in the root directory. At a minimum, set database connection info and JWT secrets:
   ```env
   PORT=8080
   GIN_MODE=debug
   # Add your Postgres connection strings, JWT secret, etc.
   ```
4. **Build and Run**
   You can run the service directly:
   ```bash
   go run cmd/api/main.go
   ```
   Or use the provided build scripts:
   ```bash
   ./build.sh    # macOS/Linux
   # or
   .\build.cmd   # Windows
   ```

## API Routes Overview

Base Path: `/api/v1`

### Authentication (`/auth`)
- `POST /register`: Register a new user
- `POST /login`: Authenticate a user and receive a JWT

### Financial Records (`/records`) - Requires Auth
- `GET /`: List financial records (All roles)
- `POST /`: Create a new record (Admin only)
- `PUT /:id`: Update an existing record (Admin only)
- `DELETE /:id`: Delete a record (Admin only)

### Dashboard (`/dashboard`) - Requires Auth
- `GET /summary`: View financial summary and aggregates (All roles)

### Admin (`/admin`) - Requires Admin Auth
- `GET /users`: List users 
- *(Other user management endpoints configured in `users.go`)*

## License

This project is licensed under the MIT License.
