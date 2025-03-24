# Go Marketplace

A modern, production-grade marketplace platform built with microservices architecture, demonstrating best practices in Go development and distributed systems.

## Key Features

- 🎯 **Modern Architecture**: GraphQL API Gateway with gRPC inter-service communication
- 🔍 **Full-Text Search**: Elasticsearch-powered product catalog with advanced search capabilities
- 📦 **Containerized**: Docker-based deployment with multi-stage builds
- 🛠 **Best Practices**: Clean code architecture and comprehensive error handling
- 🔄 **Scalable Design**: Independent services that can be scaled separately

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
docker-compose up --build

# Access GraphQL Playground
open http://localhost:8080/playground

# Stop services
docker-compose down
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
│   └── client.go    # gRPC client implementation
├── catalog/         # Product catalog service
│   ├── cmd/        # Service entry point
│   ├── pb/         # Protocol buffer definitions
│   └── client.go   # gRPC client implementation
├── order/          # Order processing service
│   ├── cmd/        # Service entry point
│   ├── pb/         # Protocol buffer definitions
│   └── client.go   # gRPC client implementation
├── graphql/        # API gateway
│   ├── schema.graphql  # GraphQL schema
│   └── resolvers/     # Query/Mutation implementations
├── docs/           # Documentation
└── docker-compose.yaml
```

## Service Details

### Account Service (Port 8081)
- Account creation and management
- PostgreSQL for data persistence
- gRPC API for service communication
- KSUID for ID generation
- Clean architecture with repository pattern

### Catalog Service (Port 8082)
- Product management with Elasticsearch
- Full-text search functionality
- Multi-match search across product fields
- Pagination support
- gRPC API for service communication

### Order Service (Port 8083)
- Order processing and management
- PostgreSQL for order data
- Integration with Account and Catalog services
- Order history tracking
- gRPC API for service communication

### GraphQL Gateway (Port 8080)
- Unified API entry point
- Interactive GraphQL Playground
- Service aggregation layer
- CORS support
- Error handling and validation

## Development Status

### Completed
- ✅ Basic project architecture
- ✅ Docker configuration with health checks
- ✅ Database schemas and setup
- ✅ gRPC service communication
- ✅ GraphQL gateway implementation
- ✅ Account service implementation
- ✅ Catalog service with Elasticsearch
- ✅ Order service implementation
- ✅ Service integration
- ✅ Error handling
- ✅ Documentation

### In Progress
- 🔄 Service resilience improvements
- 🔄 Logging enhancements
- 🔄 Performance optimization
- 🔄 Integration testing

### Coming Soon
- 📅 Service monitoring
- 📅 Caching layer
- 📅 Load balancing
- 📅 Service mesh integration
- 📅 Kubernetes deployment configuration

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
- **Service Logs**: Use `docker-compose logs -f [service]` for debugging

### Database Access
- Account DB:
  - Host: localhost:5431
  - Database: account
  - Username: postgres
  - Password: postgres

- Order DB:
  - Host: localhost:5433
  - Database: order
  - Username: postgres
  - Password: postgres

- Catalog DB (Elasticsearch):
  - Host: localhost:9200

## Documentation

For detailed information about the API, testing procedures, and architecture, please refer to the [documentation](./docs/README.md).
