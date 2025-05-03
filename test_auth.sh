#!/bin/bash

# Authenticate with test credentials
curl -X POST http://localhost:8080/auth/spotify \
  -H "Content-Type: application/json" \
  -d '{
    "spotify_uri": "spotify:user:testuser",
    "access_token": "test_access_token",
    "refresh_token": "test_refresh_token",
    "expiry_date": "2023-05-16T12:00:00Z"
  }' 