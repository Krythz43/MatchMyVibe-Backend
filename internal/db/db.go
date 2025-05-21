package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/matchmyvibe/backend/internal/models"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// GetUserByID retrieves a user by their ID
func (db *DB) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	var workJSON []byte
	var lastPlayedSongJSON []byte

	query := `SELECT id, spotify_uri, access_token, refresh_token, token_expiry, 
			 name, university_name, work, home_town, height, age, zodiac, 
			 currently_playing, "birthdayInUnix", gender, dating_preference, 
			 last_played_song, user_last_active_at, created_at, updated_at 
			 FROM users WHERE id = $1`

	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.SpotifyURI, &user.AccessToken, &user.RefreshToken, &user.TokenExpiry,
		&user.Name, &user.UniversityName, &workJSON, &user.HomeTown, &user.Height, &user.Age, &user.Zodiac,
		&user.CurrentlyPlaying, &user.BirthdayInUnix, &user.Gender, &user.DatingPreference,
		&lastPlayedSongJSON, &user.UserLastActiveAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	// Parse work JSON if it exists
	if len(workJSON) > 0 {
		var work models.WorkProfile
		if err := json.Unmarshal(workJSON, &work); err != nil {
			return nil, err
		}
		user.Work = &work
	}

	// Parse last played song JSON if it exists
	if len(lastPlayedSongJSON) > 0 {
		var lastPlayedSong models.LastPlayedSong
		if err := json.Unmarshal(lastPlayedSongJSON, &lastPlayedSong); err != nil {
			return nil, err
		}
		user.LastPlayedSong = &lastPlayedSong
	}

	return &user, nil
}

