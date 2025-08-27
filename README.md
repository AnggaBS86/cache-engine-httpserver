# Cache Engine HTTP Server

A high-performance HTTP caching server built with **Go**, **Fiber v3**, and **Allegro BigCache**.  
This project provides a simple REST API to store and retrieve data from an in-memory cache with configurable expiration.

## Features

- Fast in-memory caching with [BigCache](https://github.com/allegro/bigcache)
- REST API powered by [Fiber](https://gofiber.io)
- Middleware support for request handling
- Configurable cache expiration via `.env`
- Clean project structure (`router`, `middleware`, `model`)

## Why Cache Engine HTTP Server?

Compared to general caching solutions like **Memcached** or **Redis**, this project has several advantages for certain use cases:

### Pros

1. **Simplicity & Lightweight**
   - Runs as a single Go binary (`./cache-engine-httpserver`).
   - No external dependencies or separate server processes.
   - Easy to deploy in small/microservice environments.

2. **Performance for Local In-Memory Cache**
   - Uses [BigCache](https://github.com/allegro/bigcache), optimized for **fast concurrent access** with **low GC overhead**.
   - Often faster than Redis/Memcached for *in-process* or *single-node* caching.

3. **Low Latency**
   - Pure in-memory storage → no network round-trips required (unless exposed via HTTP).
   - Ideal for high-throughput microservices.

4. **Customizable via Go**
   - Extendable with custom middleware, authentication, or monitoring directly in Go.
   - Unlike Redis/Memcached, you don’t need external scripting or plugins.

5. **Great for Ephemeral Caching**
   - Perfect for temporary cache-only data (e.g., sessions, computed results, API responses).
   - BigCache reduces Go heap usage → handles millions of entries without heavy GC pressure.

6. **Developer-Friendly**
   - Cache and API logic live in the same codebase.
   - Easier local development/testing: no Redis/Memcached container required.

### Limitations

To be transparent, this project is **not a full replacement** for Redis/Memcached:

- ❌ No persistence → cache is lost on restart.  
- ❌ No clustering or distributed caching.  
- ❌ Only supports key/value cache → no advanced data types (lists, sets, streams, pub/sub, etc).  

---

In summary:  
This project is **lightweight, fast, and easy** for **single-node, in-memory caching** inside microservices.  
Redis/Memcached remain the better choice for **distributed caching, persistence, or advanced features**.


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
