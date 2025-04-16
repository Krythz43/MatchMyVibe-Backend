# MatchMyVibe API Postman Collection

This folder contains Postman files for testing the MatchMyVibe backend API.

## Files

- `matchmyvibe_api_collection.json` - The Postman collection with all API endpoints
- `environment.json` - Environment variables for the collection

## How to Import

1. Open Postman
2. Click the "Import" button in the top left
3. Select "Folder" and navigate to this directory
4. Select all files and click "Open"

Alternatively, you can import individual files:

1. Open Postman
2. Click the "Import" button in the top left
3. Select "File" and navigate to the desired JSON file
4. Click "Open"

## Setting Up the Environment

1. After importing, go to the "Environments" tab in Postman
2. Select the imported "MatchMyVibe Local Environment"
3. Make sure the `base_url` is set correctly (default: http://localhost:8080)
4. Click "Save"

## Using the Collection

1. Make sure your Go backend server is running
2. Select the "MatchMyVibe Local Environment" from the environment dropdown
3. Start with the "Spotify Auth" request to authenticate (this will save the token automatically)
4. Once authenticated, you can use all other API endpoints

## Automatic Token Handling

The collection is set up to automatically extract and save the auth token from the "Spotify Auth" response. If you need to manually set a token:

1. Go to the "Environments" tab
2. Select the "MatchMyVibe Local Environment"
3. Enter your JWT token in the `auth_token` field
4. Click "Save"

## Testing Workflow

1. Run "Spotify Auth" to create a user and get a token
2. Run "Get Profile" to retrieve the user profile
3. Test various profile updates using the endpoints in the "User Profile" folder
4. Test partial updates using the endpoints in the "Advanced Profile Updates" folder

## Notes

- The `$isoTimestamp` variable in the Spotify Auth request will automatically set the current timestamp
- If your API is running on a different port, update the `base_url` in the environment settings 