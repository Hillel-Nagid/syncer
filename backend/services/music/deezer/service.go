package deezer

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
	"strconv"
	"strings"
	"time"

	"syncer.net/core/services"
)

// DeezerService implements the ServiceProvider interface for Deezer API
type DeezerService struct {
	*services.BaseService
	appID     string
	appSecret string
}

// DeezerTrack represents a Deezer track with comprehensive metadata
type DeezerTrack struct {
	ID              int64             `json:"id"`
	Title           string            `json:"title"`
	Artist          DeezerArtist      `json:"artist"`
	Album           DeezerAlbum       `json:"album"`
	Duration        int               `json:"duration"`
	Rank            int               `json:"rank"`
	ExplicitContent bool              `json:"explicit_content_lyrics"`
	PreviewURL      string            `json:"preview"`
	TimeAdd         int64             `json:"time_add,omitempty"`
	ExternalRef     map[string]string `json:"external_reference,omitempty"`
	Link            string            `json:"link"`
}

// DeezerArtist represents a Deezer artist
type DeezerArtist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

// DeezerAlbum represents a Deezer album
type DeezerAlbum struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}

// DeezerPlaylist represents a Deezer playlist
type DeezerPlaylist struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Public       bool   `json:"public"`
	NbTracks     int    `json:"nb_tracks"`
	Duration     int    `json:"duration"`
	CreationDate string `json:"creation_date"`
	Link         string `json:"link"`
	Picture      string `json:"picture"`
}

// DeezerUser represents a Deezer user
type DeezerUser struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// NewDeezerService creates a new Deezer service with proper configuration
func NewDeezerService() *DeezerService {
	baseService := services.NewBaseService(services.BaseServiceConfig{
		Name:        "deezer",
		DisplayName: "Deezer",
		Category:    services.CategoryMusic,
		Scopes: []string{
			"basic_access",
			"email",
			"offline_access",
			"manage_library",
			"manage_community",
		},
		RequestsPerSecond: 10, // Deezer is more permissive than Spotify
		BurstSize:         15,
		HTTPTimeout:       30 * time.Second,
	})

	return &DeezerService{
		BaseService: baseService,
		appID:       os.Getenv("DEEZER_APP_ID"),
		appSecret:   os.Getenv("DEEZER_APP_SECRET"),
	}
}

// GetAuthURL generates the OAuth authorization URL for Deezer
func (d *DeezerService) GetAuthURL(state, redirectURL string) (string, error) {
	if d.appID == "" {
		return "", fmt.Errorf("DEEZER_APP_ID environment variable not set")
	}

	baseURL := "https://connect.deezer.com/oauth/auth.php"
	params := url.Values{
		"app_id":       {d.appID},
		"redirect_uri": {redirectURL},
		"perms":        {strings.Join(d.RequiredScopes(), ",")},
		"state":        {state},
	}

	authURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	d.LogInfo("Generated auth URL for Deezer")

	return authURL, nil
}

// ExchangeCode exchanges authorization code for access tokens
func (d *DeezerService) ExchangeCode(code, redirectURL string) (*services.OAuthTokens, error) {
	if d.appID == "" || d.appSecret == "" {
		return nil, fmt.Errorf("Deezer credentials not configured")
	}

	tokenURL := "https://connect.deezer.com/oauth/access_token.php"
	params := url.Values{
		"app_id":       {d.appID},
		"secret":       {d.appSecret},
		"code":         {code},
		"redirect_uri": {redirectURL},
	}

	// Deezer uses GET request for token exchange
	fullURL := fmt.Sprintf("%s?%s", tokenURL, params.Encode())

	if err := d.WaitForRateLimit(context.Background()); err != nil {
		return nil, err
	}

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Deezer returns access_token=TOKEN&expires=SECONDS
	responseParams, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	accessToken := responseParams.Get("access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("no access token in response: %s", body)
	}

	expiresIn, _ := strconv.Atoi(responseParams.Get("expires"))
	if expiresIn == 0 {
		expiresIn = 3600 // Default 1 hour
	}

	d.LogInfo("Successfully exchanged code for tokens")

	return &services.OAuthTokens{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}

// RefreshTokens - Deezer doesn't support refresh tokens, return error
func (d *DeezerService) RefreshTokens(refreshToken string) (*services.OAuthTokens, error) {
	return nil, fmt.Errorf("Deezer does not support refresh tokens - re-authentication required")
}

// GetUserProfile retrieves user profile information
func (d *DeezerService) GetUserProfile(ctx context.Context, tokens *services.OAuthTokens) (*services.UserProfile, error) {
	valid, err := d.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	if err := d.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.deezer.com/user/me?access_token=%s", tokens.AccessToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user profile (status %d): %s", resp.StatusCode, body)
	}

	var deezerUser struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Lastname string `json:"lastname"`
		Picture  string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deezerUser); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	profile := &services.UserProfile{
		ExternalID:  strconv.FormatInt(deezerUser.ID, 10),
		Username:    deezerUser.Name,
		Email:       deezerUser.Email,
		DisplayName: fmt.Sprintf("%s %s", deezerUser.Name, deezerUser.Lastname),
		AvatarURL:   deezerUser.Picture,
	}

	d.LogInfo("Successfully retrieved user profile")
	return profile, nil
}

