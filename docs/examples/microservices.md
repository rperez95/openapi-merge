# Microservices Gateway Example

A complete example of merging multiple microservice APIs for an API Gateway.

## Scenario

You're building an API Gateway for an e-commerce platform with:

- **Users Service** - User accounts, profiles, and authentication
- **Orders Service** - Shopping cart, checkout, and order management
- **Products Service** - Product catalog, categories, and inventory

Each service has its own OpenAPI specification that needs to be unified.

## Project Structure

```
project/
├── examples/
│   ├── users-api.json
│   ├── orders-api.json
│   ├── products-api.json
│   └── merge-config.yaml
└── output/
    └── platform-api.json
```

## Example API Files

The example APIs are included in the `examples/` folder:

| File | Description | Endpoints |
|------|-------------|-----------|
| `users-api.json` | User management | `/v1/users`, `/v1/auth/*`, `/v1/users/*/profile` |
| `orders-api.json` | Order processing | `/v1/orders`, `/v1/cart/*`, `/v1/checkout` |
| `products-api.json` | Product catalog | `/v1/products`, `/v1/categories`, `/v1/*/inventory` |

## Configuration

```yaml title="examples/merge-config.yaml"
# =============================================================================
# E-Commerce Platform API - Merge Configuration
# =============================================================================

# API Metadata
info:
  title: "E-Commerce Platform API"
  description: |
    # E-Commerce Platform API

    Unified API specification for the E-Commerce platform, combining:

    - **Users Service**: User account management and authentication
    - **Orders Service**: Order processing, cart, and checkout
    - **Products Service**: Product catalog and inventory

    ## Authentication

    All endpoints require authentication via session cookie or Bearer token.
  version: "1.0.0"
  contact:
    name: "Platform Support"
    email: "support@example.com"

# Server Definitions
servers:
  - url: "https://api.example.com"
    description: "Production"
  - url: "https://staging-api.example.com"
    description: "Staging"
  - url: "http://localhost:8080"
    description: "Local Development"

# =============================================================================
# Security Configuration
# =============================================================================

securitySchemes:
  # Session cookie authentication
  cookieAuth:
    type: apiKey
    in: cookie
    name: session_id
    description: Session cookie for browser clients

  # Bearer token for API clients
  bearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT
    description: JWT Bearer token for API clients

# Global security (any of these methods)
security:
  - cookieAuth: []
  - bearerAuth: []

# =============================================================================
# Tag Ordering
# =============================================================================

tagOrder:
  - "Authentication"
  - "Users"
  - "Profile"
  - "Products"
  - "Categories"
  - "Cart"
  - "Orders"
  - "Checkout"

# =============================================================================
# Input Files
# =============================================================================

inputs:
  # Users Service
  - inputFile: users-api.json
    dispute:
      prefix: "Users_"  # Prefix schemas to avoid collisions
    pathModification:
      stripStart: "/v1"
      prepend: "/v1/users-service"
    includeExtraParameters:
      - name: X-Request-ID
        in: header
        description: "Request tracking ID"
        required: false
        schema:
          type: string
    description:
      append: true
      title:
        value: "Users Service"
        headingLevel: 2

  # Orders Service
  - inputFile: orders-api.json
    dispute:
      prefix: "Orders_"
    pathModification:
      stripStart: "/v1"
      prepend: "/v1/orders-service"
    includeExtraParameters:
      - name: X-Request-ID
        in: header
        description: "Request tracking ID"
        required: false
        schema:
          type: string
    description:
      append: true
      title:
        value: "Orders Service"
        headingLevel: 2

  # Products Service
  - inputFile: products-api.json
    dispute:
      prefix: "Products_"
    pathModification:
      stripStart: "/v1"
      prepend: "/v1/products-service"
    includeExtraParameters:
      - name: X-Request-ID
        in: header
        description: "Request tracking ID"
        required: false
        schema:
          type: string
    description:
      append: true
      title:
        value: "Products Service"
        headingLevel: 2

# =============================================================================
# Output
# =============================================================================

output: platform-api.json
```

## Run the Merge

