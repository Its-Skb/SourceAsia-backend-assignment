# Backend Assignment - Source Asia

This repository contains the implementation for the Source Asia Backend Assignment.

The project is built using Go and Gin framework and includes:

1. Part 1 - Rate Limited API
2. Part 2 - Product Catalog API (to be implemented)

---

# Tech Stack

- Go
- Gin Framework
- In-memory storage using maps and slices
- sync.RWMutex for concurrency safety

---

# Project Structure

```text
backend-assignment/
│
├── cmd/server
├── internal/
│   ├── handlers
│   ├── limiter
│   ├── middleware
│   ├── models
│   ├── routes
│   └── storage
│
├── go.mod
├── go.sum
└── README.md
```

---

# How to Run

## Install dependencies

```bash
go mod tidy
```

## Run the server

```bash
go run ./cmd/server
```

Server runs on:

```text
http://localhost:8080
```

---

# Part 1 - Rate Limited API

## Rate Limiting Design

This implementation uses a rolling/sliding 1-minute window.

Each user is allowed:

```text
Maximum 5 accepted requests per minute
```

The limiter is concurrency-safe using:

```go
sync.RWMutex
```

This ensures parallel requests for the same user_id cannot bypass the limit.

---

# API Endpoints

---

## POST /request

Accepts a request payload for a user.

### Request Body

```json
{
  "user_id": "user1",
  "payload": {
    "message": "hello"
  }
}
```

### Success Response

Status:

```text
201 Created
```

Response:

```json
{
  "success": true,
  "message": "request accepted"
}
```

---

### Validation Errors

Status:

```text
400 Bad Request
```

Example:

```json
{
  "success": false,
  "error": "user_id is required"
}
```

---

### Rate Limit Exceeded

Status:

```text
429 Too Many Requests
```

Response:

```json
{
  "success": false,
  "error": "rate limit exceeded"
}
```

---

## GET /stats

Returns per-user statistics and global totals.

### Response

```json
{
  "success": true,
  "message": "stats fetched successfully",
  "data": {
    "users": {
      "user1": {
        "accepted_requests_current_window": 1,
        "rejected_requests_total": 0
      }
    },
    "global_totals": {
      "accepted_requests": 1,
      "rejected_requests": 0
    }
  }
}
```

---

# How to Test

## Valid Request

```bash
curl -X POST http://localhost:8080/request \
-H "Content-Type: application/json" \
-d '{
  "user_id":"user1",
  "payload":{"message":"hello"}
}'
```

---

## Invalid JSON

```bash
curl -X POST http://localhost:8080/request \
-H "Content-Type: application/json" \
-d '{invalid}'
```

---

## Missing Payload

```bash
curl -X POST http://localhost:8080/request \
-H "Content-Type: application/json" \
-d '{
  "user_id":"user1"
}'
```

---

## Rate Limiting Test

```bash
for i in {1..6}; do
curl -X POST http://localhost:8080/request \
-H "Content-Type: application/json" \
-d '{"user_id":"rate-user","payload":"test"}'
echo ""
done
```

---

## Concurrency Test

```bash
seq 1 10 | xargs -P10 -I{} curl -X POST http://localhost:8080/request \
-H "Content-Type: application/json" \
-d '{"user_id":"parallel-user","payload":"test"}'
```

Expected:
- only 5 requests accepted
- remaining requests rejected

---

# Production Limitations

Current implementation uses in-memory storage.

Limitations:
- data is lost on server restart
- single-instance only
- rate limiting is not shared across multiple servers
- no persistent database

---

# Future Improvements

Potential production improvements:

- Redis-based distributed rate limiter
- PostgreSQL persistence
- Authentication and authorization
- Metrics and monitoring
- Docker deployment
- Kubernetes deployment
- Centralized logging

---

# AI Tool Usage

AI tools were used for:
- architecture guidance
- API design suggestions
- concurrency design review
- README drafting assistance