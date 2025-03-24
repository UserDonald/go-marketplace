# Go Marketplace

A modern, production-grade marketplace platform built with microservices architecture, demonstrating best practices in Go development and distributed systems.

## Key Features

- 🔐 **Secure Authentication**: Complete user account management with JWT tokens
- 🎯 **Modern Architecture**: GraphQL API Gateway with gRPC inter-service communication
- 📦 **Containerized**: Docker-based deployment with multi-stage builds
- 🛠 **Best Practices**: Clean code architecture, error handling, and comprehensive testing
- 🔄 **Scalable Design**: Independent services that can be scaled separately

## Tech Stack

### Core Technologies
- Go 1.24.1
- GraphQL (gqlgen) for API Gateway
- gRPC for service communication
- PostgreSQL for data persistence
- Docker & Docker Compose

### Key Dependencies
- github.com/99designs/gqlgen
- google.golang.org/grpc
- github.com/lib/pq

## Quick Setup

### Prerequisites
- Go 1.24.1+
- Docker & Docker Compose
- Git

### Run the Project
```bash
# Clone the repository
git clone https://github.com/yourusername/go-marketplace.git
cd go-marketplace

# Start all services
docker-compose up -d

# View service logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Access Points
- Account Service: http://localhost:8080
  - Authentication endpoints
  - User management API
- GraphQL Playground: http://localhost:8081 (coming soon)
  - Interactive API documentation
  - Query/Mutation testing

## Project Structure
```
go-marketplace/
├── account/          # Authentication & user management
│   ├── cmd/         # Service entry points
│   ├── pb/          # Protocol buffer definitions
│   ├── repository/  # Data access layer
│   └── service/     # Business logic
├── catalog/         # Product management service
├── order/          # Order processing service
├── graphql/        # API gateway
│   ├── schema/     # GraphQL schema definitions
│   └── resolvers/  # Query/Mutation implementations
└── docker-compose.yaml
```

## Service Details

### Account Service (Complete)
- User authentication with JWT
- Password hashing and validation
- User profile management
- gRPC API for internal service communication
- PostgreSQL for user data persistence

### GraphQL Gateway (In Progress)
- Unified API entry point
- Type-safe schema generation
- Authentication middleware
- Service aggregation layer

### Catalog Service (In Development)
- Product CRUD operations
- Category management
- Search functionality
- Inventory tracking

### Order Service (Planned)
- Order processing workflow
- Payment integration
- Order status management
- History tracking

## Development Status

### Completed
- Account service implementation
- Basic project architecture
- Docker configuration
- Database schema and migrations
- Authentication flow

### In Progress
- GraphQL gateway implementation
- Catalog service development
- Service integration
- Testing infrastructure

### Coming Soon
- Order service
- Service discovery
- Monitoring and logging
- CI/CD pipeline
- Kubernetes deployment

## Troubleshooting

### Common Issues
- **Port Conflicts**: Ensure ports 8080 (Account), 8081 (GraphQL), and 5432 (PostgreSQL) are available
- **Database**: Default credentials in docker-compose.yaml
- **Logs**: Use `docker-compose logs -f [service]` for debugging

### Database Access
- Host: localhost
- Port: 5432
- Username: postgres
- Password: postgres
- Database: postgres
