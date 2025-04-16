package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matchmyvibe/backend/internal/auth"
	"github.com/matchmyvibe/backend/internal/db"
	"github.com/matchmyvibe/backend/internal/models"
	"github.com/matchmyvibe/backend/internal/spotify"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	DB            *db.DB
	JWTService    *auth.JWTService
	SpotifyClient *spotify.Client
}

// SpotifyAuthRequest represents a request for authenticating with Spotify
type SpotifyAuthRequest struct {
	SpotifyURI   string    `json:"spotify_uri" binding:"required"`
	AccessToken  string    `json:"access_token" binding:"required"`
	RefreshToken string    `json:"refresh_token" binding:"required"`
	ExpiryDate   time.Time `json:"expiry_date" binding:"required"`
}

// AuthResponse represents the response from an authentication request
type AuthResponse struct {
	Token     string              `json:"token"`
	User      *models.UserProfile `json:"user"`
	IsNewUser bool                `json:"is_new_user"`
}

// SpotifyAuth handles authentication with Spotify
func (h *AuthHandler) SpotifyAuth(c *gin.Context) {
	var req SpotifyAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	user, err := h.DB.GetUserBySpotifyURI(req.SpotifyURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking for existing user"})
		return
	}

	isNewUser := false
	// If user doesn't exist, create a new one
	if user == nil {
		isNewUser = true
		user, err = h.DB.CreateUser(req.SpotifyURI, req.AccessToken, req.RefreshToken, req.ExpiryDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating new user"})
			return
		}
	} else {
		// Update the user's Spotify tokens
		err = h.DB.UpdateSpotifyTokens(user.ID, req.AccessToken, req.RefreshToken, req.ExpiryDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user tokens"})
			return
		}
	}

	// Generate a JWT token
	token, err := h.JWTService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	// Get user profile if not a new user
	var userProfile *models.UserProfile
	if !isNewUser {
		userProfile, err = h.DB.GetFullUserProfile(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user profile"})
			return
		}
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:     token,
		User:      userProfile,
		IsNewUser: isNewUser,
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Validate the token
	userID, err := h.JWTService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Generate a new token
	newToken, err := h.JWTService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
