# E-Commerce Platform

A comprehensive e-commerce backend API built with Go.

## Features

- **User Management**: Authentication, authorization, profile management
- **Product Catalog**: Products, categories, search functionality
- **Shopping Cart**: Add/remove items, quantity management
- **Order Management**: Order placement, tracking, cancellation
- **Payment Processing**: Stripe and bKash integration
- **Inventory Management**: Stock tracking, reservation system
- **Reviews & Ratings**: Product reviews and ratings
- **Shipping**: Address management
- **Notifications**: Email, SMS, push notifications

## Project Structure

```
ecommerce_project/
├── cmd/
│   ├── api/              # API server entry point
│   └── worker/           # Background workers
├── internal/
│   ├── app/              # Application setup (router, middleware)
│   ├── config/           # Configuration management
│   ├── user/             # User domain
│   ├── product/          # Product domain
│   ├── category/         # Category domain
│   ├── cart/             # Cart domain
│   ├── order/            # Order domain
│   ├── payment/          # Payment domain
│   ├── inventory/        # Inventory domain
│   ├── auth/             # Authentication & authorization
│   ├── notification/     # Notification services
│   ├── review/           # Review domain
│   └── shipping/         # Shipping domain
├── pkg/
│   ├── db/               # Database connection
│   ├── cache/            # Redis cache
│   ├── logger/           # Logging utilities
│   ├── utils/            # Common utilities
│   ├── email/            # Email sending
│   └── payment/          # Payment gateway integrations
├── config/               # Configuration files
├── scripts/              # Database migrations and seed data
└── test/                 # Tests
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Redis 6 or higher

### Installation

1. Clone the repository:
```bash
cd ecommerce_project
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment file:
```bash
cp .env.example .env
```

4. Update `.env` with your configuration

5. Run database migrations:
```bash
chmod +x scripts/migrate.sh
./scripts/migrate.sh
```

6. Seed the database (optional):
```bash
go run scripts/seed_data.go
```

### Running the Application

Start the API server:
```bash
go run cmd/api/main.go
```

Start the background workers:
```bash
go run cmd/worker/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Users
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update profile
- `PUT /api/v1/users/me/password` - Change password

### Products
- `GET /api/v1/products` - List products
- `GET /api/v1/products/{id}` - Get product details
- `GET /api/v1/products/search?q=query` - Search products
- `POST /api/v1/admin/products` - Create product (admin)
- `PUT /api/v1/admin/products/{id}` - Update product (admin)
- `DELETE /api/v1/admin/products/{id}` - Delete product (admin)

### Categories
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/{id}` - Get category details
- `POST /api/v1/admin/categories` - Create category (admin)
- `PUT /api/v1/admin/categories/{id}` - Update category (admin)
- `DELETE /api/v1/admin/categories/{id}` - Delete category (admin)

### Cart
- `GET /api/v1/cart` - Get cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items/{id}` - Update cart item
- `DELETE /api/v1/cart/items/{id}` - Remove item from cart
- `DELETE /api/v1/cart/clear` - Clear cart

### Orders
- `GET /api/v1/orders` - List user orders
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/{id}` - Get order details
- `POST /api/v1/orders/{id}/cancel` - Cancel order

### Payments
- `POST /api/v1/payments` - Create payment
- `GET /api/v1/payments/{id}` - Get payment details
- `POST /api/v1/payments/webhook/stripe` - Stripe webhook
- `POST /api/v1/payments/webhook/bkash` - bKash webhook

### Reviews
- `GET /api/v1/products/{id}/reviews` - Get product reviews
- `POST /api/v1/reviews` - Create review
- `PUT /api/v1/reviews/{id}` - Update review
- `DELETE /api/v1/reviews/{id}` - Delete review

### Shipping
- `GET /api/v1/shipping/addresses` - List addresses
- `POST /api/v1/shipping/addresses` - Create address
- `PUT /api/v1/shipping/addresses/{id}` - Update address
- `DELETE /api/v1/shipping/addresses/{id}` - Delete address

## Environment Variables

See `.env.example` for all available environment variables.

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
