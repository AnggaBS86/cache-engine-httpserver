# Cache Engine HTTP Server

A high-performance HTTP caching server built with **Go**, **Fiber v3**, and **Allegro BigCache**.  
This project provides a simple REST API to store and retrieve data from an in-memory cache with configurable expiration.

## Features

- Fast in-memory caching with [BigCache](https://github.com/allegro/bigcache)
- REST API powered by [Fiber](https://gofiber.io)
- Middleware support for request handling
- Configurable cache expiration via `.env`
- Clean project structure (`router`, `middleware`, `model`)

## Getting Started

### Prerequisites
- Go 1.21+
- Git
- Make (optional)

### Installation

```bash
git clone https://github.com/yourusername/cache-engine-httpserver.git
cd cache-engine-httpserver
go mod tidy
```
### Configuration

Create a .env file in the project root:

```bash
DEFAULT_CACHE_DURATION_IN_SECONDS=60
PORT=3000
```

### Build and Run
```bash
go build
./cache-engine-httpserver
```
Server will run on http://localhost:3000.

### API Endpoints
| Method | Endpoint      | Description           |
| ------ | ------------- | --------------------- |
| GET    | `/cache/:key` | Retrieve cached value |
| POST   | `/cache`      | Store data in cache   |
| DELETE | `/cache/:key` | Delete a cached entry |


### Project Structure
```
.
├── main.go
├── go.mod
├── internal/
│   └── api/
│       ├── router/      # Route definitions
│       ├── model/       # API models
│       └── middleware/  # Middlewares
└── .env                 # Environment variables
```

### Run Test
```bash
go test ./...
```
