# Simple Ecommerce API

A RESTful API for an online shop built with Go and PostgreSQL. This API provides endpoints for managing products, processing orders, and handling payments.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
  - [Products](#products)
  - [Orders](#orders)
  - [Admin](#admin)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Features

- Product management (CRUD operations)
- Order processing and checkout
- Payment confirmation
- Admin-only routes with authentication
- Secure passcode generation for orders

## Project Structure

```
ecommerce-API/
├── handler/             # HTTP request handlers
│   ├── order.go         # Order-related handlers
│   └── product.go       # Product-related handlers
├── middleware/          # HTTP middleware
│   └── admin.go         # Admin authentication middleware
├── model/               # Data models and database operations
│   ├── order.go         # Order-related models and DB operations
│   └── product.go       # Product-related models and DB operations
├── .env                 # Environment variables (not in version control)
├── .env-example         # Example environment variables
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point and server setup
└── README.md            # Project documentation
```

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```
git clone <repository-url>
cd ecommerce-API
```

2. Install dependencies:
```
go mod download
```

3. Create a `.env` file based on `.env-example` and configure your environment variables.

4. Run the application:
```
go run .
```

The server will start on port 8080 by default.

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_URI=postgres://username:password@localhost:5432/dbname
ADMIN_SECRET=your-admin-secret-key
```

## Database Schema

The application automatically creates the following tables on startup:

### Products Table
- `id`: UUID, primary key
- `name`: Product name
- `price`: Product price in cents/smallest currency unit
- `is_deleted`: Soft delete flag

### Orders Table
- `id`: UUID, primary key
- `email`: Customer email
- `address`: Shipping address
- `passcode`: Hashed passcode for order authentication
- `paid_at`: Timestamp of payment
- `paid_bank`: Bank used for payment
- `paid_account`: Account number used for payment
- `grand_total`: Total order amount

### Order Details Table
- `id`: UUID, primary key
- `order_id`: Reference to Orders table
- `product_id`: Reference to Products table
- `quantity`: Number of products ordered
- `price`: Price of product at time of order
- `total`: Total price for this line item

## API Endpoints

### Products

#### GET /api/v1/products
- Description: Get all available products
- Authentication: None
- Response: 200 OK with array of products

#### GET /api/v1/products/:id
- Description: Get a specific product by ID
- Authentication: None
- Response: 200 OK with product details or 404 Not Found

### Orders

#### POST /api/v1/checkout
- Description: Create a new order
- Authentication: None
- Request Body:
  ```json
  {
    "email": "customer@example.com",
    "address": "123 Main St",
    "products": [
      {
        "id": "product-uuid",
        "quantity": 2
      }
    ]
  }
  ```
- Response: 200 OK with order details including passcode

#### POST /api/v1/orders/:id/confirm
- Description: Confirm payment for an order
- Authentication: None (requires passcode)
- Request Body:
  ```json
  {
    "amount": 10000,
    "bank": "Example Bank",
    "accountNumber": "123456789",
    "passcode": "abc123"
  }
  ```
- Response: 200 OK with updated order details

#### GET /api/v1/orders/:id?passcode=abc123
- Description: Get order details
- Authentication: None (requires passcode as query parameter)
- Response: 200 OK with order details

### Admin

#### POST /admin/products
- Description: Create a new product
- Authentication: Admin (Authorization header)
- Request Body:
  ```json
  {
    "name": "Product Name",
    "price": 10000
  }
  ```
- Response: 201 Created with product details

#### PUT /admin/products/:id
- Description: Update an existing product
- Authentication: Admin (Authorization header)
- Request Body:
  ```json
  {
    "name": "Updated Product Name",
    "price": 15000
  }
  ```
- Response: 200 OK with updated product details

#### DELETE /admin/products/:id
- Description: Delete a product (soft delete)
- Authentication: Admin (Authorization header)
- Response: 204 No Content

## Authentication

### Admin Authentication
Admin endpoints require an `Authorization` header with the admin secret key:

```
Authorization: your-admin-secret-key
```

### Order Authentication
Orders are protected with a unique passcode generated during checkout. This passcode is required for confirming payment and viewing order details.

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful request
- `201 Created`: Resource successfully created
- `204 No Content`: Resource successfully deleted
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication failed
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

## Examples

### Creating a Product (Admin)

```bash
curl -X POST http://localhost:8080/admin/products \
  -H "Authorization: your-admin-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"name": "Smartphone", "price": 1000000}'
```

### Placing an Order

```bash
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "address": "123 Main St",
    "products": [
      {
        "id": "product-uuid",
        "quantity": 1
      }
    ]
  }'
```

### Confirming Payment

```bash
curl -X POST http://localhost:8080/api/v1/orders/order-uuid/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000000,
    "bank": "Example Bank",
    "accountNumber": "123456789",
    "passcode": "abc123"
  }'
```