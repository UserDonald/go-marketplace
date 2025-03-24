# API Reference

This document provides detailed information about the Go Marketplace GraphQL API.

## Schema Types

### Account
```graphql
type Account {
  id: String!
  name: String!
  orders: [Order!]!
}

input AccountInput {
  name: String!
}
```

### Product
```graphql
type Product {
  id: String!
  name: String!
  description: String!
  price: Float!
}

input ProductInput {
  name: String!
  description: String!
  price: Float!
}
```

### Order
```graphql
type Order {
  id: String!
  createdAt: Time!
  totalPrice: Float!
  products: [OrderedProduct!]!
}

type OrderedProduct {
  id: String!
  name: String!
  description: String!
  price: Float!
  quantity: Int!
}

input OrderInput {
  accountId: String!
  products: [OrderProductInput!]!
}

input OrderProductInput {
  id: String!
  quantity: Int!
}
```

### Common Types
```graphql
input PaginationInput {
  skip: Int
  take: Int
}

scalar Time
```

## Queries

### accounts
Retrieves a list of accounts with optional pagination and filtering.

```graphql
accounts(pagination: PaginationInput, id: String): [Account!]!
```

Parameters:
- `pagination`: Optional pagination input
  - `skip`: Number of records to skip
  - `take`: Number of records to take
- `id`: Optional account ID to filter by

Returns:
- Array of Account objects
- Empty array if no accounts found

### products
Retrieves a list of products with optional pagination, search, and filtering.

```graphql
products(pagination: PaginationInput, query: String, id: String, ids: [String!]): [Product!]!
```

Parameters:
- `pagination`: Optional pagination input
  - `skip`: Number of records to skip
  - `take`: Number of records to take
- `query`: Optional search term to filter products
- `id`: Optional product ID to filter by
- `ids`: Optional array of product IDs to filter by

Returns:
- Array of Product objects
- Empty array if no products found
- For `id` or `ids` queries, returns only matching products

## Mutations

### createAccount
Creates a new account.

```graphql
createAccount(account: AccountInput!): Account
```

Parameters:
- `account`: Account input object
  - `name`: Account holder's name

Returns:
- Created Account object or null if creation fails

### createProduct
Creates a new product.

```graphql
createProduct(product: ProductInput!): Product
```

Parameters:
- `product`: Product input object
  - `name`: Product name
  - `description`: Product description
  - `price`: Product price (must be positive)

Returns:
- Created Product object or null if creation fails

### createOrder
Creates a new order for an account.

```graphql
createOrder(order: OrderInput!): Order
```

Parameters:
- `order`: Order input object
  - `accountId`: ID of the account placing the order
  - `products`: Array of products to order
    - `id`: Product ID
    - `quantity`: Quantity to order (must be positive)

Returns:
- Created Order object or null if creation fails

Error Responses:
- Account not found: `"account with ID {accountId} does not exist"`
- Product not found: `"one or more products in your order could not be found"`
- Invalid quantity: `"invalid quantity for product {productId}: quantity must be greater than 0"`
- General error: `"failed to create order, please try again"`

## Error Handling

The API implements comprehensive error handling with detailed messages and proper error classification. Errors are returned in the following format:

```json
{
  "errors": [
    {
      "message": "User-friendly error message",
      "path": ["path", "to", "field"],
      "extensions": {
        "code": "ERROR_CODE",
        "details": "Additional error context"
      }
    }
  ],
  "data": {
    "fieldName": null
  }
}
```

### Error Types

1. **Validation Errors**
   - Missing required fields
   - Invalid data types
   - Out of range values
   - Pagination limits exceeded

2. **Resource Errors**
   - Resource not found
   - Resource already exists
   - Resource unavailable
   - Resource state conflict

3. **Authorization Errors**
   - Invalid credentials
   - Insufficient permissions
   - Token expired
   - Invalid token

4. **System Errors**
   - Internal server error
   - Service unavailable
   - Database error
   - Network timeout

### Common Error Scenarios

#### Account Service
- `ACCOUNT_NOT_FOUND`: Account with specified ID does not exist
- `INVALID_ACCOUNT_NAME`: Account name is empty or invalid
- `ACCOUNT_EXISTS`: Account with same name already exists
- `INVALID_PAGINATION`: Pagination parameters exceed limits

#### Catalog Service
- `PRODUCT_NOT_FOUND`: Product with specified ID does not exist
- `INVALID_PRODUCT_PRICE`: Product price must be positive
- `INVALID_PRODUCT_NAME`: Product name is empty or invalid
- `INVALID_SEARCH_QUERY`: Search query is empty or invalid

#### Order Service
- `INVALID_ORDER`: Order validation failed
  - Empty order
  - Invalid quantities
  - Invalid product IDs
- `ORDER_TOTAL_EXCEEDED`: Order total exceeds maximum allowed
- `PRODUCT_UNAVAILABLE`: One or more products not available
- `INVALID_QUANTITY`: Product quantity must be positive

### Error Response Examples

1. **Validation Error**
```json
{
  "errors": [
    {
      "message": "Invalid pagination parameters: take must not exceed 100",
      "path": ["products"],
      "extensions": {
        "code": "INVALID_PAGINATION",
        "details": "Requested take: 150, Maximum allowed: 100"
      }
    }
  ]
}
```

2. **Resource Error**
```json
{
  "errors": [
    {
      "message": "Product not found",
      "path": ["createOrder"],
      "extensions": {
        "code": "PRODUCT_NOT_FOUND",
        "details": "Product ID: abc123 does not exist"
      }
    }
  ]
}
```

3. **System Error**
```json
{
  "errors": [
    {
      "message": "Internal server error",
      "path": ["createOrder"],
      "extensions": {
        "code": "INTERNAL_ERROR",
        "details": "Request ID: xyz789 for debugging"
      }
    }
  ]
}
```

### Error Handling Best Practices

1. **Client-Side**
   - Always check for error responses
   - Handle specific error codes appropriately
   - Display user-friendly error messages
   - Implement retry logic for transient errors
   - Log errors with context for debugging

2. **Error Recovery**
   - Implement proper error recovery strategies
   - Use exponential backoff for retries
   - Cache valid responses when appropriate
   - Maintain consistent UI state during errors
   - Provide clear user feedback

3. **Monitoring**
   - Track error frequencies
   - Monitor error patterns
   - Set up alerts for critical errors
   - Analyze error trends
   - Use request IDs for tracing

## Rate Limiting

Currently, there are no rate limits implemented on the API. However, it's recommended to:
- Use pagination for large result sets
- Implement client-side caching where appropriate
- Avoid excessive concurrent requests

## Best Practices

1. Always use pagination when fetching lists of items
2. Include only the fields you need in your queries
3. Handle errors gracefully on the client side
4. Use appropriate HTTP headers for caching
5. Consider implementing client-side caching for frequently accessed data
6. Test error scenarios with invalid inputs
7. Use the multiple product query (ids parameter) when fetching specific products
8. Handle empty results appropriately (empty arrays instead of null) 