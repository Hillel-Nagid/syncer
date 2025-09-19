package spotify

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"syncer.net/core/services"
)

// SpotifyService implements the ServiceProvider interface for Spotify Web API
type SpotifyService struct {
	*services.BaseService
	clientID     string
	clientSecret string
}

// SpotifyTrack represents a Spotify track with comprehensive metadata
type SpotifyTrack struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Artists         []SpotifyArtist   `json:"artists"`
	Album           SpotifyAlbum      `json:"album"`
	Duration        int               `json:"duration_ms"`
	Popularity      int               `json:"popularity"`
	ExternalIDs     map[string]string `json:"external_ids"`
	PreviewURL      *string           `json:"preview_url"`
	ExplicitContent bool              `json:"explicit"`
	AddedAt         *time.Time        `json:"added_at,omitempty"`
	URI             string            `json:"uri"`
}

// SpotifyArtist represents a Spotify artist
type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// SpotifyAlbum represents a Spotify album
type SpotifyAlbum struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Artists     []SpotifyArtist `json:"artists"`
	ReleaseDate string          `json:"release_date"`
	Images      []SpotifyImage  `json:"images"`
	URI         string          `json:"uri"`
}

// SpotifyImage represents a Spotify image
type SpotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// SpotifyPlaylist represents a Spotify playlist
type SpotifyPlaylist struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Public      bool           `json:"public"`
	Owner       SpotifyUser    `json:"owner"`
	Tracks      SpotifyTracks  `json:"tracks"`
	URI         string         `json:"uri"`
	Images      []SpotifyImage `json:"images"`
}

// SpotifyUser represents a Spotify user
type SpotifyUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// SpotifyTracks represents tracks in a playlist
type SpotifyTracks struct {
	Total int `json:"total"`
}

// NewSpotifyService creates a new Spotify service with proper configuration
func NewSpotifyService() *SpotifyService {
	baseService := services.NewBaseService(services.BaseServiceConfig{
		Name:        "spotify",
		DisplayName: "Spotify",
		Category:    services.CategoryMusic,
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"user-library-read",
			"user-library-modify",
			"playlist-read-private",
			"playlist-modify-public",
			"playlist-modify-private",
			"user-read-recently-played",
			"user-top-read",
		},
		RequestsPerSecond: 5,
		BurstSize:         10,
		HTTPTimeout:       30 * time.Second,
	})

	return &SpotifyService{
		BaseService:  baseService,
		clientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}

// GetAuthURL generates the OAuth authorization URL
func (s *SpotifyService) GetAuthURL(state, redirectURL string) (string, error) {
	if s.clientID == "" {
		return "", fmt.Errorf("SPOTIFY_CLIENT_ID environment variable not set")
	}

	baseURL := "https://accounts.spotify.com/authorize"
	params := url.Values{
		"client_id":     {s.clientID},
		"response_type": {"code"},
		"redirect_uri":  {redirectURL},
		"scope":         {strings.Join(s.RequiredScopes(), " ")},
		"state":         {state},
		"show_dialog":   {"true"},
	}

	authURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	s.LogInfo("Generated auth URL for Spotify")

	return authURL, nil
}

// ExchangeCode exchanges authorization code for access tokens
func (s *SpotifyService) ExchangeCode(code, redirectURL string) (*services.OAuthTokens, error) {
	if s.clientID == "" || s.clientSecret == "" {
		return nil, fmt.Errorf("Spotify credentials not configured")
	}

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirectURL},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.DoRequest(context.Background(), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, body)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	s.LogInfo("Successfully exchanged code for tokens")

	return &services.OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scope:        tokenResp.Scope,
	}, nil
}