// GetUserBySpotifyURI retrieves a user by their Spotify URI
func (db *DB) GetUserBySpotifyURI(spotifyURI string) (*models.User, error) {
	var user models.User
	var workJSON []byte
	var lastPlayedSongJSON []byte

	query := `SELECT id, spotify_uri, access_token, refresh_token, token_expiry, 
			 name, university_name, work, home_town, height, age, zodiac, 
			 currently_playing, "birthdayInUnix", gender, dating_preference,
			 last_played_song, user_last_active_at, created_at, updated_at 
			 FROM users WHERE spotify_uri = $1`

	fmt.Println("[DEBUG] Query:", query)
	fmt.Println("[DEBUG] Spotify URI:", spotifyURI)

	err := db.QueryRow(query, spotifyURI).Scan(
		&user.ID, &user.SpotifyURI, &user.AccessToken, &user.RefreshToken, &user.TokenExpiry,
		&user.Name, &user.UniversityName, &workJSON, &user.HomeTown, &user.Height, &user.Age, &user.Zodiac,
		&user.CurrentlyPlaying, &user.BirthdayInUnix, &user.Gender, &user.DatingPreference,
		&lastPlayedSongJSON, &user.UserLastActiveAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		fmt.Println("[DEBUG] Error fetching user by Spotify URI:", err)
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	// Parse work JSON if it exists
	if len(workJSON) > 0 {
		var work models.WorkProfile
		if err := json.Unmarshal(workJSON, &work); err != nil {
			return nil, err
		}
		user.Work = &work
	}

	// Parse last played song JSON if it exists
	if len(lastPlayedSongJSON) > 0 {
		var lastPlayedSong models.LastPlayedSong
		if err := json.Unmarshal(lastPlayedSongJSON, &lastPlayedSong); err != nil {
			return nil, err
		}
		user.LastPlayedSong = &lastPlayedSong
	}

	return &user, nil
}

// CreateUser creates a new user with Spotify authentication details
func (db *DB) CreateUser(spotifyURI, accessToken, refreshToken string, tokenExpiry time.Time) (*models.User, error) {
	id := uuid.New()
	query := `INSERT INTO users (id, spotify_uri, access_token, refresh_token, token_expiry, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`

	var user models.User
	err := db.QueryRow(query, id, spotifyURI, accessToken, refreshToken, tokenExpiry).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.SpotifyURI = spotifyURI
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	user.TokenExpiry = tokenExpiry

	return &user, nil
}

// UpdateUser updates user profile information
func (db *DB) UpdateUser(user *models.User) error {
	fmt.Printf("[DEBUG] UpdateUser - Updating user %s with BirthdayInUnix=%v, Gender=%v, DatingPreference=%v\n",
		user.ID, user.BirthdayInUnix, user.Gender, user.DatingPreference)

	workJSON, err := json.Marshal(user.Work)
	if err != nil {
		fmt.Printf("[ERROR] UpdateUser - Error marshaling work JSON: %v\n", err)
		return err
	}

	var lastPlayedSongJSON []byte
	if user.LastPlayedSong != nil {
		lastPlayedSongJSON, err = json.Marshal(user.LastPlayedSong)
		if err != nil {
			fmt.Printf("[ERROR] UpdateUser - Error marshaling last played song JSON: %v\n", err)
			return err
		}
	}

	query := `UPDATE users SET 
			 name = $1, university_name = $2, work = $3, home_town = $4, 
			 height = $5, zodiac = $6, currently_playing = $7, "birthdayInUnix" = $8,
			 gender = $9, dating_preference = $10, last_played_song = $11, user_last_active_at = $12, updated_at = NOW() 
			 WHERE id = $13`

	fmt.Printf("[DEBUG] UpdateUser - Executing query: %s\n", query)
	fmt.Printf("[DEBUG] UpdateUser - Query params: name=%v, birthdayInUnix=%v, gender=%v, dating_preference=%v\n",
		user.Name, user.BirthdayInUnix, user.Gender, user.DatingPreference)

	result, err := db.Exec(query,
		user.Name, user.UniversityName, workJSON, user.HomeTown,
		user.Height, user.Zodiac, user.CurrentlyPlaying, user.BirthdayInUnix,
		user.Gender, user.DatingPreference, lastPlayedSongJSON, user.UserLastActiveAt, user.ID,
	)

	if err != nil {
		fmt.Printf("[ERROR] UpdateUser - Error executing query: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("[DEBUG] UpdateUser - Rows affected: %d\n", rowsAffected)

	return nil
}

// UpdateSpotifyTokens updates a user's Spotify access token, refresh token, and expiry
func (db *DB) UpdateSpotifyTokens(userID uuid.UUID, accessToken, refreshToken string, tokenExpiry time.Time) error {
	query := `UPDATE users SET access_token = $1, refresh_token = $2, token_expiry = $3, updated_at = NOW() WHERE id = $4`
	_, err := db.Exec(query, accessToken, refreshToken, tokenExpiry, userID)
	return err
}

// SaveImage saves a user's image
func (db *DB) SaveImage(userID uuid.UUID, imageData []byte) error {
	imageID := uuid.New()
	query := `INSERT INTO images (id, user_id, data) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, imageID, userID, imageData)
	return err
}

// ClearUserImages removes all images for a user
func (db *DB) ClearUserImages(userID uuid.UUID) error {
	query := `DELETE FROM images WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	return err
}

// GetUserImages retrieves all images for a user
func (db *DB) GetUserImages(userID uuid.UUID) ([][]byte, error) {
	query := `SELECT data FROM images WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images [][]byte
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		images = append(images, data)
	}

	return images, nil
}

// SaveInterest saves a user's interest
func (db *DB) SaveInterest(userID uuid.UUID, interestName string) error {
	interestID := uuid.New()
	query := `INSERT INTO interests (id, user_id, name) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, interestID, userID, interestName)
	return err
}

// ClearUserInterests removes all interests for a user
func (db *DB) ClearUserInterests(userID uuid.UUID) error {
	query := `DELETE FROM interests WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	return err
}

// GetUserInterests retrieves all interests for a user
func (db *DB) GetUserInterests(userID uuid.UUID) ([]string, error) {
	query := `SELECT name FROM interests WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		interests = append(interests, name)
	}

	return interests, nil
}

// SaveInterestRating saves a user's interest rating
func (db *DB) SaveInterestRating(userID uuid.UUID, interestName string, rating int) error {
	interestRatingID := uuid.New()
	query := `INSERT INTO interest_ratings (id, user_id, name, rating) 
			 VALUES ($1, $2, $3, $4) 
			 ON CONFLICT (user_id, name) DO UPDATE SET rating = $4`
	_, err := db.Exec(query, interestRatingID, userID, interestName, rating)
	return err
}

// ClearUserInterestRatings removes all interest ratings for a user
func (db *DB) ClearUserInterestRatings(userID uuid.UUID) error {
	query := `DELETE FROM interest_ratings WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	return err
}

// GetUserInterestRatings retrieves all interest ratings for a user
func (db *DB) GetUserInterestRatings(userID uuid.UUID) (map[string]int, error) {
	query := `SELECT name, rating FROM interest_ratings WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := make(map[string]int)
	for rows.Next() {
		var name string
		var rating int
		if err := rows.Scan(&name, &rating); err != nil {
			return nil, err
		}
		ratings[name] = rating
	}

	return ratings, nil
}

// SavePrompt saves a user's prompt
func (db *DB) SavePrompt(userID uuid.UUID, question, answer string) error {
	promptID := uuid.New()
	query := `INSERT INTO prompts (id, user_id, question, answer) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, promptID, userID, question, answer)
	return err
}

// ClearUserPrompts removes all prompts for a user
func (db *DB) ClearUserPrompts(userID uuid.UUID) error {
	query := `DELETE FROM prompts WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	return err
}

// GetUserPrompts retrieves all prompts for a user
func (db *DB) GetUserPrompts(userID uuid.UUID) ([]models.Prompt, error) {
	query := `SELECT id, question, answer FROM prompts WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []models.Prompt
	for rows.Next() {
		var prompt models.Prompt
		if err := rows.Scan(&prompt.ID, &prompt.Question, &prompt.Answer); err != nil {
			return nil, err
		}
		prompt.UserID = userID
		prompts = append(prompts, prompt)
	}

	return prompts, nil
}

// SaveArtist saves a user's top artist
func (db *DB) SaveArtist(userID uuid.UUID, name, uri string, imageURL *string) error {
	artistID := uuid.New()
	query := `INSERT INTO artists (id, user_id, name, uri, image_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, artistID, userID, name, uri, imageURL)
	return err
}

// GetUserArtists retrieves all top artists for a user
func (db *DB) GetUserArtists(userID uuid.UUID) ([]models.Artist, error) {
	query := `SELECT id, name, uri, image_url FROM artists WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.Artist
	for rows.Next() {
		var artist models.Artist
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Uri, &artist.ImageURL); err != nil {
			return nil, err
		}
		artist.UserID = userID
		artists = append(artists, artist)
	}

	return artists, nil
}

// SaveSong saves a user's top song
func (db *DB) SaveSong(userID uuid.UUID, name, artist, uri string, imageURL *string) error {
	songID := uuid.New()
	query := `INSERT INTO songs (id, user_id, name, artist, uri, image_url) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, songID, userID, name, artist, uri, imageURL)
	return err
}

