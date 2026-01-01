# Basic Merge Example

A simple example of merging two OpenAPI specifications.

## Scenario

You have two microservices:

- **Users Service** - Manages user accounts
- **Orders Service** - Manages customer orders

You want to create a unified API for your API Gateway.

## Input Files

=== "users-api.json"

    ```json
    {
      "openapi": "3.0.0",
      "info": {
        "title": "Users API",
        "version": "1.0.0",
        "description": "User management service"
      },
      "paths": {
        "/users": {
          "get": {
            "summary": "List all users",
            "tags": ["Users"],
            "responses": {
              "200": {
                "description": "List of users",
                "content": {
                  "application/json": {
                    "schema": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/User" }
                    }
                  }
                }
              }
            }
          },
          "post": {
            "summary": "Create a user",
            "tags": ["Users"],
            "requestBody": {
              "content": {
                "application/json": {
                  "schema": { "$ref": "#/components/schemas/CreateUser" }
                }
              }
            },
            "responses": {
              "201": { "description": "User created" }
            }
          }
        },
        "/users/{id}": {
          "get": {
            "summary": "Get user by ID",
            "tags": ["Users"],
            "parameters": [
              {
                "name": "id",
                "in": "path",
                "required": true,
                "schema": { "type": "string" }
              }
            ],
            "responses": {
              "200": {
                "description": "User details",
                "content": {
                  "application/json": {
                    "schema": { "$ref": "#/components/schemas/User" }
                  }
                }
              }
            }
          }
        }
      },
      "components": {
        "schemas": {
          "User": {
            "type": "object",
            "properties": {
              "id": { "type": "string" },
              "email": { "type": "string" },
              "name": { "type": "string" }
            }
          },
          "CreateUser": {
            "type": "object",
            "required": ["email", "name"],
            "properties": {
              "email": { "type": "string" },
              "name": { "type": "string" }
            }
          }
        }
      }
    }
    ```

=== "orders-api.json"

    ```json
    {
      "openapi": "3.0.0",
      "info": {
        "title": "Orders API",
        "version": "1.0.0",
        "description": "Order management service"
      },
      "paths": {
        "/orders": {
          "get": {
            "summary": "List all orders",
            "tags": ["Orders"],
            "responses": {
              "200": {
                "description": "List of orders",
                "content": {
                  "application/json": {
                    "schema": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Order" }
                    }
                  }
                }
              }
            }
          },
          "post": {
            "summary": "Create an order",
            "tags": ["Orders"],
            "requestBody": {
              "content": {
                "application/json": {
                  "schema": { "$ref": "#/components/schemas/CreateOrder" }
                }
              }
            },
            "responses": {
              "201": { "description": "Order created" }
            }
          }
        }
      },
      "components": {
        "schemas": {
          "Order": {
            "type": "object",
            "properties": {
              "id": { "type": "string" },
              "userId": { "type": "string" },
              "total": { "type": "number" },
              "status": { "type": "string" }
            }
          },
          "CreateOrder": {
            "type": "object",
            "required": ["userId", "items"],
            "properties": {
              "userId": { "type": "string" },
              "items": {
                "type": "array",
                "items": { "type": "string" }
              }
            }
          }
        }
      }
    }
    ```

## Configuration

```yaml title="merge-config.yaml"
# API information
info:
  title: "E-Commerce Platform API"
  description: |
    Unified API for the E-Commerce platform.
    
    ## Services
    - Users: User account management
    - Orders: Order processing
  version: "1.0.0"
  contact:
    name: "API Support"
    email: "api@example.com"

# Server configuration
servers:
  - url: "https://api.example.com"
    description: "Production"
  - url: "https://staging-api.example.com"
    description: "Staging"

# Input files
inputs:
  - inputFile: users-api.json
    pathModification:
      prepend: "/users-service"
  
  - inputFile: orders-api.json
    pathModification:
      prepend: "/orders-service"

# Output file
output: platform-api.json
```

## Run the Merge

```bash
openapi-merge merge --config merge-config.yaml
```

Output:
```
Successfully merged 2 specifications into platform-api.json
```

## Result

The merged `platform-api.json` contains:

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "E-Commerce Platform API",
    "description": "Unified API for the E-Commerce platform.\n\n## Services\n- Users: User account management\n- Orders: Order processing",
    "version": "1.0.0",
    "contact": {
      "name": "API Support",
      "email": "api@example.com"
    }
  },
  "servers": [
    {
      "url": "https://api.example.com",
      "description": "Production"
    },
    {
      "url": "https://staging-api.example.com",
      "description": "Staging"
    }
  ],
  "paths": {
    "/users-service/users": { ... },
    "/users-service/users/{id}": { ... },
    "/orders-service/orders": { ... }
  },
  "components": {
    "schemas": {
      "User": { ... },
      "CreateUser": { ... },
      "Order": { ... },
      "CreateOrder": { ... }
    }
  }
}
```

## Path Summary

| Original | Merged |
|----------|--------|
| `/users` (users-api) | `/users-service/users` |
| `/users/{id}` (users-api) | `/users-service/users/{id}` |
| `/orders` (orders-api) | `/orders-service/orders` |

## Next Steps

- [Microservices Gateway Example](microservices.md) - More complex scenario
- [Advanced Configuration](advanced.md) - Security, filtering, and more