// RefreshTokens refreshes expired access tokens
func (s *SpotifyService) RefreshTokens(refreshToken string) (*services.OAuthTokens, error) {
	if s.clientID == "" || s.clientSecret == "" {
		return nil, fmt.Errorf("Spotify credentials not configured")
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.DoRequest(context.Background(), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, body)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	if tokenResp.RefreshToken == "" {
		tokenResp.RefreshToken = refreshToken
	}

	s.LogInfo("Successfully refreshed tokens")

	return &services.OAuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scope:        tokenResp.Scope,
	}, nil
}

// GetUserProfile retrieves user profile information
func (s *SpotifyService) GetUserProfile(ctx context.Context, tokens *services.OAuthTokens) (*services.UserProfile, error) {
	valid, err := s.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	req, err := s.CreateAuthenticatedRequest(ctx, "GET", "https://api.spotify.com/v1/me", tokens)
	if err != nil {
		return nil, err
	}

	resp, err := s.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user profile (status %d): %s", resp.StatusCode, body)
	}

	var spotifyUser struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Images      []struct {
			URL string `json:"url"`
		} `json:"images"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&spotifyUser); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	profile := &services.UserProfile{
		ExternalID:  spotifyUser.ID,
		Username:    spotifyUser.ID,
		Email:       spotifyUser.Email,
		DisplayName: spotifyUser.DisplayName,
	}

	if len(spotifyUser.Images) > 0 {
		profile.AvatarURL = spotifyUser.Images[0].URL
	}

	s.LogInfo("Successfully retrieved user profile")
	return profile, nil
}

// SyncUserData fetches user data for real-time cross-service sync
func (s *SpotifyService) SyncUserData(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) (*services.UserDataResult, error) {
	s.LogInfo("Starting Spotify sync for user")

	valid, err := s.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	var allItems []services.SyncItem
	var errors []services.SyncError

	tracks, err := s.fetchSavedTracks(ctx, tokens, lastSync)
	if err != nil {
		s.LogError("Failed to fetch saved tracks: %v", err)
		errors = append(errors, s.CreateSyncError("saved_tracks", err.Error(), "", "fetching saved tracks"))
	} else {
		allItems = append(allItems, tracks...)
		s.LogInfo("Fetched %d saved tracks", len(tracks))
	}

	playlists, err := s.fetchUserPlaylists(ctx, tokens, lastSync)
	if err != nil {
		s.LogError("Failed to fetch playlists: %v", err)
		errors = append(errors, s.CreateSyncError("playlists", err.Error(), "", "fetching playlists"))
	} else {
		allItems = append(allItems, playlists...)
		s.LogInfo("Fetched %d playlist tracks", len(playlists))
	}

	recent, err := s.fetchRecentlyPlayed(ctx, tokens, lastSync)
	if err != nil {
		s.LogError("Failed to fetch recently played: %v", err)
		errors = append(errors, s.CreateSyncError("recently_played", err.Error(), "", "fetching recently played"))
	} else {
		allItems = append(allItems, recent...)
		s.LogInfo("Fetched %d recently played tracks", len(recent))
	}

	s.LogInfo("Spotify sync completed: %d total items", len(allItems))
	result := &services.UserDataResult{
		Success: len(errors) == 0,
		Items:   allItems, // Transient data - never persisted
		Errors:  errors,
		Metadata: map[string]any{
			"sync_time":    time.Now(),
			"service":      "spotify",
			"items_synced": len(allItems),
			"last_sync":    lastSync,
		},
	}
	if err != nil {
		return nil, err
	}

	// Convert SyncResult[SpotifyTrack] to SyncResult[any]
	anyResult := &services.UserDataResult{
		Success:  result.Success,
		Errors:   result.Errors,
		Metadata: result.Metadata,
	}

	// Convert items
	for _, item := range result.Items {
		anyItem := services.SyncItem{
			ExternalID: item.ExternalID,
			ItemType:   item.ItemType,
			Action:     item.Action,
			Data:       item.Data, // SpotifyTrack -> any
		}
		anyResult.Items = append(anyResult.Items, anyItem)
	}
	return anyResult, nil
}

// fetchSavedTracks retrieves user's liked songs
func (s *SpotifyService) fetchSavedTracks(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem
	offset := 0
	limit := 50

	for {
		if err := s.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?offset=%d&limit=%d", offset, limit)
		req, err := s.CreateAuthenticatedRequest(ctx, "GET", url, tokens)
		if err != nil {
			return nil, err
		}

		resp, err := s.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Spotify API error (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Items []struct {
				AddedAt time.Time    `json:"added_at"`
				Track   SpotifyTrack `json:"track"`
			} `json:"items"`
			Total int     `json:"total"`
			Next  *string `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode saved tracks response: %w", err)
		}

		for _, item := range result.Items {
			if item.AddedAt.After(lastSync) {
				item.Track.AddedAt = &item.AddedAt

				syncItem := services.SyncItem{
					ExternalID: item.Track.ID,
					ItemType:   "saved_track",
					Action:     services.ActionCreate,
					Data:       item.Track,
				}
				items = append(items, syncItem)
			}
		}

		if result.Next == nil {
			break
		}
		offset += limit

		if offset > result.Total {
			break
		}
	}

	return items, nil
}

