# GraphQL Playground Testing Guide

This guide will help you test the Go Marketplace API using the GraphQL Playground interface.

## Getting Started

1. Ensure all services are running:
   ```bash
   docker-compose up --build
   ```

2. Open your browser and navigate to [http://localhost:8080/playground](http://localhost:8080/playground)

## Testing Queries

### List All Products
```graphql
# Get all products
query {
  products {
    id
    name
    description
    price
  }
}
```

Success Response:
```json
{
  "data": {
    "products": [
      {
        "id": "product1",
        "name": "iPhone 15",
        "description": "Latest iPhone model",
        "price": 999.99
      }
    ]
  }
}
```

No Products Response:
```json
{
  "data": {
    "products": []
  }
}
```

### Search Products
```graphql
# Search by text
query {
  products(query: "phone") {
    id
    name
    description
    price
  }
}
```

Success Response:
```json
{
  "data": {
    "products": [
      {
        "id": "product1",
        "name": "iPhone 15",
        "description": "Latest iPhone model",
        "price": 999.99
      }
    ]
  }
}
```

No Matches Response:
```json
{
  "data": {
    "products": []
  }
}
```

### Get Product by ID
```graphql
query {
  products(id: "your_product_id") {
    id
    name
    description
    price
  }
}
```

Success Response:
```json
{
  "data": {
    "products": [
      {
        "id": "your_product_id",
        "name": "iPhone 15",
        "description": "Latest iPhone model",
        "price": 999.99
      }
    ]
  }
}
```

Product Not Found Response:
```json
{
  "data": {
    "products": []
  }
}
```

### Get Multiple Products by IDs
```graphql
query {
  products(ids: ["id1", "id2", "id3"]) {
    id
    name
    description
    price
  }
}
```

Success Response:
```json
{
  "data": {
    "products": [
      {
        "id": "id1",
        "name": "iPhone 15",
        "description": "Latest iPhone model",
        "price": 999.99
      },
      {
        "id": "id2",
        "name": "Samsung Galaxy S24",
        "description": "Latest Samsung phone",
        "price": 899.99
      }
    ]
  }
}
```

Some Products Not Found Response:
```json
{
  "data": {
    "products": [
      {
        "id": "id1",
        "name": "iPhone 15",
        "description": "Latest iPhone model",
        "price": 999.99
      }
    ]
  }
}
```

### List All Accounts
```graphql
query {
  accounts {
    id
    name
    orders {
      id
      totalPrice
      createdAt
    }
  }
}
```

Success Response:
```json
{
  "data": {
    "accounts": [
      {
        "id": "account1",
        "name": "John Doe",
        "orders": []
      }
    ]
  }
}
```

No Accounts Response:
```json
{
  "data": {
    "accounts": []
  }
}
```

## Testing Mutations

### Create Account
```graphql
mutation {
  createAccount(account: {
    name: "John Doe"
  }) {
    id
    name
  }
}
```

Success Response:
```json
{
  "data": {
    "createAccount": {
      "id": "account1",
      "name": "John Doe"
    }
  }
}
```

Invalid Input Response:
```json
{
  "errors": [
    {
      "message": "name is required",
      "path": ["createAccount"]
    }
  ],
  "data": {
    "createAccount": null
  }
}
```

### Create Product
```graphql
mutation {
  createProduct(product: {
    name: "iPhone 15"
    description: "Latest iPhone model"
    price: 999.99
  }) {
    id
    name
    description
    price
  }
}
```

Success Response:
```json
{
  "data": {
    "createProduct": {
      "id": "product1",
      "name": "iPhone 15",
      "description": "Latest iPhone model",
      "price": 999.99
    }
  }
}
```

Invalid Input Responses:
```json
{
  "errors": [
    {
      "message": "price cannot be negative",
      "path": ["createProduct"]
    }
  ],
  "data": {
    "createProduct": null
  }
}
```

```json
{
  "errors": [
    {
      "message": "name is required",
      "path": ["createProduct"]
    }
  ],
  "data": {
    "createProduct": null
  }
}
```

### Create Order
```graphql
mutation {
  createOrder(order: {
    accountId: "account1"
    products: [
      {
        id: "product1"
        quantity: 2
      },
      {
        id: "product2"
        quantity: 0  # Products with quantity <= 0 will be filtered out
      }
    ]
  }) {
    id
    totalPrice
    createdAt
    products {
      id
      name
      price
      quantity
    }
  }
}
```

Success Response:
```json
{
  "data": {
    "createOrder": {
      "id": "order1",
      "totalPrice": 1999.98,
      "createdAt": "2024-03-20T10:30:00Z",
      "products": [
        {
          "id": "product1",
          "name": "iPhone 15",
          "price": 999.99,
          "quantity": 2
        }
      ]
    }
  }
}
```

Error Responses:

1. Missing Account ID:
```json
{
  "errors": [
    {
      "message": "invalid parameter: accountId is required",
      "path": ["createOrder"]
    }
  ],
  "data": {
    "createOrder": null
  }
}
```

2. No Valid Products (all quantities <= 0):
```json
{
  "errors": [
    {
      "message": "invalid parameter: order must contain at least one product with quantity greater than 0",
      "path": ["createOrder"]
    }
  ],
  "data": {
    "createOrder": null
  }
}
```

3. Missing Product ID:
```json
{
  "errors": [
    {
      "message": "invalid parameter: product ID is required at index 0",
      "path": ["createOrder"]
    }
  ],
  "data": {
    "createOrder": null
  }
}
```

4. Account Not Found:
```json
{
  "errors": [
    {
      "message": "not found: account with ID account1 does not exist",
      "path": ["createOrder"]
    }
  ],
  "data": {
    "createOrder": null
  }
}
```

5. Product Not Found:
```json
{
  "errors": [
    {
      "message": "not found: one or more products in your order could not be found",
      "path": ["createOrder"]
    }
  ],
  "data": {
    "createOrder": null
  }
}
```

Note: Products with quantity <= 0 will be silently filtered out from the order. If all products have quantity <= 0, an error will be returned.

## Error Handling Best Practices

1. Always check for both `errors` and `data` in the response
2. Handle empty arrays gracefully (for queries that return lists)
3. Provide user-friendly error messages to end users
4. Implement retry logic for transient errors
5. Log detailed error information on the client side
6. Handle pagination errors by resetting to the first page
7. Validate input before sending to the server

Example Error Handling (JavaScript):
```javascript
const handleGraphQLResponse = (response) => {
  // Check for errors
  if (response.errors) {
    // Handle specific error types
    const error = response.errors[0];
    switch (true) {
      case error.message.includes('not found'):
        // Handle not found error
        showNotFoundMessage(error.message);
        break;
      case error.message.includes('invalid'):
        // Handle validation error
        showValidationError(error.message);
        break;
      default:
        // Handle unexpected errors
        showGeneralError('An unexpected error occurred');
    }
    return;
  }

  // Handle empty results
  if (Array.isArray(response.data?.products) && response.data.products.length === 0) {
    showNoResultsMessage();
    return;
  }

  // Process successful response
  processData(response.data);
};
```

## Tips
1. Use the "Schema" tab in GraphQL Playground to explore available queries and mutations
2. Variables can be defined in the "Query Variables" panel
3. Headers can be added in the "HTTP Headers" panel
4. Each request can be shared via the "Share" button
5. Use pagination for large result sets to improve performance
6. Include only the fields you need in your queries
7. Test error scenarios to ensure proper error handling
8. Use the Network tab to inspect actual HTTP responses
9. Save commonly used queries in the Playground
