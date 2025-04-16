// Tests for Profile endpoints

// Get Profile Tests
pm.test("Status code is 200 OK", function () {
    pm.response.to.have.status(200);
});

pm.test("Response is valid user profile", function () {
    var jsonData = pm.response.json();
    
    // Check basic profile properties
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData.id).to.be.a('string');
    
    // These may be null for new users
    pm.expect(jsonData).to.have.property('name');
    pm.expect(jsonData).to.have.property('university_name');
    pm.expect(jsonData).to.have.property('home_town');
    pm.expect(jsonData).to.have.property('height');
    pm.expect(jsonData).to.have.property('age');
    pm.expect(jsonData).to.have.property('zodiac');
    
    // Arrays should exist even if empty
    pm.expect(jsonData).to.have.property('interests').that.is.an('array');
    pm.expect(jsonData).to.have.property('top_artists').that.is.an('array');
    pm.expect(jsonData).to.have.property('top_songs').that.is.an('array');
    pm.expect(jsonData).to.have.property('saved_playlists').that.is.an('array');
    pm.expect(jsonData).to.have.property('prompts').that.is.an('array');
    
    // Interest ratings should be an object
    pm.expect(jsonData).to.have.property('interest_rating').that.is.an('object');
});

// Update Profile Tests
pm.test("Profile update returns 200 OK", function () {
    pm.response.to.have.status(200);
});

pm.test("Profile update returns updated values", function () {
    var jsonData = pm.response.json();
    var requestData = JSON.parse(pm.request.body.raw);
    
    // Check if the response includes the updated values
    Object.keys(requestData).forEach(function(key) {
        if (key === 'work') {
            if (requestData.work && jsonData.work) {
                Object.keys(requestData.work).forEach(function(workKey) {
                    pm.expect(jsonData.work[workKey]).to.eql(requestData.work[workKey]);
                });
            }
        } else if (key === 'interests') {
            if (requestData.interests) {
                requestData.interests.forEach(function(interest) {
                    pm.expect(jsonData.interests).to.include(interest);
                });
            }
        } else if (key === 'interest_rating') {
            if (requestData.interest_rating) {
                Object.keys(requestData.interest_rating).forEach(function(interestKey) {
                    pm.expect(jsonData.interest_rating[interestKey]).to.eql(requestData.interest_rating[interestKey]);
                });
            }
        } else if (key === 'prompts') {
            // Prompts are more complex to verify, we'll just check the count matches
            if (requestData.prompts) {
                pm.expect(jsonData.prompts.length).to.eql(requestData.prompts.length);
            }
        } else if (key !== 'images') { // Skip image verification
            if (requestData[key] !== null && requestData[key] !== undefined) {
                pm.expect(jsonData[key]).to.eql(requestData[key]);
            }
        }
    });
});

// Currently Playing Tests
pm.test("Currently playing update returns 200 OK", function () {
    pm.response.to.have.status(200);
});

pm.test("Currently playing update returns track info", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('currently_playing');
}); 