# Concurrent Load-Balancing Reverse Proxy

A simple, concurrent reverse proxy in Go with health monitoring and an admin API.

## Features
- Round-Robin & Least-Connections load balancing.
- Periodic background health checking.
- Admin API for dynamic backend management.
- Thread-safe operations.

## Quick Start

1. **Run the Proxy**:
   ```bash
   go run main.go --config=config.json
   ```

2. **Manage Backends (Port 8081)**:
   - **Add**: `curl -X POST http://localhost:8081/backends -d '{"url": "http://host:port"}'`
   - **Remove**: `curl -X DELETE http://localhost:8081/backends -d '{"url": "http://host:port"}'`
   - **Status**: `curl http://localhost:8081/status`

3. **Proxy Traffic (Port 8080)**:
   - Direct your HTTP requests to `http://localhost:8080`.