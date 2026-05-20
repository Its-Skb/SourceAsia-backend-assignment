# Backend Assignment - Source Asia

This repository contains the implementation for the Source Asia Backend Assignment.

The project is built using Go and Gin framework and includes:

1. Part 1 - Rate Limited API
2. Part 2 - Product Catalog API

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

---

# Part 2 - Product Catalog API

This section implements a scalable in-memory product catalog system with optimized list and detail APIs.

---

# Product Architecture

Each product contains:

- core product metadata
- image URLs
- video URLs

Media are stored as URL strings only.

Example:

```json
{
  "name": "Widget A",
  "sku": "SKU-001",
  "image_urls": [
    "https://cdn.example.com/products/sku-001/img-1.jpg"
  ],
  "video_urls": [
    "https://cdn.example.com/products/sku-001/demo.mp4"
  ]
}
```

---

# Storage Design

The implementation uses separate in-memory structures for optimized performance.

## Product Metadata Store

Stores lightweight product information:

```go
Products map[string]*Product
```

Contains:
- id
- name
- sku
- created_at

---

## Product Media Store

Stores heavy media arrays separately:

```go
ProductMediaStore map[string]*ProductMedia
```

Contains:
- image_urls
- video_urls

---

## SKU Index

Used for fast uniqueness validation:

```go
SKUIndex map[string]string
```

Maps:

```text
sku -> product_id
```

This provides:

```text
O(1)
```

duplicate SKU lookups.

---

## Product Order Store

Maintains deterministic pagination order:

```go
ProductOrder []string
```

This avoids unstable ordering caused by Go map iteration.

---

# Why List vs Detail APIs Are Separate

The assignment specifically requires:

```text
GET /products must not load or serialize all media URLs
```

To satisfy this efficiently:

## GET /products

Returns only lightweight fields:

- id
- name
- sku
- image_count
- video_count
- thumbnail_url
- created_at

This keeps list responses fast even with:
- 1000+ products
- thousands of stored media URLs

The endpoint intentionally does NOT return:
- image_urls arrays
- video_urls arrays

---

## GET /products/{id}

Returns the complete product including:
- all image_urls
- all video_urls

This separation improves:
- memory usage
- response size
- serialization performance
- scalability

---

# Validation Rules

The API enforces:

| Validation | Rule |
|---|---|
| Product name | required, non-empty |
| SKU | required, non-empty, unique |
| URL scheme | must be http:// or https:// |
| URL max length | 2048 characters |
| Maximum URLs | 20 URLs per request array |

---

# Pagination

GET /products supports:

```http
GET /products?limit=20&offset=0
```

Defaults:
- limit = 20
- offset = 0

Maximum limit:
- 100

Pagination is deterministic using:

```go
ProductOrder []string
```

---

# Concurrency Safety

The implementation uses:

```go
sync.RWMutex
```

to protect:
- product creation
- media updates
- concurrent reads/writes

This prevents:
- race conditions
- duplicate SKU conflicts
- concurrent map access issues

---

# Optional Dataset Seeding

An optional seed utility is included:

```go
utils.SeedProducts(1000)
```

This can generate:
- 1000 products
- 10 image URLs per product

Useful for testing list endpoint scalability.

The seed utility is intentionally disabled by default.

---

# Example API Usage

---

## Create Product

```bash
curl -X POST http://localhost:8080/products \
-H "Content-Type: application/json" \
-d '{
  "name":"Widget A",
  "sku":"SKU-001",
  "image_urls":[
    "https://cdn.example.com/products/sku-001/img-1.jpg"
  ],
  "video_urls":[
    "https://cdn.example.com/products/sku-001/demo.mp4"
  ]
}'
```

---

## List Products

```bash
curl "http://localhost:8080/products?limit=20&offset=0"
```

---

## Get Product Detail

```bash
curl http://localhost:8080/products/{id}
```

---

## Add Product Media

```bash
curl -X POST http://localhost:8080/products/{id}/media \
-H "Content-Type: application/json" \
-d '{
  "image_urls":[
    "https://cdn.example.com/products/sku-001/img-2.jpg"
  ]
}'
```

---

# Production Improvements (Future)

For a real production deployment, the following improvements would be recommended:

---

## PostgreSQL

Instead of in-memory maps:
- products table
- product_media table
- indexed SKU column
- indexed product_id foreign keys

Benefits:
- persistence
- transactional consistency
- scalable querying
- indexing support

---

## Redis

Could be used for:
- distributed rate limiting
- caching hot product lists
- reducing database load

---

## CDN

Media URLs would normally point to:
- CloudFront
- Cloudflare
- Akamai
- S3-backed CDN

The backend would store only metadata and media references.

---

## Object Storage

Images/videos would typically be stored in:
- AWS S3
- Google Cloud Storage
- Azure Blob Storage

instead of being managed directly by the API.

---

## Cursor-Based Pagination

For very large datasets, offset pagination could be replaced with:
- cursor pagination
- keyset pagination

for improved scalability.

---

# Assignment Notes

- In-memory storage is intentionally used as allowed by the assignment
- Persistence across restarts is not implemented
- APIs are designed to prioritize clarity, scalability, and performance