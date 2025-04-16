# Postman Testing Guide for MatchMyVibe API

This guide explains how to use the provided Postman collection to test the MatchMyVibe API endpoints.

## Setup

1. Install [Postman](https://www.postman.com/downloads/) if you haven't already.
2. Import the collection:
   - Open Postman
   - Click "Import" in the top left
   - Select the `MatchMyVibe.postman_collection.json` file
   - The collection will appear in the left sidebar

3. Import the environment:
   - Click "Import" again
   - Select the `MatchMyVibe.postman_environment.json` file
   - Select the "MatchMyVibe Local" environment from the dropdown in the top right

## Authentication Flow

To test the API endpoints, you'll need to follow this workflow:

1. First, authenticate using Spotify credentials:
   - Open the "Spotify Auth" request in the "Authentication" folder
   - Replace the example values in the request body with your actual Spotify credentials
   - Send the request
   - From the response, copy the `token` value

2. Set up the token for subsequent requests:
   - Click the "eye" icon next to the environment dropdown
   - Edit the "MatchMyVibe Local" environment
   - Paste the token you copied into the `auth_token` variable
   - Click "Update"

3. Now you can use all other endpoints that require authentication

## Testing Endpoints

### Authentication

- **Spotify Auth**: Use this to authenticate with Spotify credentials and get a JWT token.
- **Refresh Token**: Use this to get a new JWT token using an existing valid token.

### User Profile

- **Get Profile**: Retrieves the current user's profile information.
- **Update Profile**: Updates specific fields in the user's profile. You can edit the request body to include only the fields you want to update.
- **Update Currently Playing**: Updates the currently playing track information from Spotify.

## Using Variables

The collection uses two key variables:

- `{{base_url}}`: The base URL of the API (default: http://localhost:8080)
- `{{auth_token}}`: The JWT token received after authentication

You can update these in the environment settings as needed.

## Testing Tips

1. After successful authentication, the token is valid for 24 hours.
2. You can test token refresh by using the "Refresh Token" endpoint before the token expires.
3. For the "Update Profile" endpoint, you can modify the request body to include only the fields you want to update.
4. The "Update Currently Playing" endpoint doesn't require a request body as it fetches the data directly from Spotify. 