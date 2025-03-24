# Go Marketplace Documentation

Welcome to the Go Marketplace documentation. This documentation provides comprehensive information about the microservices architecture, API endpoints, and testing procedures.

## Table of Contents

1. [Architecture Overview](./architecture/README.md)
   - System Components
   - Service Communication
   - Data Flow

2. [API Reference](./api-reference/README.md)
   - GraphQL Schema
   - Available Queries and Mutations
   - Data Types

3. [GraphQL Playground Guide](./graphql-playground/README.md)
   - Getting Started
   - Testing Queries
   - Testing Mutations
   - Example Scenarios

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/donaldnash/go-marketplace.git
   ```

2. Start the services:
   ```bash
   docker-compose up --build
   ```

3. Access the GraphQL Playground:
   - Open your browser and navigate to [http://localhost:8080/playground](http://localhost:8080/playground)

## Services

The marketplace consists of the following microservices:

1. **Account Service** (Port 8081)
   - User account management
   - Authentication and authorization

2. **Catalog Service** (Port 8082)
   - Product management
   - Product search and filtering

3. **Order Service** (Port 8083)
   - Order processing
   - Order history

4. **GraphQL Gateway** (Port 8080)
   - API Gateway
   - Service aggregation
   - GraphQL interface

## Contributing

Please read our [Contributing Guidelines](../CONTRIBUTING.md) before submitting any changes.

## Support

If you encounter any issues or have questions, please:
1. Check the documentation in the relevant sections
2. Look for existing GitHub issues
3. Create a new issue if needed 