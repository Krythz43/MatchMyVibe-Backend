# MatchMyVibe Backend

MatchMyVibe is a music-based dating app that matches users based on their music tastes from Spotify. This repository contains the backend API built with Go.

## Features

- Spotify authentication integration
- JWT-based authentication for all API requests
- User profile management
- Storing and retrieving user's music preferences from Spotify
- Real-time "currently playing" track updates
- Detailed last played song information with user activity tracking
- Dating profile with gender and preferences

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

5. Apply the migrations:
   ```
   psql -U postgres -d matchmyvibe -f internal/db/migrations/add_last_played_song.sql
   ```

6. Build and run the application:
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
  - Response: Full user profile including dating preferences, last played song and activity timestamp

- `PUT /api/profile` - Update the user's profile
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Request body: Any profile fields to update including the new fields
    ```json
    {
      "name": "John Doe",
      "birthdayInUnix": 631152000,
      "gender": "Man",  // Can be "Man", "Woman", or "Non-binary"
      "dating_preference": "Everyone"  // Can be "Men", "Women", or "Everyone"
    }
    ```
  - Response: Updated full user profile

- `PUT /api/profile/currently-playing` - Update the user's currently playing track
  - Headers:
    ```
    Authorization: Bearer <token>
    ```
  - Can be called in two ways:
    1. Without a request body (legacy): Fetches the currently playing track from Spotify
    2. With a request body: Updates with detailed last played song information
       ```json
       {
         "track": "Midnight Rain",
         "artist": "Taylor Swift",
         "uri": "spotify:track:4eKMqf9ZMSclDX7V9Ptg7x",
         "album": "Midnights (The Til Dawn Edition)",
         "album_uri": "spotify:album:1fnJ7k0bllNfL1kVdNVW1A",
         "duration": 174782,
         "context_title": "Midnights (The Til Dawn Edition)",
         "context_uri": "spotify:album:1fnJ7k0bllNfL1kVdNVW1A"
       }
       ```
  - Response:
    ```json
    {
      "currently_playing": "Midnight Rain - Taylor Swift",
      "last_played_song": {
        "track": "Midnight Rain",
        "artist": "Taylor Swift",
        "uri": "spotify:track:4eKMqf9ZMSclDX7V9Ptg7x",
        "album": "Midnights (The Til Dawn Edition)",
        "album_uri": "spotify:album:1fnJ7k0bllNfL1kVdNVW1A",
        "duration": 174782,
        "context_title": "Midnights (The Til Dawn Edition)",
        "context_uri": "spotify:album:1fnJ7k0bllNfL1kVdNVW1A"
      },
      "user_last_active_at": 1693245678
    }
    ```

## Recent Updates

- Added last played song functionality with detailed track information
- Added userLastActiveAt timestamp to track when a user last played a song (as UNIX timestamp)
- Added new profile fields:
  - `birthdayInUnix`: User's birthday as Unix timestamp
  - `gender`: User's gender (Man, Woman, or Non-binary)
  - `dating_preference`: User's dating preference (Men, Women, or Everyone)

## License

MIT