```bash
cd examples
openapi-merge merge --config merge-config.yaml --verbose
```

Output:
```
Starting merge with 3 input files
Output file: /path/to/examples/platform-api.json
Processing input 1: users-api.json
Processing input 2: orders-api.json
Processing input 3: products-api.json
Successfully merged 3 specifications into platform-api.json
```

## Result Structure

The merged API has a clear version-first structure:

```
Platform API
├── /v1/users-service/
│   ├── /users
│   ├── /users/{userId}
│   ├── /users/{userId}/profile
│   ├── /auth/login
│   └── /auth/logout
├── /v1/orders-service/
│   ├── /orders
│   ├── /orders/{orderId}
│   ├── /cart
│   ├── /cart/items
│   └── /checkout
└── /v1/products-service/
    ├── /products
    ├── /products/{productId}
    ├── /categories
    ├── /products/{productId}/inventory
    └── /products/{productId}/reviews
```

## Path Transformation Summary

| Original Path | Service | Merged Path |
|---------------|---------|-------------|
| `/v1/users` | Users | `/v1/users-service/users` |
| `/v1/auth/login` | Users | `/v1/users-service/auth/login` |
| `/v1/orders` | Orders | `/v1/orders-service/orders` |
| `/v1/cart` | Orders | `/v1/orders-service/cart` |
| `/v1/products` | Products | `/v1/products-service/products` |
| `/v1/categories` | Products | `/v1/products-service/categories` |

## API Gateway Integration

### Kong Gateway

```yaml title="kong.yml"
services:
  - name: users-service
    url: http://users-backend:8081
    routes:
      - name: users-routes
        paths:
          - /v1/users-service

  - name: orders-service
    url: http://orders-backend:8082
    routes:
      - name: orders-routes
        paths:
          - /v1/orders-service

  - name: products-service
    url: http://products-backend:8083
    routes:
      - name: products-routes
        paths:
          - /v1/products-service
```

### AWS API Gateway

Import the merged `platform-api.json` directly into AWS API Gateway as an OpenAPI 3.0 specification.

### Nginx

```nginx
upstream users_backend {
    server users-service:8081;
}

upstream orders_backend {
    server orders-service:8082;
}

upstream products_backend {
    server products-service:8083;
}

server {
    listen 80;

    location /v1/users-service/ {
        proxy_pass http://users_backend/;
    }

    location /v1/orders-service/ {
        proxy_pass http://orders_backend/;
    }

    location /v1/products-service/ {
        proxy_pass http://products_backend/;
    }
}
```

### Traefik

```yaml title="traefik.yml"
http:
  routers:
    users:
      rule: "PathPrefix(`/v1/users-service`)"
      service: users-service
    orders:
      rule: "PathPrefix(`/v1/orders-service`)"
      service: orders-service
    products:
      rule: "PathPrefix(`/v1/products-service`)"
      service: products-service

  services:
    users-service:
      loadBalancer:
        servers:
          - url: "http://users:8081"
    orders-service:
      loadBalancer:
        servers:
          - url: "http://orders:8082"
    products-service:
      loadBalancer:
        servers:
          - url: "http://products:8083"
```

## Generated Documentation

The merged specification includes combined API documentation:

```markdown
# E-Commerce Platform API

Unified API specification for the E-Commerce platform...

## Authentication

All endpoints require authentication via session cookie or Bearer token.

## Users Service

User management microservice for handling user accounts, authentication, and profiles.

## Orders Service

Order management microservice for handling customer orders, cart, and checkout.

## Products Service

Product catalog microservice for managing products, categories, and inventory.
```

## Try It Yourself

1. Clone the repository
2. Navigate to the examples folder
3. Run the merge:

```bash
cd examples
openapi-merge merge --config merge-config.yaml -o my-api.json
```

4. View the result in Swagger Editor or import into Postman

## Next Steps

- [Advanced Configuration](advanced.md) - OAuth2, filtering, and more
- [Security Schemes](../configuration/security.md) - Authentication options
- [Path Modification](../configuration/paths.md) - Path transformation guide
