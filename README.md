# Go Marketplace

A microservices-based marketplace application built with Go, featuring GraphQL API gateway and gRPC communication between services.

## Project Overview

This project implements a marketplace system with multiple microservices architecture. It uses modern technologies and best practices including:

- GraphQL API Gateway for client interactions
- gRPC for inter-service communication
- PostgreSQL for data persistence
- Docker for containerization
- Microservices architecture

## Services

### 1. Account Service (Implemented)
- User account management
- Authentication and authorization
- gRPC-based service
- PostgreSQL database for user data
- Features:
  - User registration
  - User authentication
  - Account management

### 2. Catalog Service (Planned)
- Product catalog management
- Features to be implemented:
  - Product listing
  - Product categories
  - Product search
  - Inventory management

### 3. Order Service (Planned)
- Order processing and management
- Features to be implemented:
  - Order creation
  - Order status tracking
  - Order history
  - Payment integration

### 4. GraphQL Gateway (In Progress)
- Central API gateway
- GraphQL schema and resolvers
- Client-facing API
- Features:
  - Account mutations and queries
  - Integration with other services (planned)

## Technology Stack

- **Language:** Go 1.24.1
- **API Gateway:** GraphQL with gqlgen
- **Service Communication:** gRPC
- **Database:** PostgreSQL
- **Containerization:** Docker
- **Dependencies:**
  - github.com/99designs/gqlgen
  - google.golang.org/grpc
  - github.com/lib/pq
  - Other utility packages

## Project Structure

```
.
├── account/           # Account service
│   ├── cmd/          # Service entry points
│   ├── pb/           # Protocol buffer definitions
│   ├── repository.go # Data access layer
│   ├── server.go     # gRPC server implementation
│   └── service.go    # Business logic
├── catalog/          # Catalog service (planned)
├── order/           # Order service (planned)
├── graphql/         # GraphQL gateway
│   ├── schema.graphql    # GraphQL schema
│   ├── resolvers/        # GraphQL resolvers
│   └── generated/        # Generated GraphQL code
└── docker-compose.yaml   # Docker composition
```

## Current Status

### Completed
- Basic project structure and architecture
- Account service with gRPC implementation
- GraphQL gateway setup
- Docker configuration
- Database schema for account service

### In Progress
- GraphQL resolvers implementation
- Service integration
- Testing infrastructure

### To Be Implemented
- Catalog service
- Order service
- Authentication and authorization
- Service discovery
- Logging and monitoring
- CI/CD pipeline