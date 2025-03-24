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

The API follows standard GraphQL error handling practices with user-friendly messages. Errors are returned in the following format:

```json
{
  "errors": [
    {
      "message": "User-friendly error message",
      "path": ["path", "to", "field"]
    }
  ],
  "data": {
    "fieldName": null
  }
}
```

Common error scenarios:
- `NOT_FOUND`: Resource not found (account, product)
- `INVALID_INPUT`: Invalid input (quantity, price)
- `INTERNAL_ERROR`: Internal server error

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