// fetchUserPlaylists retrieves tracks from user's playlists
func (s *SpotifyService) fetchUserPlaylists(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem

	playlists, err := s.getUserPlaylists(ctx, tokens)
	if err != nil {
		return nil, err
	}

	for _, playlist := range playlists {
		tracks, err := s.getPlaylistTracks(ctx, tokens, playlist.ID, lastSync)
		if err != nil {
			s.LogWarn("Failed to fetch tracks from playlist %s: %v", playlist.ID, err)
			continue
		}
		items = append(items, tracks...)
	}

	return items, nil
}

// getUserPlaylists gets user's playlists
func (s *SpotifyService) getUserPlaylists(ctx context.Context, tokens *services.OAuthTokens) ([]SpotifyPlaylist, error) {
	var playlists []SpotifyPlaylist
	offset := 0
	limit := 50

	for {
		if err := s.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?offset=%d&limit=%d", offset, limit)
		req, err := s.CreateAuthenticatedRequest(ctx, "GET", url, tokens)
		if err != nil {
			return nil, err
		}

		resp, err := s.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to get playlists (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Items []SpotifyPlaylist `json:"items"`
			Total int               `json:"total"`
			Next  *string           `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode playlists response: %w", err)
		}

		playlists = append(playlists, result.Items...)

		if result.Next == nil {
			break
		}
		offset += limit

		if offset > result.Total {
			break
		}
	}

	return playlists, nil
}

// getPlaylistTracks gets tracks from a specific playlist
func (s *SpotifyService) getPlaylistTracks(ctx context.Context, tokens *services.OAuthTokens, playlistID string, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem
	offset := 0
	limit := 100

	for {
		if err := s.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=%d&limit=%d", playlistID, offset, limit)
		req, err := s.CreateAuthenticatedRequest(ctx, "GET", url, tokens)
		if err != nil {
			return nil, err
		}

		resp, err := s.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to get playlist tracks (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Items []struct {
				AddedAt time.Time    `json:"added_at"`
				Track   SpotifyTrack `json:"track"`
			} `json:"items"`
			Total int     `json:"total"`
			Next  *string `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode playlist tracks response: %w", err)
		}

		for _, item := range result.Items {
			if item.AddedAt.After(lastSync) {
				item.Track.AddedAt = &item.AddedAt

				syncItem := services.SyncItem{
					ExternalID: item.Track.ID,
					ItemType:   "playlist_track",
					Action:     services.ActionCreate,
					Data:       item.Track,
				}
				items = append(items, syncItem)
			}
		}

		if result.Next == nil {
			break
		}
		offset += limit

		if offset > result.Total {
			break
		}
	}

	return items, nil
}

