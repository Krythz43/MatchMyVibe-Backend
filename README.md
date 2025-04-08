# MatchMyVibe - Music-Based Dating App Backend

A dating application that matches users based on their music preferences using Spotify API integration with Supabase for authentication and database storage.

## Project Status

- [x] Supabase project created
- [x] Database schema configured with necessary tables
- [x] Spotify authentication flow implemented
- [x] Basic frontend with "Sign in with Spotify" functionality
- [x] Fetching and displaying user's top artists and tracks
- [ ] Implement user matching algorithm based on music preferences
- [ ] Add chat functionality
- [ ] Create mobile (iOS) application
- [ ] Add profile customization features

## Tech Stack

- **Backend:** Node.js, Express
- **Database:** PostgreSQL (via Supabase)
- **Authentication:** Supabase Auth with Spotify OAuth
- **Frontend:** Plain HTML, CSS, JavaScript (currently)
- **API Integration:** Spotify Web API

## Project Structure

```
MatchMyVibe-Backend/
├── public/                  # Static frontend files
│   ├── index.html           # Main HTML page with login button
│   └── style.css            # CSS styles
├── src/                     # Source files (not fully utilized yet)
├── .env                     # Environment variables
├── .gitignore               # Git ignore file
├── index.js                 # Main Express server
├── package.json             # Project dependencies
└── README.md                # Project documentation
```

## Database Schema

The application uses the following database tables:

1. **profiles** - Stores user profile information
   - id (UUID, primary key)
   - spotify_id (text, unique)
   - display_name (text)
   - avatar_url (text)
   - created_at (timestamp)
   - updated_at (timestamp)

2. **top_artists** - Stores user's top artists from Spotify
   - id (UUID, primary key)
   - profile_id (UUID, foreign key to profiles)
   - artist_id (text)
   - artist_name (text)
   - popularity (integer)
   - genres (text array)
   - image_url (text)
   - time_range (text)
   - timestamp (timestamp)

3. **top_tracks** - Stores user's top tracks from Spotify
   - id (UUID, primary key)
   - profile_id (UUID, foreign key to profiles)
   - track_id (text)
   - track_name (text)
   - artist_name (text)
   - album_name (text)
   - popularity (integer)
   - image_url (text)
   - time_range (text)
   - timestamp (timestamp)

## Setup Instructions

### Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- Supabase account
- Spotify Developer account

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Supabase configuration
SUPABASE_URL=your_supabase_project_url
SUPABASE_ANON_KEY=your_supabase_anon_key

# Spotify API credentials
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
SPOTIFY_REDIRECT_URI=http://localhost:3000/callback
```

### Supabase Setup

1. Create a new Supabase project
2. Set up Spotify OAuth provider in Supabase Auth settings
3. Configure the database schema by running the SQL migrations

### Installation and Running

1. Clone the repository
2. Install dependencies:
   ```
   npm install
   ```
3. Start the development server:
   ```
   npm run dev
   ```
4. Open http://localhost:3000 in your browser

## Next Steps

1. Implement user matching algorithm based on music taste
2. Add swiping functionality for potential matches
3. Develop iOS mobile application
4. Add chat functionality between matched users
5. Enhance user profiles with additional customization options

## License

ISC
