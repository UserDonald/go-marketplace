# Go Marketplace

A modern, production-grade marketplace platform built with microservices architecture, demonstrating best practices in Go development and distributed systems.

## Key Features

- 🎯 **Modern Architecture**: GraphQL API Gateway with gRPC inter-service communication
- 🔍 **Full-Text Search**: Elasticsearch-powered product catalog with advanced search capabilities
- 📦 **Containerized**: Docker-based deployment with multi-stage builds
- 🛠 **Best Practices**: Clean code architecture with comprehensive error handling and validation
- 🔄 **Scalable Design**: Independent services that can be scaled separately
- 🔒 **Robust Error Handling**: Detailed error messages, input validation, and graceful recovery
- 📝 **Structured Logging**: Comprehensive logging with context and file information

## Tech Stack

### Core Technologies
- Go 1.23.0
- GraphQL (gqlgen) for API Gateway
- gRPC for service communication
- PostgreSQL for structured data
- Elasticsearch for product search
- Docker & Docker Compose

### Key Dependencies
- github.com/99designs/gqlgen
- google.golang.org/grpc
- github.com/lib/pq
- github.com/olivere/elastic/v7
- github.com/segmentio/ksuid
- github.com/kelseyhightower/envconfig
- github.com/tinrab/retry

## Quick Setup

### Prerequisites
- Go 1.23.0+
- Docker & Docker Compose
- Git

### Run the Project
```bash
# Clone the repository
git clone https://github.com/donaldnash/go-marketplace.git
cd go-marketplace

# Start all services
docker compose up --build

# Access GraphQL Playground
open http://localhost:8080/playground

# Stop services
docker compose down
```

### Access Points
- GraphQL Gateway: http://localhost:8080
  - GraphQL Playground: http://localhost:8080/playground
  - API Documentation
  - Query/Mutation testing
- Account Service: http://localhost:8081 (gRPC)
- Catalog Service: http://localhost:8082 (gRPC)
- Order Service: http://localhost:8083 (gRPC)

## Project Structure
```
go-marketplace/
├── account/          # User account management
│   ├── cmd/         # Service entry point
│   ├── pb/          # Protocol buffer definitions
│   ├── client.go    # gRPC client implementation
│   ├── server.go    # gRPC server implementation
│   ├── service.go   # Business logic
│   └── repository.go # Data access layer
├── catalog/         # Product catalog service
│   ├── cmd/        # Service entry point
│   ├── pb/         # Protocol buffer definitions
│   ├── client.go   # gRPC client implementation
│   ├── server.go   # gRPC server implementation
│   ├── service.go  # Business logic
│   └── repository.go # Data access layer
├── order/          # Order processing service
│   ├── cmd/        # Service entry point
│   ├── pb/         # Protocol buffer definitions
│   ├── client.go   # gRPC client implementation
│   ├── server.go   # gRPC server implementation
│   ├── service.go  # Business logic
│   └── repository.go # Data access layer
├── graphql/        # API gateway
│   ├── schema/     # GraphQL schema
│   ├── generated/  # Generated GraphQL code
│   └── resolvers/  # Query/Mutation implementations
├── docs/           # Documentation
└── docker-compose.yaml
```

## Service Details

### Account Service (Port 8081)
- Account creation and management
- PostgreSQL for data persistence
- gRPC API for service communication
- KSUID for unique ID generation
- Clean architecture with repository pattern
- Input validation for account creation
- Error handling with detailed messages
- Graceful shutdown with resource cleanup

### Catalog Service (Port 8082)
- Product management with Elasticsearch
- Full-text search with multi-match queries
- Product creation and retrieval
- Pagination support (max 100 items)
- Efficient bulk product retrieval
- Error handling with detailed messages
- Graceful shutdown with resource cleanup

### Order Service (Port 8083)
- Order processing and management
- PostgreSQL for order data
- Integration with Account and Catalog services
- Transaction support for order creation
- Product validation and price calculation
- Order history tracking by account
- Error handling with detailed messages
- Graceful shutdown with resource cleanup

### GraphQL Gateway (Port 8080)
- Unified API entry point
- Interactive GraphQL Playground
- Service aggregation and orchestration
- Request timeout handling (3s default)
- Detailed error messages with context
- Input validation and sanitization
- Graceful shutdown with resource cleanup

## Development Status

### Completed
- ✅ Basic project architecture and setup
- ✅ Docker configuration
- ✅ Database schemas and initialization
- ✅ gRPC service communication
- ✅ GraphQL gateway implementation
- ✅ Account service CRUD operations
- ✅ Catalog service with Elasticsearch
- ✅ Order service with transactions
- ✅ Service integration and testing
- ✅ Error handling
  - Input validation
  - Detailed error messages
  - Transaction management
  - Resource cleanup
  - Context handling
- ✅ Structured logging
  - Request context
  - Error details
  - Service operations
- ✅ Documentation
  - API reference
  - GraphQL playground guide
  - Architecture overview

### In Progress
- 🔄 Health check implementation
- 🔄 Service metrics collection
- 🔄 Integration testing
- 🔄 API documentation updates

### Planned
- 📅 Product stock management
- 📅 Caching layer
- 📅 Authentication and authorization
- 📅 Rate limiting
- 📅 Service monitoring
- 📅 Performance optimization
- 📅 Load balancing
- 📅 Kubernetes deployment

## Development Setup

### Dependencies
```bash
# Download dependencies
go mod download

# Vendor dependencies
go mod vendor
```

### Generate Protobuf Files
```bash
# Install protoc compiler (if not already installed)
# For macOS:
brew install protobuf
# For Linux:
apt-get install -y protobuf-compiler
# For Windows:
# Download from https://github.com/protocolbuffers/protobuf/releases

# Install protoc-gen-go and protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    account/pb/*.proto catalog/pb/*.proto order/pb/*.proto
```

## Troubleshooting

### Common Issues
- **Port Conflicts**: Ensure ports 8080-8083 are available
- **Database Connections**: Check service URLs in docker-compose.yaml
- **Service Logs**: Use `docker compose logs -f [service]` for debugging
- **Protobuf Generation**: Ensure protoc and required plugins are in PATH

### Database Access
- Account DB:
  - Host: localhost:5432
  - Database: account
  - Username: postgres
  - Password: postgres

- Order DB:
  - Host: localhost:5432
  - Database: order
  - Username: postgres
  - Password: postgres

- Catalog DB (Elasticsearch):
  - Host: localhost:9200
  - No authentication required

## Documentation

For detailed information about the API, testing procedures, and architecture, please refer to:
- [API Reference](./docs/api-reference/README.md)
- [Architecture Overview](./docs/architecture/README.md)
- [GraphQL Playground Guide](./docs/graphql-playground/README.md)
