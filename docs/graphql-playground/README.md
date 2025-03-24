# GraphQL Playground Testing Guide

This guide will help you test the Go Marketplace API using the GraphQL Playground interface.

## Getting Started

1. Ensure all services are running:
   ```bash
   docker-compose up --build
   ```

2. Open your browser and navigate to [http://localhost:8080/playground](http://localhost:8080/playground)

## Testing Queries

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

### Search Products
```graphql
query {
  products(query: "phone") {
    id
    name
    description
    price
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

### Create Order
```graphql
mutation {
  createOrder(order: {
    accountId: "your_account_id"
    products: [
      {
        id: "product_id_1"
        quantity: 2
      },
      {
        id: "product_id_2"
        quantity: 1
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

## Example Testing Scenarios

### 1. Complete Order Flow
1. Create an account
2. Create some products
3. Create an order with the products
4. Query the account to verify the order appears

```graphql
# 1. Create Account
mutation {
  createAccount(account: { name: "Test User" }) {
    id
    name
  }
}

# 2. Create Products
mutation {
  product1: createProduct(product: {
    name: "Test Product 1"
    description: "First test product"
    price: 29.99
  }) {
    id
  }
  
  product2: createProduct(product: {
    name: "Test Product 2"
    description: "Second test product"
    price: 39.99
  }) {
    id
  }
}

# 3. Create Order
mutation {
  createOrder(order: {
    accountId: "account_id_from_step_1"
    products: [
      { id: "product_id_1_from_step_2", quantity: 2 },
      { id: "product_id_2_from_step_2", quantity: 1 }
    ]
  }) {
    id
    totalPrice
    products {
      name
      price
      quantity
    }
  }
}

# 4. Verify Order in Account
query {
  accounts(id: "account_id_from_step_1") {
    name
    orders {
      id
      totalPrice
      products {
        name
        price
        quantity
      }
    }
  }
}
```

### 2. Product Search and Filtering
1. Create multiple products
2. Test search functionality
3. Test pagination

```graphql
# 1. Create Multiple Products
mutation {
  phone1: createProduct(product: {
    name: "iPhone 15"
    description: "Latest iPhone model"
    price: 999.99
  }) {
    id
  }
  
  phone2: createProduct(product: {
    name: "Samsung Galaxy S24"
    description: "Latest Samsung phone"
    price: 899.99
  }) {
    id
  }
  
  laptop: createProduct(product: {
    name: "MacBook Pro"
    description: "Powerful laptop"
    price: 1499.99
  }) {
    id
  }
}

# 2. Search Products
query {
  phones: products(query: "phone") {
    name
    description
    price
  }
  
  samsung: products(query: "samsung") {
    name
    description
    price
  }
}

# 3. Test Pagination
query {
  firstPage: products(pagination: { skip: 0, take: 2 }) {
    name
    price
  }
  
  secondPage: products(pagination: { skip: 2, take: 2 }) {
    name
    price
  }
}
```

## Error Handling

The API returns appropriate error messages when:
- Required fields are missing
- Invalid IDs are provided
- Business rules are violated (e.g., ordering non-existent products)

Example error response:
```json
{
  "errors": [
    {
      "message": "Product not found",
      "path": ["createOrder"]
    }
  ],
  "data": null
}
```

## Tips
1. Use the "Schema" tab in GraphQL Playground to explore available queries and mutations
2. Variables can be defined in the "Query Variables" panel
3. Headers can be added in the "HTTP Headers" panel
4. Each request can be shared via the "Share" button 