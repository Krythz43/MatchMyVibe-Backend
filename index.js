/**
 * MatchMyVibe Backend Server
 * 
 * This Express server handles Supabase authentication with Spotify OAuth
 * and provides endpoints for retrieving user's top artists and tracks.
 */
const express = require('express');
const dotenv = require('dotenv');
const { createClient } = require('@supabase/supabase-js');
const cors = require('cors');
const path = require('path');

// Load environment variables
dotenv.config();

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, 'public')));

// Initialize Supabase client
const supabaseUrl = process.env.SUPABASE_URL;
const supabaseKey = process.env.SUPABASE_ANON_KEY;
const supabase = createClient(supabaseUrl, supabaseKey);

console.log(supabaseUrl, supabaseKey);
console.log(supabase);

// Routes
app.get('/', (req, res) => {
  console.log('Home page was requested');
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// Serve the callback page for Spotify OAuth
app.get('/callback', (req, res) => {
  console.log('Callback page was requested');
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// API endpoint to fetch current user's profile
app.get('/api/user', async (req, res) => {
  console.log('User endpoint was requested');
  const token = req.headers.authorization?.split(' ')[1];
  
  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }
  
  try {
    // Fetch user from Supabase
    const { data: { user }, error } = await supabase.auth.getUser(token);
    
    if (error) throw error;
    
    // Fetch profile from our database
    const { data: profile, error: profileError } = await supabase
      .from('profiles')
      .select('*')
      .eq('id', user.id)
      .single();
    
    if (profileError && profileError.code !== 'PGRST116') {
      throw profileError;
    }
    
    res.json({ user, profile });
  } catch (error) {
    console.error('Error fetching user:', error);
    res.status(500).json({ error: error.message });
  }
});

// API endpoint to fetch top artists and tracks
app.get('/api/spotify-data', async (req, res) => {
  console.log('Spotify data endpoint was requested');
  const token = req.headers.authorization?.split(' ')[1];
  
  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }
  
  try {
    // Get the current user
    const { data: { user }, error } = await supabase.auth.getUser(token);
    
    if (error) throw error;
    
    // Fetch the user's top artists
    const { data: topArtists, error: artistsError } = await supabase
      .from('top_artists')
      .select('*')
      .eq('profile_id', user.id)
      .eq('time_range', 'short_term') // Last 4 weeks
      .order('popularity', { ascending: false })
      .limit(5);
    
    if (artistsError) throw artistsError;
    
    // Fetch the user's top tracks
    const { data: topTracks, error: tracksError } = await supabase
      .from('top_tracks')
      .select('*')
      .eq('profile_id', user.id)
      .eq('time_range', 'short_term') // Last 4 weeks
      .order('popularity', { ascending: false })
      .limit(5);
    
    if (tracksError) throw tracksError;
    
    res.json({ topArtists, topTracks });
  } catch (error) {
    console.error('Error fetching Spotify data:', error);
    res.status(500).json({ error: error.message });
  }
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
}); 