// SyncUserData fetches user data from Deezer for cross-service sync
func (d *DeezerService) SyncUserData(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) (*services.UserDataResult, error) {
	d.LogInfo("Starting Deezer sync for user")

	valid, err := d.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	var allItems []services.SyncItem
	var errors []services.SyncError

	// Fetch user's favorite tracks
	favorites, err := d.fetchFavoriteTracks(ctx, tokens, lastSync)
	if err != nil {
		d.LogError("Failed to fetch favorite tracks: %v", err)
		errors = append(errors, d.CreateSyncError("favorites", err.Error(), "", "fetching favorites"))
	} else {
		allItems = append(allItems, favorites...)
		d.LogInfo("Fetched %d favorite tracks", len(favorites))
	}

	// Fetch user's playlists
	playlists, err := d.fetchUserPlaylists(ctx, tokens, lastSync)
	if err != nil {
		d.LogError("Failed to fetch playlists: %v", err)
		errors = append(errors, d.CreateSyncError("playlists", err.Error(), "", "fetching playlists"))
	} else {
		allItems = append(allItems, playlists...)
		d.LogInfo("Fetched %d playlist tracks", len(playlists))
	}

	// Fetch listening history
	history, err := d.fetchListeningHistory(ctx, tokens, lastSync)
	if err != nil {
		d.LogError("Failed to fetch listening history: %v", err)
		errors = append(errors, d.CreateSyncError("history", err.Error(), "", "fetching history"))
	} else {
		allItems = append(allItems, history...)
		d.LogInfo("Fetched %d history tracks", len(history))
	}

	d.LogInfo("Deezer sync completed: %d total items", len(allItems))

	result := &services.UserDataResult{
		Success: len(errors) == 0,
		Items:   allItems, // Transient data - never persisted
		Errors:  errors,
		Metadata: map[string]any{
			"sync_time":    time.Now(),
			"service":      "deezer",
			"items_synced": len(allItems),
			"last_sync":    lastSync,
		},
	}
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
			Data:       item.Data, // DeezerTrack -> any
		}
		anyResult.Items = append(anyResult.Items, anyItem)
	}

	return anyResult, nil
}

// fetchFavoriteTracks retrieves user's favorite tracks from Deezer
func (d *DeezerService) fetchFavoriteTracks(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem
	index := 0
	limit := 100

	for {
		if err := d.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.deezer.com/user/me/tracks?access_token=%s&index=%d&limit=%d",
			tokens.AccessToken, index, limit)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := d.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Deezer API error (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Data []struct {
				TimeAdd int64 `json:"time_add"`
				DeezerTrack
			} `json:"data"`
			Total int     `json:"total"`
			Next  *string `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode favorites response: %w", err)
		}

		// Process items
		for _, item := range result.Data {
			addedTime := time.Unix(item.TimeAdd, 0)
			if addedTime.After(lastSync) {
				item.DeezerTrack.TimeAdd = item.TimeAdd

				syncItem := services.SyncItem{
					ExternalID: strconv.FormatInt(item.DeezerTrack.ID, 10),
					ItemType:   "favorite_track",
					Action:     services.ActionCreate,
					Data:       item.DeezerTrack,
				}
				items = append(items, syncItem)
			}
		}

		if result.Next == nil || len(result.Data) < limit {
			break
		}
		index += limit

		// Safety check
		if index > result.Total {
			break
		}
	}

	return items, nil
}

// fetchUserPlaylists retrieves tracks from user's playlists
func (d *DeezerService) fetchUserPlaylists(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem

	// First, get user's playlists
	playlists, err := d.getUserPlaylists(ctx, tokens)
	if err != nil {
		return nil, err
	}

	// Then, get tracks from each playlist
	for _, playlist := range playlists {
		tracks, err := d.getPlaylistTracks(ctx, tokens, playlist.ID, lastSync)
		if err != nil {
			d.LogWarn("Failed to fetch tracks from playlist %d: %v", playlist.ID, err)
			continue
		}
		items = append(items, tracks...)
	}

	return items, nil
}

// getUserPlaylists gets user's playlists
func (d *DeezerService) getUserPlaylists(ctx context.Context, tokens *services.OAuthTokens) ([]DeezerPlaylist, error) {
	var playlists []DeezerPlaylist
	index := 0
	limit := 50

	for {
		if err := d.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.deezer.com/user/me/playlists?access_token=%s&index=%d&limit=%d",
			tokens.AccessToken, index, limit)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := d.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to get playlists (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Data  []DeezerPlaylist `json:"data"`
			Total int              `json:"total"`
			Next  *string          `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode playlists response: %w", err)
		}

		playlists = append(playlists, result.Data...)

		if result.Next == nil || len(result.Data) < limit {
			break
		}
		index += limit

		if index > result.Total {
			break
		}
	}

	return playlists, nil
}

