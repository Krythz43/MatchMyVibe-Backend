# MatchMyVibe Backend

MatchMyVibe is a music-based dating app that matches users based on their music tastes from Spotify. This repository contains the backend API built with Go.

## Features

- Spotify authentication integration
- JWT-based authentication for all API requests
- User profile management
- Storing and retrieving user's music preferences from Spotify
- Real-time "currently playing" track updates

## Tech Stack

- Go (Golang)
- Gin web framework
- PostgreSQL database
- JWT for authentication
- Spotify Web API integration

## Project Structure

```
├── cmd/
│   └── api/               # Main entry point for the API
├── config/                # Configuration files
├── internal/
│   ├── auth/              # Authentication logic
│   ├── db/                # Database access and models
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # Middleware components
│   ├── models/            # Data models
│   └── spotify/           # Spotify API integration
```

## Getting Started

### Prerequisites

- Go 1.19+
- PostgreSQL
- Spotify Developer API credentials

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/matchmyvibe/backend.git
   cd backend
   ```

2. Copy the example environment file and fill in your values:
   ```
   cp .env.example .env
   ```

3. Set up the PostgreSQL database:
   ```
   psql -U postgres -c "CREATE DATABASE matchmyvibe;"
   ```

4. Run the database schema:
   ```
   psql -U postgres -d matchmyvibe -f internal/db/schema.sql
   ```

5. Build and run the application:
   ```
   go build -o matchmyvibe-backend ./cmd/api
   ./matchmyvibe-backend
   ```

## API Endpoints

### Authentication

- `POST /auth/spotify` - Authenticate with Spotify credentials
  - Request body:
    ```json
    {
      "spotify_uri": "spotify:user:1234567890",
      "access_token": "spotify_access_token",
      "refresh_token": "spotify_refresh_token",
      "expiry_date": "2023-04-16T12:00:00Z"
    }
    ```
  - Response:
    ```json
    {
      "token": "jwt_token",
      "user": { ... },
      "is_new_user": true|false
    }
    ```

- `POST /auth/refresh` - Refresh JWT token
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Response:
    ```json
    {
      "token": "new_jwt_token"
    }
    ```

### User Profile

- `GET /api/profile` - Get the user's profile
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Response: Full user profile

- `PUT /api/profile` - Update the user's profile
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Request body: Any profile fields to update
  - Response: Updated full user profile

- `PUT /api/profile/currently-playing` - Update the user's currently playing track
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Response:
    ```json
    {
      "currently_playing": "Track Name - Artist Name"
    }
    ```

## License

MIT
