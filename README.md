# PAF - Go Authentication REST API Service

This project is a backend REST API service built with Go, providing user authentication and authorization functionalities using JSON Web Tokens (JWT). It includes features for user signup, signin, token generation, token validation, token revocation, and token refresh. MongoDB is used as the persistent data store for user credentials and revoked tokens.

## Features

*   User Signup (Email & Password)
*   User Signin (Returns JWT Access Token & Refresh Token)
*   Token Authorization (Bearer Token in Authorization header)
*   Token Expiry Checks
*   Token Revocation (Persistent in MongoDB with TTL)
*   Token Refresh Mechanism

## Prerequisites

*   Docker and Docker Compose must be installed on your system.
*   Git
(Optional) For local development outside of Docker:
*   Go (version 1.23 or higher, matching the Dockerfile)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/deepaksingh26/PAF.git
    cd PAF
    ```

2.  **(Optional) For local development outside of Docker:**
    If you plan to develop or run the Go application natively (without Docker), ensure Go dependencies are fetched:
    ```bash
    go mod tidy
    ```

## Running the Service

This project is configured to run using Docker and Docker Compose, which handles both the Go application and the MongoDB database.

### Starting the Service

From the project root directory, run the following command:

To build and run the application and a MongoDB database in containers.

1.  Navigate to the project root directory.
2.  Run the following command:
    ```bash
    docker compose up
    ```
    To run in detached mode (in the background):
    ```bash
    docker compose up -d
    ```
    To rebuild the Go application image if you've made code changes:
    ```bash
    docker compose up --build
    ```

You should see output similar to:
```
Successfully connected to MongoDB and database initialized.
Starting server on :8080
```
The service will be running on `http://localhost:8080`.

## API Endpoints & Testing with cURL

Replace placeholders like `<YOUR_EMAIL>`, `<YOUR_PASSWORD>`, `<ACCESS_TOKEN>`, and `<REFRESH_TOKEN>` with actual values obtained from previous steps.

### 1. Sign Up

Creates a new user.

```bash
curl -X POST \
  http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "<YOUR_EMAIL>",
    "password": "<YOUR_PASSWORD>"
  }'
```
**Expected Success Response:** `{"message":"User created successfully"}` (or tokens if implemented)

### 2. Sign In

Authenticates a user and returns an access token and a refresh token.

```bash
curl -X POST \
  http://localhost:8080/signin \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "<YOUR_EMAIL>",
    "password": "<YOUR_PASSWORD>"
  }'
```
**Expected Success Response:** `{"token":"<ACCESS_TOKEN>","refresh_token":"<REFRESH_TOKEN>"}`
*(Save the `token` and `refresh_token` for subsequent requests)*

### 3. Access a Protected Route

Requires a valid access token in the Authorization header.

```bash
curl -X GET \
  http://localhost:8080/protected \
  -H 'Authorization: Bearer <ACCESS_TOKEN>'
```
**Expected Success Response:** `{"message":"Hello <YOUR_EMAIL>! You have accessed a protected route."}`

### 4. Refresh an Access Token

Uses a refresh token to obtain a new access token.

```bash
curl -X POST \
  http://localhost:8080/refresh \
  -H 'Content-Type: application/json' \
  -d '{
    "refresh_token": "<REFRESH_TOKEN>"
  }'
```
**Expected Success Response:** `{"token":"<NEW_ACCESS_TOKEN>"}`

### 5. Revoke Token(s)

Revokes the current access token (sent in header) and optionally a refresh token (sent in body). Requires authentication.

*   **Revoke current access token only:**
    ```bash
    curl -X POST \
      http://localhost:8080/revoke \
      -H 'Content-Type: application/json' \
      -H 'Authorization: Bearer <ACCESS_TOKEN>' \
      -d '{}'
    ```

*   **Revoke current access token AND a specific refresh token:**
    ```bash
    curl -X POST \
      http://localhost:8080/revoke \
      -H 'Content-Type: application/json' \
      -H 'Authorization: Bearer <ACCESS_TOKEN>' \
      -d '{
        "refresh_token": "<REFRESH_TOKEN_TO_REVOKE>"
      }'
    ```
**Expected Success Response:** `{"message":"Token(s) revoked successfully"}`

After revoking, try using the revoked token(s) with the `/protected` or `/refresh` endpoints; you should receive a `401 Unauthorized` error.

## Project Structure

*   `main.go`: Entry point of the application, sets up routes.
*   `database.go`: Handles MongoDB connection and initialization.
*   `models.go`: Defines data structures (User, RevokedToken) and database interaction logic.
*   `tokens.go`: JWT generation, validation, and revocation logic.
*   `middleware.go`: Authentication middleware.
*   `utils.go`: Utility functions like `respondJSON`.
*   `constants.go`: Defines constants like `minute` and `hour`.
<!-- *   `.env`: (To be created by user) Stores environment-specific configurations. -->
*   `go.mod`, `go.sum`: Go module files.