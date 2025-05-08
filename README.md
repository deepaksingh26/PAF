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

*   Go (version 1.18 or higher recommended)
*   MongoDB (running instance, local or Atlas)
*   Git

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd PAF
    ```

2.  **Install dependencies:**
    The project uses Go modules. Dependencies will be fetched automatically when you build or run the project. You can also run:
    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables:**
    Create a `.env` file in the root of the project with the following content.

    ```env
    JWT_SECRET_KEY=""
    MONGO_URI=""
    MONGO_DATABASE_NAME=""
    ```

## Running the Service

Ensure your MongoDB instance is running and accessible.

From the project root directory, run the following command:

```bash
go run .
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
*   `.env`: (To be created by user) Stores environment-specific configurations.
*   `go.mod`, `go.sum`: Go module files.