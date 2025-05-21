package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matchmyvibe/backend/internal/db"
	"github.com/matchmyvibe/backend/internal/middleware"
	"github.com/matchmyvibe/backend/internal/models"
	"github.com/matchmyvibe/backend/internal/spotify"
)

// ProfileHandler handles profile-related requests
type ProfileHandler struct {
	DB            *db.DB
	SpotifyClient *spotify.Client
}

// GetProfile retrieves the user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.DB.GetFullUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfileRequest represents a request to update a user's profile
type UpdateProfileRequest struct {
	Name             *string             `json:"name"`
	UniversityName   *string             `json:"university_name"`
	Work             *models.WorkProfile `json:"work"`
	HomeTown         *string             `json:"home_town"`
	Height           *string             `json:"height"`
	Zodiac           *string             `json:"zodiac"`
	BirthdayInUnix   *int64              `json:"birthdayInUnix"`
	Gender           *string             `json:"gender"`
	DatingPreference *string             `json:"dating_preference"`
	Images           [][]byte            `json:"images"`
	Interests        []string            `json:"interests"`
	InterestRating   map[string]int      `json:"interest_rating"`
	Prompts          []struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	} `json:"prompts"`
}

// UpdateProfile updates the user's profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[ERROR] UpdateProfile - Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[DEBUG] UpdateProfile - Request: %+v\n", req)
	fmt.Printf("[DEBUG] UpdateProfile - BirthdayInUnix: %v, Gender: %v, DatingPreference: %v\n",
		req.BirthdayInUnix, req.Gender, req.DatingPreference)

	// Get the current user to update
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		fmt.Printf("[ERROR] UpdateProfile - Error retrieving user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user"})
		return
	}

	// Update the user's fields if provided
	if req.Name != nil {
		user.Name = req.Name
	}
	if req.UniversityName != nil {
		user.UniversityName = req.UniversityName
	}
	if req.Work != nil {
		user.Work = req.Work
	}
	if req.HomeTown != nil {
		user.HomeTown = req.HomeTown
	}
	if req.Height != nil {
		user.Height = req.Height
	}
	if req.Zodiac != nil {
		user.Zodiac = req.Zodiac
	}
	if req.BirthdayInUnix != nil {
		fmt.Printf("[DEBUG] UpdateProfile - Setting BirthdayInUnix to: %v\n", *req.BirthdayInUnix)
		user.BirthdayInUnix = req.BirthdayInUnix
	}
	if req.Gender != nil {
		// Validate gender
		if *req.Gender != "Man" && *req.Gender != "Woman" && *req.Gender != "Non-binary" {
			fmt.Printf("[ERROR] UpdateProfile - Invalid gender value: %v\n", *req.Gender)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gender value"})
			return
		}
		fmt.Printf("[DEBUG] UpdateProfile - Setting Gender to: %v\n", *req.Gender)
		user.Gender = req.Gender
	}
	if req.DatingPreference != nil {
		// Validate dating preference
		if *req.DatingPreference != "Men" && *req.DatingPreference != "Women" && *req.DatingPreference != "Everyone" {
			fmt.Printf("[ERROR] UpdateProfile - Invalid dating preference value: %v\n", *req.DatingPreference)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dating preference value"})
			return
		}
		fmt.Printf("[DEBUG] UpdateProfile - Setting DatingPreference to: %v\n", *req.DatingPreference)
		user.DatingPreference = req.DatingPreference
	}

	fmt.Printf("[DEBUG] UpdateProfile - Before DB update: BirthdayInUnix=%v, Gender=%v, DatingPreference=%v\n",
		user.BirthdayInUnix, user.Gender, user.DatingPreference)

	// Update the user in the database
	if err := h.DB.UpdateUser(user); err != nil {
		fmt.Printf("[ERROR] UpdateProfile - Error updating user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user"})
		return
	}

	// Update images if provided
	if req.Images != nil {
		// Clear existing images and add new ones
		if err := h.DB.ClearUserImages(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error clearing user images"})
			return
		}

		for _, imageData := range req.Images {
			if err := h.DB.SaveImage(userID, imageData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving image"})
				return
			}
		}
	}

	// Update interests if provided
	if req.Interests != nil {
		// Clear existing interests and add new ones
		if err := h.DB.ClearUserInterests(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error clearing user interests"})
			return
		}

		for _, interest := range req.Interests {
			if err := h.DB.SaveInterest(userID, interest); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving interest"})
				return
			}
		}
	}

	// Update interest ratings if provided
	if req.InterestRating != nil {
		// Clear existing interest ratings and add new ones
		if err := h.DB.ClearUserInterestRatings(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error clearing user interest ratings"})
			return
		}

		for interest, rating := range req.InterestRating {
			if err := h.DB.SaveInterestRating(userID, interest, rating); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving interest rating"})
				return
			}
		}
	}

	// Update prompts if provided
	if req.Prompts != nil {
		// Clear existing prompts and add new ones
		if err := h.DB.ClearUserPrompts(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error clearing user prompts"})
			return
		}

		for _, prompt := range req.Prompts {
			if err := h.DB.SavePrompt(userID, prompt.Question, prompt.Answer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving prompt"})
				return
			}
		}
	}

	// Get the updated profile
	profile, err := h.DB.GetFullUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving updated user profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateCurrentlyPlayingRequest represents the request body for updating the user's currently playing track
type UpdateCurrentlyPlayingRequest struct {
	Track        string `json:"track"`
	Artist       string `json:"artist"`
	URI          string `json:"uri"`
	Album        string `json:"album"`
	AlbumURI     string `json:"album_uri"`
	Duration     int    `json:"duration"`
	ContextTitle string `json:"context_title"`
	ContextURI   string `json:"context_uri"`
}

// UpdateCurrentlyPlaying updates the user's currently playing track
func (h *ProfileHandler) UpdateCurrentlyPlaying(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the user
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user"})
		return
	}

	// Parse the request body
	var req UpdateCurrentlyPlayingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no JSON is provided, fall back to the old behavior of fetching from Spotify
		handleSpotifyCurrentlyPlaying(c, h, user)
		return
	}

	// Create formatted currently playing string
	currentlyPlaying := fmt.Sprintf("%s - %s", req.Track, req.Artist)
	user.CurrentlyPlaying = &currentlyPlaying

	// Create last played song object
	lastPlayedSong := &models.LastPlayedSong{
		Track:        req.Track,
		Artist:       req.Artist,
		URI:          req.URI,
		Album:        req.Album,
		AlbumURI:     req.AlbumURI,
		Duration:     req.Duration,
		ContextTitle: req.ContextTitle,
		ContextURI:   req.ContextURI,
	}
	user.LastPlayedSong = lastPlayedSong

	// Update user's last active timestamp
	now := time.Now().Unix()
	user.UserLastActiveAt = &now

	// Update the user in the database
	if err := h.DB.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currently_playing":   currentlyPlaying,
		"last_played_song":    lastPlayedSong,
		"user_last_active_at": now,
	})
}

// handleSpotifyCurrentlyPlaying is the legacy function to fetch currently playing from Spotify
func handleSpotifyCurrentlyPlaying(c *gin.Context, h *ProfileHandler, user *models.User) {
	// Check if the access token needs to be refreshed
	if time.Now().After(user.TokenExpiry) {
		// Refresh the token
		tokenResponse, err := h.SpotifyClient.RefreshToken(user.RefreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error refreshing Spotify token"})
			return
		}

		// Calculate the new expiry time
		newExpiry := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

		// Update the user's tokens in the database
		err = h.DB.UpdateSpotifyTokens(user.ID, tokenResponse.AccessToken, tokenResponse.RefreshToken, newExpiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user tokens"})
			return
		}

		// Update the access token for the request
		user.AccessToken = tokenResponse.AccessToken
	}

	// Get the currently playing track from Spotify
	currentlyPlaying, err := h.SpotifyClient.GetCurrentlyPlaying(user.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching currently playing track"})
		return
	}

	// Update the user's currently playing track in the database
	user.CurrentlyPlaying = &currentlyPlaying

	// Update user's last active timestamp
	now := time.Now().Unix()
	user.UserLastActiveAt = &now

	if err := h.DB.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating currently playing track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currently_playing":   currentlyPlaying,
		"user_last_active_at": now,
	})
}
