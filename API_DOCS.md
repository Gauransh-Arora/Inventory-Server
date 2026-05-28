# Inventory Server API Documentation

A comprehensive guide to the HTTP API endpoints provided by the Inventory Server.

## 🌐 Base URL
All API endpoints (except the health check endpoint) are prefixed with:
```http
/api/v1
```

---

## 🔒 Authentication & Security

### Bearer Token Authentication
Most endpoints in this API require authentication. These endpoints are protected using JWT (JSON Web Tokens) with RSA-256 signing.

To access protected endpoints, you must include the access token in the request headers:
```http
Authorization: Bearer <access_token>
```

### Auth Middleware Error Responses
If the authentication check fails, the API returns a `401 Unauthorized` status code with one of the following JSON error payloads:

*   **Missing or Incorrect Authorization Header Format:**
    ```json
    {
      "error": "Bearer token required"
    }
    ```
*   **Invalid or Expired Access Token:**
    ```json
    {
      "error": "invalid or expired token"
    }
    ```
*   **Revoked/Denylisted Token (after logging out):**
    ```json
    {
      "error": "token has been revoked"
    }
    ```

---

## 📡 Public Endpoints

### Health Check
Check the server's availability and health status.

*   **Endpoint:** `GET /ping`
*   **Auth Required:** No
*   **Response (200 OK):**
    ```json
    {
      "message": "pong"
    }
    ```

---

## 🔑 Authentication Endpoints

### 1. Register User
Create a new user account.

*   **Endpoint:** `POST /api/v1/auth/register`
*   **Auth Required:** No
*   **Request Body (JSON):**
    | Field | Type | Binding Constraints | Description |
    | :--- | :--- | :--- | :--- |
    | `username` | `string` | Required | The user's unique username. |
    | `password` | `string` | Required, Min length 8 | The password for the new account. |
    
    *Example:*
    ```json
    {
      "username": "john_doe",
      "password": "Password123!"
    }
    ```

*   **Success Response (201 Created):**
    ```json
    {
      "id": "e4a2d815-56fa-4654-8c88-29cf6cf8bf27",
      "username": "john_doe",
      "message": "User created successfully"
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If the request body fails validation (invalid username or password less than 8 characters).
        ```json
        {
          "error": "invalid input"
        }
        ```
    *   **409 Conflict:** If a user with the provided username is already registered.
        ```json
        {
          "error": "user already exists"
        }
        ```

---

### 2. Login
Authenticate user credentials and issue an Access Token and a Refresh Token.

*   **Endpoint:** `POST /api/v1/auth/login`
*   **Auth Required:** No
*   **Request Body (JSON):**
    | Field | Type | Binding Constraints | Description |
    | :--- | :--- | :--- | :--- |
    | `username` | `string` | Required | User's username. |
    | `password` | `string` | Required, Min length 8 | User's password. |
    
    *Example:*
    ```json
    {
      "username": "john_doe",
      "password": "Password123!"
    }
    ```

*   **Success Response (200 OK):**
    ```json
    {
      "access_token": "eyJhbGciOiJSUzI1NiIs...",
      "refresh_token": "a1b2c3d4e5f6g7h8..."
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If required login credentials are not provided.
        ```json
        {
          "error": "Credentials required"
        }
        ```
    *   **401 Unauthorized:** If the credentials are invalid, or if the account is deactivated.
        ```json
        {
          "error": "invalid credentials"
        }
        ```
        *or*
        ```json
        {
          "error": "account disabled"
        }
        ```

---

### 3. Refresh Token
Issue a new Access Token and rotate the Refresh Token.

*   **Endpoint:** `POST /api/v1/auth/refresh`
*   **Auth Required:** No
*   **Request Body (JSON):**
    | Field | Type | Binding Constraints | Description |
    | :--- | :--- | :--- | :--- |
    | `refresh_token` | `string` | Required | The active refresh token. |
    
    *Example:*
    ```json
    {
      "refresh_token": "a1b2c3d4e5f6g7h8..."
    }
    ```

*   **Success Response (200 OK):**
    *Returns a rotated token pair. The old refresh token will be revoked.*
    ```json
    {
      "access_token": "eyJhbGciOiJSUzI1NiIs...",
      "refresh_token": "z9y8x7w6v5u4t3s2..."
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If the refresh token is missing.
        ```json
        {
          "error": "Refresh token required"
        }
        ```
    *   **401 Unauthorized:** If the refresh token is expired or invalid.
        ```json
        {
          "error": "invalid refresh token"
        }
        ```

---

## 🔓 Protected Endpoints

All endpoints below require authentication.

### 4. Logout
Revoke the current session, delete active refresh tokens, and add the access token to the JWT denylist.

*   **Endpoint:** `POST /api/v1/logout`
*   **Auth Required:** Yes
*   **Success Response (200 OK):**
    ```json
    {
      "message": "Logged out successfully"
    }
    ```
*   **Error Responses:**
    *   **500 Internal Server Error:** If there is a backend failure revoking tokens.
        ```json
        {
          "error": "failed to logout"
        }
        ```

---

### 5. Create Product (Bulk)
Create new product records in the system. Accepts a JSON object containing an array of products.

*   **Endpoint:** `POST /api/v1/products`
*   **Auth Required:** Yes
*   **Request Body (JSON):**
    | Field | Type | Description |
    | :--- | :--- | :--- |
    | `products` | `array` | A list of product objects to insert. |
    
    Each object within the `products` array must contain:
    | Field | Type | Description |
    | :--- | :--- | :--- |
    | `icode` | `integer` | ID of the product from legacy system. |
    | `item_name` | `string` | Name of the product. |
    | `batch_no` | `integer` | Batch number of the product. |
    | `mrp` | `number` | Maximum Retail Price of the product. |
    | `barcode` | `string` | Barcode identifying the product (can be duplicate across different batches). |
    
    *Example:*
    ```json
    {
      "products": [
        {
          "icode": 1,
          "item_name": "Premium Soap",
          "batch_no": 1,
          "mrp": 45.00,
          "barcode": "8901030752538"
        },
        {
          "icode": 2,
          "item_name": "Premium Shampoo",
          "batch_no": 2,
          "mrp": 120.00,
          "barcode": "8901030752539"
        }
      ]
    }
    ```

*   **Success Response (201 Created):**
    ```json
    {
      "message": "product created"
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If the request JSON cannot be parsed.
        ```json
        {
          "error": "invalid character..."
        }
        ```
    *   **500 Internal Server Error:** If the database operation fails (e.g., if a product already exists, the entire batch will be rolled back).
        ```json
        {
          "error": "product already exists (icode: 1, name: Premium Soap)"
        }
        ```

---

### 6. Get All Products
Retrieve a paginated list of products. All query parameters are optional.

*   **Endpoint:** `GET /api/v1/products`
*   **Auth Required:** Yes
*   **Query Parameters:**
    | Parameter | Type | Default | Description | Example |
    | :--- | :--- | :--- | :--- | :--- |
    | `icode` | `integer` | — | Filter by legacy product icode. | `?icode=1001` |
    | `page` | `integer` | `1` | Page number (1-indexed). | `?page=2` |
    | `page_size` | `integer` | `20` | Number of records per page. Max is `100`. | `?page_size=50` |

    *Examples:*
    ```http
    GET /api/v1/products
    GET /api/v1/products?page=2&page_size=50
    GET /api/v1/products?icode=1001
    GET /api/v1/products?icode=1001&page=1&page_size=10
    ```

*   **Success Response (200 OK):**
    ```json
    {
      "data": [
        {
          "id": 1,
          "icode": 1001,
          "item_name": "Premium Soap",
          "batch_no": 1,
          "mrp": 45.00,
          "barcode": "8901030752538"
        }
      ],
      "page": 1,
      "page_size": 20,
      "total_count": 142,
      "total_pages": 8
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If a query parameter has an invalid value.
        ```json
        { "error": "invalid value for 'icode', must be an integer" }
        ```
        ```json
        { "error": "invalid value for 'page', must be a positive integer" }
        ```
    *   **500 Internal Server Error:** If the database operation fails.
        ```json
        {
          "error": "internal server error"
        }
        ```

---

### 7. Get Product by Barcode
Retrieve a list of products that match a specific barcode.

*   **Endpoint:** `GET /api/v1/products/:barcode`
*   **Auth Required:** Yes
*   **Parameters:**
    *   `barcode` (Path parameter, string): The barcode to search for.
*   **Success Response (200 OK):**
    *Returns a list of matching product objects (typically containing one or more matches).*
    ```json
    [
      {
        "id": 1,
        "icode": 1,
        "item_name": "Premium Soap",
        "batch_no": 1,
        "mrp": 45.00,
        "barcode": "8901030752538"
      }
    ]
    ```

*   **Error Responses:**
    *   **404 Not Found:** If no products are found matching the provided barcode.
        ```json
        {
          "message": "no products found"
        }
        ```
    *   **500 Internal Server Error:** If the database operation fails.
        ```json
        {
          "error": "internal server error"
        }
        ```

---

### 8. Create Inventory Log
Record a change in inventory levels for a product.

*   **Endpoint:** `POST /api/v1/logs`
*   **Auth Required:** Yes
*   **Request Body (JSON):**
    | Field | Type | Description |
    | :--- | :--- | :--- |
    | `product_id` | `integer` | ID of the product. |
    | `quantity` | `integer` | Quantity added (positive) or removed (negative). |
    
    *Example:*
    ```json
    {
      "product_id": 1,
      "quantity": 25
    }
    ```

*   **Success Response (200 OK):**
    ```json
    {
      "message": "log created"
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If the request JSON cannot be parsed.
        ```json
        {
          "error": "invalid character..."
        }
        ```
    *   **500 Internal Server Error:** If the database operation fails.
        ```json
        {
          "error": "database error..."
        }
        ```

---

### 9. Get All Inventory Logs
Retrieve a paginated list of inventory logs, ordered by creation date descending. All query parameters are optional and can be combined freely.

*   **Endpoint:** `GET /api/v1/logs`
*   **Auth Required:** Yes
*   **Query Parameters:**
    | Parameter | Type | Default | Description | Example |
    | :--- | :--- | :--- | :--- | :--- |
    | `updated` | `boolean` | — | Filter by export status. `false` = not yet exported, `true` = already exported. | `?updated=false` |
    | `product_id` | `integer` | — | Filter logs for a specific product. | `?product_id=5` |
    | `date_from` | `YYYY-MM-DD` | — | Logs created on or after this date (inclusive, start of day UTC). | `?date_from=2026-05-01` |
    | `date_to` | `YYYY-MM-DD` | — | Logs created on or before this date (inclusive, end of day UTC). | `?date_to=2026-05-26` |
    | `page` | `integer` | `1` | Page number (1-indexed). | `?page=2` |
    | `page_size` | `integer` | `20` | Number of records per page. Max is `100`. | `?page_size=50` |

    *Examples:*
    ```http
    GET /api/v1/logs
    GET /api/v1/logs?updated=false
    GET /api/v1/logs?product_id=5&updated=false
    GET /api/v1/logs?date_from=2026-05-01&date_to=2026-05-26&page=1&page_size=100
    GET /api/v1/logs?updated=false&product_id=3&date_from=2026-05-26
    ```

*   **Success Response (200 OK):**
    ```json
    {
      "data": [
        {
          "id": 1,
          "product_id": 1,
          "quantity": 25,
          "updated": false,
          "created_at": "2026-05-25T14:30:15Z"
        }
      ],
      "page": 1,
      "page_size": 20,
      "total_count": 540,
      "total_pages": 27
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If a query parameter has an invalid value or format.
        ```json
        { "error": "invalid value for 'updated', must be true or false" }
        ```
        ```json
        { "error": "invalid value for 'product_id', must be an integer" }
        ```
        ```json
        { "error": "invalid 'date_from' format, use YYYY-MM-DD" }
        ```
        ```json
        { "error": "invalid value for 'page', must be a positive integer" }
        ```
    *   **500 Internal Server Error:** If the database operation fails.
        ```json
        {
          "error": "internal server error"
        }
        ```

---

### 10. Mark Inventory Logs as Updated
Mark inventory logs as exported/updated. Used by the end-of-day operator after storing logs to the local machine.

*   **Endpoint:** `PATCH /api/v1/logs`
*   **Auth Required:** Yes
*   **Request Body (JSON, optional):**
    | Field | Type | Description |
    | :--- | :--- | :--- |
    | `ids` | `[]integer` | List of specific log IDs to mark as updated. If omitted or empty, **all** currently un-updated logs are marked. |

    *Example — mark specific logs:*
    ```json
    {
      "ids": [1, 4, 7]
    }
    ```

    *Example — mark all pending logs (bulk):*
    ```json
    {}
    ```

*   **Success Response (200 OK) — bulk:**
    ```json
    {
      "message": "all pending logs marked as updated"
    }
    ```

*   **Success Response (200 OK) — specific IDs:**
    ```json
    {
      "message": "specified logs marked as updated",
      "ids": [1, 4, 7]
    }
    ```

*   **Error Responses:**
    *   **400 Bad Request:** If the request body cannot be parsed.
        ```json
        {
          "error": "invalid character..."
        }
        ```
    *   **500 Internal Server Error:** If the database operation fails.
        ```json
        {
          "error": "internal server error"
        }
        ```
