# Inventory Manager Server

This is the backend API for the Inventory Manager application, built with **Go** and the **Gin Gonic** web framework. It provides a robust and scalable architecture to manage inventories and products.

## 🚀 Overview

The server follows a clean architecture pattern, separating concerns into distinct layers:
- **Handlers**: Handle HTTP requests and responses.
- **Services**: Contain business logic.
- **Repositories**: Manage data persistence and database interactions.
- **Models**: Define the data structures for the application.

## 📂 Folder Structure

```text
server/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── config/          # Database connection and configuration
│   ├── handler/         # HTTP request handlers (Controllers)
│   ├── models/          # Data structures (Inventory, Product)
│   ├── repository/      # Data access layer (SQL queries)
│   └── service/         # Business logic layer
├── routes/
│   └── routes.go        # API route definitions
├── .env                 # Environment variables
├── docker-compose.yml   # Infrastructure setup (PostgreSQL)
├── go.mod               # Go module dependencies
└── go.sum               # Dependency checksums
```

## 🛠️ Tech Stack

- **Language:** [Go (Golang)](https://golang.org/)
- **Framework:** [Gin Gonic](https://gin-gonic.com/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Drivers:** `pgx` for PostgreSQL integration
- **Configuration:** `godotenv` for environment variable management

## ⚙️ Setup & Installation

### Prerequisites
- [Go 1.25+](https://go.dev/doc/install)
- [Docker](https://www.docker.com/get-started) (for running the database)

### 1. Start the Database
The server requires a PostgreSQL database. You can start it easily using Docker Compose:
```bash
docker-compose up -d
```
*This will start a Postgres instance on port `5432` with the credentials defined in `docker-compose.yml`.*

### 2. Configure Environment Variables
Ensure you have a `.env` file in the `server/` directory with the following content:
```env
PORT=8080
DB_URL=postgres://inventory_user:inventory123@127.0.0.1:5432/inventory?sslmode=disable
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run the Server
```bash
go run cmd/main.go
```
The server will start on `http://localhost:8080`.

## 📡 API Endpoints

- `GET /ping`: Health check endpoint.
- `GET /api/v1/inventories`: List all inventories.
- `GET /api/v1/products`: List all products.
- *(Other CRUD endpoints are managed within the routes layer)*
