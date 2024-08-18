# Tarantool KV API

This project provides a simple RESTful API for storing and retrieving key-value pairs using Tarantool as the backend storage. The API is secured with JWT authentication.

## Table of Contents
- [Installation](#installation)
- [API Documentation](#api-documentation)
  - [Login](#login)
  - [Write Data](#write-data)
  - [Read Data](#read-data)

## Installation

### Prerequisites

- Docker and Docker Compose installed on your machine.
- Git installed.

### Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/DmitryNai/RESTful-API-practice
    cd RESTful-API-practice
    ```

2. Start the services with Docker Compose:

    ```bash
    sudo docker-compose up --build
    ```

3. The API will be available at `http://localhost:8080`.

## API Documentation

### Login

**Endpoint:** `/api/login`

**Method:** `POST`

**Description:** Authenticates the user and returns a JWT token.

**Request Body:**

```json
{
  "username": "admin",
  "password": "presale"
}
```

**Response:**

```json

{
  "token": "your.jwt.token"
}
```

Error Codes:

    401 Unauthorized - Invalid credentials.

### Write Data

**Endpoint:** `/api/write`

**Method:** `POST`

**Description:** Stores key-value pairs in the Tarantool database.

**Headers:**

    Authorization: Bearer your.jwt.token

**Request Body:**

```json

{
  "data": {
    "key1": "value1",
    "key2": "value2",
    "key3": 1
  }
}
```

**Response:**

```json

{
  "status": "success"
}
```

**Error Codes:**

    400 Bad Request - Malformed request.
    500 Internal Server Error - Failed to write some keys.

### Read Data

**Endpoint:** `/api/read`

**Method:** `POST`

**Description:** Retrieves values for specified keys from the Tarantool database.

**Headers:**

    Authorization: Bearer your.jwt.token

**Request Body:**

```json

{
  "keys": ["key1", "key2", "key3"]
}
```

**Response:**

```json

{
  "data": {
    "key1": "value1",
    "key2": "value2",
    "key3": 1
  }
}
```

**Error Codes:**

    400 Bad Request - Malformed request.
    404 Not Found - Key not found.
