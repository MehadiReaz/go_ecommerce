# API Documentation

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Response Format

All API responses follow this format:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## Endpoints

### Authentication

#### Sign Up
```http
POST /api/v1/auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGc...",
    "refresh_token": "uuid...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

### Products

#### List Products
```http
GET /api/v1/products?category_id=1&min_price=10&max_price=100&limit=20&offset=0
```

#### Get Product
```http
GET /api/v1/products/{id}
```

#### Search Products
```http
GET /api/v1/products/search?q=laptop&limit=20
```

### Cart

#### Get Cart
```http
GET /api/v1/cart
Authorization: Bearer <token>
```

#### Add to Cart
```http
POST /api/v1/cart/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 2
}
```

### Orders

#### Create Order
```http
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "shipping_address_id": 1,
  "billing_address_id": 1,
  "payment_method": "stripe"
}
```

#### Get Orders
```http
GET /api/v1/orders?limit=20&offset=0
Authorization: Bearer <token>
```

### Payments

#### Create Payment
```http
POST /api/v1/payments
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": 1,
  "payment_method": "stripe",
  "currency": "USD"
}
```

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |
