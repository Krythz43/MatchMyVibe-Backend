package handlers

import (
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
	Name           *string             `json:"name"`
	UniversityName *string             `json:"university_name"`
	Work           *models.WorkProfile `json:"work"`
	HomeTown       *string             `json:"home_town"`
	Height         *string             `json:"height"`
	Zodiac         *string             `json:"zodiac"`
	Images         [][]byte            `json:"images"`
	Interests      []string            `json:"interests"`
	InterestRating map[string]int      `json:"interest_rating"`
	Prompts        []struct {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current user to update
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
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

	// Update the user in the database
	if err := h.DB.UpdateUser(user); err != nil {
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

// UpdateCurrentlyPlaying updates the user's currently playing track
func (h *ProfileHandler) UpdateCurrentlyPlaying(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the user to check if the token is valid or needs refresh
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user"})
		return
	}

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
	if err := h.DB.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating currently playing track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"currently_playing": currentlyPlaying})
}