// GetUserSongs retrieves all top songs for a user
func (db *DB) GetUserSongs(userID uuid.UUID) ([]models.Song, error) {
	query := `SELECT id, name, artist, uri, image_url FROM songs WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Name, &song.Artist, &song.Uri, &song.ImageURL); err != nil {
			return nil, err
		}
		song.UserID = userID
		songs = append(songs, song)
	}

	return songs, nil
}

// SavePlaylist saves a user's saved playlist
func (db *DB) SavePlaylist(userID uuid.UUID, name, uri string, imageURL *string) error {
	playlistID := uuid.New()
	query := `INSERT INTO playlists (id, user_id, name, uri, image_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, playlistID, userID, name, uri, imageURL)
	return err
}

// GetUserPlaylists retrieves all saved playlists for a user
func (db *DB) GetUserPlaylists(userID uuid.UUID) ([]models.Playlist, error) {
	query := `SELECT id, name, uri, image_url FROM playlists WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var playlist models.Playlist
		if err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.Uri, &playlist.ImageURL); err != nil {
			return nil, err
		}
		playlist.UserID = userID
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// GetFullUserProfile retrieves the complete user profile
func (db *DB) GetFullUserProfile(userID uuid.UUID) (*models.UserProfile, error) {
	fmt.Println("[DEBUG] GetFullUserProfile - Retrieving profile for user:", userID)

	// First, get the user from GetUserByID to ensure we have all fields
	user, err := db.GetUserByID(userID)
	if err != nil {
		fmt.Printf("[ERROR] GetFullUserProfile - Error calling GetUserByID: %v\n", err)
		return nil, fmt.Errorf("error fetching user from GetUserByID: %v", err)
	}

	// Create a profile based on the user
	userProfile := &models.UserProfile{
		ID:               user.ID,
		Name:             user.Name,
		UniversityName:   user.UniversityName,
		Work:             user.Work,
		HomeTown:         user.HomeTown,
		Height:           user.Height,
		Age:              user.Age,
		Zodiac:           user.Zodiac,
		CurrentlyPlaying: user.CurrentlyPlaying,
		LastPlayedSong:   user.LastPlayedSong,
		UserLastActiveAt: user.UserLastActiveAt,
		BirthdayInUnix:   user.BirthdayInUnix,
		Gender:           user.Gender,
		DatingPreference: user.DatingPreference,
	}

	fmt.Printf("[DEBUG] UserProfile created with: BirthdayInUnix=%v, Gender=%v, DatingPreference=%v\n",
		userProfile.BirthdayInUnix, userProfile.Gender, userProfile.DatingPreference)

	// Get images
	images, err := db.GetUserImages(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching images: %v\n", err)
		return nil, fmt.Errorf("error fetching images: %v", err)
	}
	userProfile.Images = images

	// Get interests
	interests, err := db.GetUserInterests(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching interests: %v\n", err)
		return nil, fmt.Errorf("error fetching interests: %v", err)
	}
	userProfile.Interests = interests

	// Get interest ratings
	interestRatings, err := db.GetUserInterestRatings(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching interest ratings: %v\n", err)
		return nil, fmt.Errorf("error fetching interest ratings: %v", err)
	}
	userProfile.InterestRating = interestRatings

	// Get prompts
	prompts, err := db.GetUserPrompts(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching prompts: %v\n", err)
		return nil, fmt.Errorf("error fetching prompts: %v", err)
	}
	userProfile.Prompts = prompts

	// Get top artists
	topArtists, err := db.GetUserArtists(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching top artists: %v\n", err)
		return nil, fmt.Errorf("error fetching top artists: %v", err)
	}
	userProfile.TopArtists = topArtists

	// Get top songs
	topSongs, err := db.GetUserSongs(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching top songs: %v\n", err)
		return nil, fmt.Errorf("error fetching top songs: %v", err)
	}
	userProfile.TopSongs = topSongs

	// Get saved playlists
	savedPlaylists, err := db.GetUserPlaylists(userID)
	if err != nil {
		fmt.Printf("[ERROR] Error fetching saved playlists: %v\n", err)
		return nil, fmt.Errorf("error fetching saved playlists: %v", err)
	}
	userProfile.SavedPlaylists = savedPlaylists

	fmt.Printf("[DEBUG] Final userProfile: BirthdayInUnix=%v, Gender=%v, DatingPreference=%v\n",
		userProfile.BirthdayInUnix, userProfile.Gender, userProfile.DatingPreference)

	return userProfile, nil
}
