package models

import (
	"time"

	"github.com/google/uuid"
)

// LastPlayedSong represents a user's last played song from Spotify
type LastPlayedSong struct {
	Track        string `json:"track" db:"track"`
	Artist       string `json:"artist" db:"artist"`
	URI          string `json:"uri" db:"uri"`
	Album        string `json:"album" db:"album"`
	AlbumURI     string `json:"album_uri" db:"album_uri"`
	Duration     int    `json:"duration" db:"duration"`
	ContextTitle string `json:"context_title" db:"context_title"`
	ContextURI   string `json:"context_uri" db:"context_uri"`
}

// User represents the main user profile
type User struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	SpotifyURI       string          `json:"spotify_uri" db:"spotify_uri"`
	AccessToken      string          `json:"-" db:"access_token"`
	RefreshToken     string          `json:"-" db:"refresh_token"`
	TokenExpiry      time.Time       `json:"-" db:"token_expiry"`
	Name             *string         `json:"name" db:"name"`
	UniversityName   *string         `json:"university_name" db:"university_name"`
	Work             *WorkProfile    `json:"work" db:"work"`
	HomeTown         *string         `json:"home_town" db:"home_town"`
	Height           *string         `json:"height" db:"height"`
	Age              *string         `json:"age" db:"age"`
	Zodiac           *string         `json:"zodiac" db:"zodiac"`
	CurrentlyPlaying *string         `json:"currently_playing" db:"currently_playing"`
	LastPlayedSong   *LastPlayedSong `json:"last_played_song" db:"last_played_song"`
	UserLastActiveAt *int64          `json:"user_last_active_at" db:"user_last_active_at"`
	BirthdayInUnix   *int64          `json:"birthdayInUnix" db:"birthdayInUnix"`
	Gender           *string         `json:"gender" db:"gender"`
	DatingPreference *string         `json:"dating_preference" db:"dating_preference"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// WorkProfile represents a user's work information
type WorkProfile struct {
	Company *string `json:"company" db:"company"`
	Role    *string `json:"role" db:"role"`
}

// Image represents a user's profile image
type Image struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Data   []byte    `json:"data" db:"data"`
}

// Interest represents a user's interest
type Interest struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Name   string    `json:"name" db:"name"`
}

// InterestRating represents a user's rating of an interest
type InterestRating struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Name   string    `json:"name" db:"name"`
	Rating int       `json:"rating" db:"rating"`
}

// Prompt represents a user's prompt response
type Prompt struct {
	ID       uuid.UUID `json:"id" db:"id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Question string    `json:"question" db:"question"`
	Answer   string    `json:"answer" db:"answer"`
}

// Artist represents a Spotify artist
type Artist struct {
	ID       uuid.UUID `json:"id" db:"id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Name     string    `json:"name" db:"name"`
	Uri      string    `json:"uri" db:"uri"`
	ImageURL *string   `json:"image_url" db:"image_url"`
}

// Song represents a Spotify song
type Song struct {
	ID       uuid.UUID `json:"id" db:"id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Name     string    `json:"name" db:"name"`
	Artist   string    `json:"artist" db:"artist"`
	Uri      string    `json:"uri" db:"uri"`
	ImageURL *string   `json:"image_url" db:"image_url"`
}

// Playlist represents a Spotify playlist
type Playlist struct {
	ID       uuid.UUID `json:"id" db:"id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Name     string    `json:"name" db:"name"`
	Uri      string    `json:"uri" db:"uri"`
	ImageURL *string   `json:"image_url" db:"image_url"`
}

// UserProfile represents the complete user profile to be returned by the API
type UserProfile struct {
	ID               uuid.UUID       `json:"id"`
	Name             *string         `json:"name"`
	UniversityName   *string         `json:"university_name"`
	Work             *WorkProfile    `json:"work"`
	HomeTown         *string         `json:"home_town"`
	Height           *string         `json:"height"`
	Age              *string         `json:"age"`
	Zodiac           *string         `json:"zodiac"`
	Images           [][]byte        `json:"images"`
	Interests        []string        `json:"interests"`
	InterestRating   map[string]int  `json:"interest_rating"`
	Prompts          []Prompt        `json:"prompts"`
	TopArtists       []Artist        `json:"top_artists"`
	TopSongs         []Song          `json:"top_songs"`
	SavedPlaylists   []Playlist      `json:"saved_playlists"`
	CurrentlyPlaying *string         `json:"currently_playing"`
	LastPlayedSong   *LastPlayedSong `json:"last_played_song"`
	UserLastActiveAt *int64          `json:"user_last_active_at"`
	BirthdayInUnix   *int64          `json:"birthdayInUnix"`
	Gender           *string         `json:"gender"`
	DatingPreference *string         `json:"dating_preference"`
}