// getPlaylistTracks gets tracks from a specific playlist
func (d *DeezerService) getPlaylistTracks(ctx context.Context, tokens *services.OAuthTokens, playlistID int64, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem
	index := 0
	limit := 100

	for {
		if err := d.WaitForRateLimit(ctx); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("https://api.deezer.com/playlist/%d/tracks?access_token=%s&index=%d&limit=%d",
			playlistID, tokens.AccessToken, index, limit)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := d.DoRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to get playlist tracks (status %d): %s", resp.StatusCode, body)
		}

		var result struct {
			Data  []DeezerTrack `json:"data"`
			Total int           `json:"total"`
			Next  *string       `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode playlist tracks response: %w", err)
		}

		// Process items - Note: Deezer playlist tracks don't have individual timestamps
		// so we use the last sync time as approximation
		for _, track := range result.Data {
			syncItem := services.SyncItem{
				ExternalID: strconv.FormatInt(track.ID, 10),
				ItemType:   "playlist_track",
				Action:     services.ActionCreate,
				Data:       track,
			}
			items = append(items, syncItem)
		}

		if result.Next == nil || len(result.Data) < limit {
			break
		}
		index += limit

		if index > result.Total {
			break
		}
	}

	return items, nil
}

// fetchListeningHistory retrieves user's listening history (flow)
func (d *DeezerService) fetchListeningHistory(ctx context.Context, tokens *services.OAuthTokens, lastSync time.Time) ([]services.SyncItem, error) {
	var items []services.SyncItem

	if err := d.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	// Get user's flow (listening history/recommendations)
	url := fmt.Sprintf("https://api.deezer.com/user/me/flow?access_token=%s&limit=50", tokens.AccessToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get flow (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Data []DeezerTrack `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode flow response: %w", err)
	}

	// Process items
	for _, track := range result.Data {
		syncItem := services.SyncItem{
			ExternalID: strconv.FormatInt(track.ID, 10),
			ItemType:   "flow_track",
			Action:     services.ActionCreate,
			Data:       track,
		}
		items = append(items, syncItem)
	}

	return items, nil
}

// generateTrackChecksum creates a checksum for change detection
func (d *DeezerService) generateTrackChecksum(track DeezerTrack) string {
	data := fmt.Sprintf("%d|%s|%s|%d",
		track.ID, track.Title, track.Artist.Name, track.Duration)
	hash := sha256.Sum256([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// HealthCheck performs a health check on the Deezer API
func (d *DeezerService) HealthCheck() error {
	// Simple health check by calling a public endpoint
	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://api.deezer.com/genre", nil)
	if err != nil {
		return err
	}

	resp, err := d.DoRequest(context.Background(), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("Deezer API health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Cross-service sync methods for adding tracks to Deezer

// AddTrack searches for a universal track on Deezer and adds it to user's favorites
func (d *DeezerService) AddTrack(ctx context.Context, tokens *services.OAuthTokens, universalTrack interface{}) error {
	// This would be implemented with the UniversalTrack type from the sync engine
	// For now, we'll create a placeholder that can be expanded when the universal types are ready
	d.LogInfo("AddTrack method called - will be enhanced with universal track support")
	return fmt.Errorf("AddTrack not yet implemented - awaiting universal track type implementation")
}

// SearchTrack searches for a track on Deezer using universal track data
func (d *DeezerService) SearchTrack(ctx context.Context, tokens *services.OAuthTokens, title, artist, album string) (*DeezerTrack, error) {
	valid, err := d.ValidateTokens(tokens)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid tokens: %w", err)
	}

	if err := d.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	// Construct search query
	query := fmt.Sprintf("%s %s", title, artist)
	if album != "" {
		query += " " + album
	}

	url := fmt.Sprintf("https://api.deezer.com/search?q=%s&access_token=%s&limit=1",
		url.QueryEscape(query), tokens.AccessToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Data []DeezerTrack `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("track not found on Deezer")
	}

	return &result.Data[0], nil
}

// AddToFavorites adds a track to user's favorite tracks
func (d *DeezerService) AddToFavorites(ctx context.Context, tokens *services.OAuthTokens, trackID int64) error {
	valid, err := d.ValidateTokens(tokens)
	if err != nil || !valid {
		return fmt.Errorf("invalid tokens: %w", err)
	}

	if err := d.WaitForRateLimit(ctx); err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.deezer.com/user/me/tracks?access_token=%s&track_id=%d",
		tokens.AccessToken, trackID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := d.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add to favorites (status %d): %s", resp.StatusCode, body)
	}

	d.LogInfo("Successfully added track %d to user's favorites", trackID)
	return nil
}

// RegisterWithRegistry allows the service to register itself with a service registry
func (d *DeezerService) RegisterWithRegistry(registry interface{}) error {
	if reg, ok := registry.(interface {
		Register(services.ServiceProvider) error
	}); ok {
		return reg.Register(d)
	}
	return fmt.Errorf("registry does not support service registration")
}