// fetchRecentlyPlayed retrieves user's recently played tracks
func (s *SpotifyService) fetchRecentlyPlayed(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem

	if err := s.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/recently-played?limit=50&after=%d",
		lastSync.UnixNano()/1000000)

	req, err := s.CreateAuthenticatedRequest(ctx, "GET", url, tokens)
	if err != nil {
		return nil, err
	}

	resp, err := s.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get recently played (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Items []struct {
			PlayedAt time.Time    `json:"played_at"`
			Track    SpotifyTrack `json:"track"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode recently played response: %w", err)
	}

	for _, item := range result.Items {
		if item.PlayedAt.After(lastSync) {
			item.Track.AddedAt = &item.PlayedAt

			syncItem := services.SyncItem{
				ExternalID: item.Track.ID,
				ItemType:   "recently_played",
				Action:     services.ActionCreate,
				Data:       item.Track,
			}
			items = append(items, syncItem)
		}
	}

	return items, nil
}

// generateTrackChecksum creates a checksum for change detection
func (s *SpotifyService) generateTrackChecksum(track SpotifyTrack) string {
	data := fmt.Sprintf("%s|%s|%s|%d",
		track.ID, track.Name, track.Artists[0].Name, track.Duration)
	hash := sha256.Sum256([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// HealthCheck performs a health check on the Spotify API
func (s *SpotifyService) HealthCheck() error {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://api.spotify.com/v1/browse/featured-playlists?limit=1", nil)
	if err != nil {
		return err
	}

	resp, err := s.DoRequest(context.Background(), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("Spotify API health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Cross-service sync methods for adding tracks to Spotify

// AddTrack searches for a universal track on Spotify and adds it to user's library
func (s *SpotifyService) AddTrack(ctx context.Context, tokens *services.OAuthTokens, universalTrack interface{}) error {
	// This would be implemented with the UniversalTrack type from the sync engine
	// For now, we'll create a placeholder that can be expanded when the universal types are ready
	s.LogInfo("AddTrack method called - will be enhanced with universal track support")
	return fmt.Errorf("AddTrack not yet implemented - awaiting universal track type implementation")
}

// SearchTrack searches for a track on Spotify using universal track data
func (s *SpotifyService) SearchTrack(ctx context.Context, tokens *services.OAuthTokens, title, artist, album string) (*SpotifyTrack, error) {
	valid, err := s.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	if err := s.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("track:\"%s\" artist:\"%s\"", title, artist)
	if album != "" {
		query += fmt.Sprintf(" album:\"%s\"", album)
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1",
		url.QueryEscape(query))

	req, err := s.CreateAuthenticatedRequest(ctx, "GET", url, tokens)
	if err != nil {
		return nil, err
	}

	resp, err := s.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Tracks struct {
			Items []SpotifyTrack `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(result.Tracks.Items) == 0 {
		return nil, fmt.Errorf("track not found on Spotify")
	}

	return &result.Tracks.Items[0], nil
}

// SaveTrack adds a track to user's saved tracks
func (s *SpotifyService) SaveTrack(ctx context.Context, tokens *services.OAuthTokens, trackID string) error {
	valid, err := s.ValidateTokens(tokens)
	if err != nil || !valid {
		return fmt.Errorf("invalid tokens: %w", err)
	}

	if err := s.WaitForRateLimit(ctx); err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?ids=%s", trackID)
	req, err := s.CreateAuthenticatedRequest(ctx, "PUT", url, tokens)
	if err != nil {
		return err
	}

	resp, err := s.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to save track (status %d): %s", resp.StatusCode, body)
	}

	s.LogInfo("Successfully saved track %s to user's library", trackID)
	return nil
}

// RegisterWithRegistry allows the service to register itself with a service registry
func (s *SpotifyService) RegisterWithRegistry(registry interface{}) error {
	if reg, ok := registry.(interface {
		Register(services.ServiceProvider) error
	}); ok {
		return reg.Register(s)
	}
	return fmt.Errorf("registry does not support service registration")
}
