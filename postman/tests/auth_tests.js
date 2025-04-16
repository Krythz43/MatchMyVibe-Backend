// Tests for Spotify Auth endpoint
pm.test("Status code is 200 OK", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has token property", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('token');
    pm.expect(jsonData.token).to.be.a('string');
    pm.expect(jsonData.token).to.not.be.empty;
});

pm.test("Response has is_new_user property", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('is_new_user');
    pm.expect(jsonData.is_new_user).to.be.a('boolean');
});

// Store the token in environment variable for use in other requests
pm.test("Save token to environment", function () {
    var jsonData = pm.response.json();
    if (jsonData.token) {
        pm.environment.set("auth_token", jsonData.token);
        console.log("Auth token saved to environment");
    }
});

// Check for user profile in the response
pm.test("Response has user property for existing users", function () {
    var jsonData = pm.response.json();
    if (jsonData.is_new_user === false) {
        pm.expect(jsonData).to.have.property('user');
        pm.expect(jsonData.user).to.be.an('object');
    }
}); 