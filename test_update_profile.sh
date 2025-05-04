#!/bin/bash

# Use the JWT token from authentication
AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzE1ZTJiOTgtMGE4Ni00OWI5LTg5ZjktZmJjNzA5YTE5MDEzIiwiaXNzIjoibWF0Y2hteXZpYmUtYmFja2VuZCIsInN1YiI6IjcxNWUyYjk4LTBhODYtNDliOS04OWY5LWZiYzcwOWExOTAxMyIsImV4cCI6MTc0NjM5NzQ5MywibmJmIjoxNzQ2MzExMDkzLCJpYXQiOjE3NDYzMTEwOTN9.NK3O-KzxeLVI0fXHw9p1jzlVszqD-WZiIiumPXcmAKs"

# Update a user profile with the new fields
curl -X PUT http://localhost:8080/api/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -d '{
    "name": "Test User",
    "birthdayInUnix": 631152000,
    "gender": "Man",
    "dating_preference": "Women"
  }'

echo -e "\n\nGet the updated profile:"
# Get the user profile to verify the changes
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $AUTH_TOKEN" 