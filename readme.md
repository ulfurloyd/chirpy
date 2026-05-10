# Chirpy

Chirpy is a custom-built social media backend server written in Go.
It handles user authentication, chirps (short posts), and
administrative tools like health checks and metrics.

## Features
- RESTful API for managing users and chirps.
- User authentication using JWTs and Refresh Tokens.
- Integration with PostreSQL for persistent storage.
- Support for API versioning (e.g., `/api/v1/`).
- Health check and administrative metrics endpoints.


## Installation

### Prerequisites
- Go 1.21+
- PostreSQL

### Running the Server
1. Clone the repository
2. Set up your environment variables (DB connection string, JWT secrets, etc.)
3. Run the application:
```bash
go build -o out && ./out
```

The server will start on port `8082` by default.

## API Documentataion
- `GET /api/healthz` - Check if the server is alive.
- `POST /api/users` - Create a new user.
- `PUT /api/users` - Update a user.
- `POST /api/login` - Authenticate and receive tokens.
- `GET /api/chirps` - Fetch chirps.
- `POST /api/chirps` - Post a new chirp.
- `GET /api/chirps` - Fetch chirps.
- `POST /api/refresh` - Refresh the API.

> **_NOTE:_** This project was built as part of the [Boot.dev](https://www.boot.dev) DevOps Engineering learning path, in the [Learn HTTP Servers](https://www.boot.dev/courses/learn-http-servers-golang) course.
