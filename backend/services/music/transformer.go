package music

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"syncer.net/core/sync"
)

// MusicTrackTransformer handles conversion between different music service formats
// Implements the DataTransformer interface from core/sync
type MusicTrackTransformer struct {
	logger *log.Logger
}

// NewMusicTrackTransformer creates a new music track transformer
func NewMusicTrackTransformer() *MusicTrackTransformer {
	return &MusicTrackTransformer{
		logger: log.New(log.Writer(), "[MusicTransformer] ", log.LstdFlags),
	}
}

// TransformToUniversal converts service-specific data to universal format
func (t *MusicTrackTransformer) TransformToUniversal(serviceName string, data any) (sync.UniversalItem, error) {
	switch serviceName {
	case "spotify":
		return t.SpotifyToUniversal(data), nil
	case "deezer":
		return t.DeezerToUniversal(data), nil
	default:
		return UniversalTrack{}, fmt.Errorf("unsupported music service: %s", serviceName)
	}
}

// FindBestMatch finds the best matching track using multiple strategies
func (t *MusicTrackTransformer) FindBestMatch(sourceItem sync.UniversalItem, candidates []sync.UniversalItem, threshold float64) sync.UniversalMatch {
	sourceTrack, ok := sourceItem.(UniversalTrack)
	if !ok {
		t.logger.Printf("Source item is not a UniversalTrack")
		return sync.UniversalMatch{}
	}

	var universalCandidates []UniversalTrack
	for _, candidate := range candidates {
		if track, ok := candidate.(UniversalTrack); ok {
			universalCandidates = append(universalCandidates, track)
		}
	}

	bestMatch, bestScore := t.findBestMatchInternal(sourceTrack, universalCandidates, threshold)
	if bestMatch != nil {
		return sync.UniversalMatch{
			Target:     *bestMatch,
			Confidence: bestScore,
		}
	}
	return sync.UniversalMatch{}
}

// AnalyzeMatches provides statistics about match quality
func (t *MusicTrackTransformer) AnalyzeMatches(matches []sync.UniversalMatch) map[string]any {
	if len(matches) == 0 {
		return map[string]any{
			"total":              0,
			"average_confidence": 0.0,
			"perfect_matches":    0,
			"good_matches":       0,
			"poor_matches":       0,
		}
	}

	var totalConfidence float64
	perfectMatches := 0
	goodMatches := 0
	poorMatches := 0

	for _, match := range matches {
		totalConfidence += match.Confidence
		switch {
		case match.Confidence >= 0.95:
			perfectMatches++
		case match.Confidence >= 0.8:
			goodMatches++
		default:
			poorMatches++
		}
	}

	return map[string]any{
		"total":              len(matches),
		"average_confidence": totalConfidence / float64(len(matches)),
		"perfect_matches":    perfectMatches,
		"good_matches":       goodMatches,
		"poor_matches":       poorMatches,
	}
}

// SpotifyToUniversal converts Spotify track data to universal format
func (t *MusicTrackTransformer) SpotifyToUniversal(trackData any) UniversalTrack {
	// Use reflection or map access for generic data conversion
	if trackMap, ok := trackData.(map[string]any); ok {
		return t.mapToUniversalTrack(trackMap, "spotify")
	}

	t.logger.Printf("Warning: Unexpected Spotify track data type: %T", trackData)
	return UniversalTrack{}
}

// DeezerToUniversal converts Deezer track data to universal format
func (t *MusicTrackTransformer) DeezerToUniversal(trackData any) UniversalTrack {
	// Use reflection or map access for generic data conversion
	if trackMap, ok := trackData.(map[string]any); ok {
		return t.mapToUniversalTrack(trackMap, "deezer")
	}

	t.logger.Printf("Warning: Unexpected Deezer track data type: %T", trackData)
	return UniversalTrack{}
}

// mapToUniversalTrack converts a generic track map to UniversalTrack based on service type
func (t *MusicTrackTransformer) mapToUniversalTrack(trackMap map[string]any, serviceName string) UniversalTrack {
	switch serviceName {
	case "spotify":
		return t.spotifyMapToUniversal(trackMap)
	case "deezer":
		return t.deezerMapToUniversal(trackMap)
	default:
		// Generic conversion - extract common fields
		universal := UniversalTrack{
			ExternalIDs: make(map[string]string),
			Metadata:    make(map[string]any),
		}
		universal.Title = t.getStringField(trackMap, "name", "title")
		universal.Artist = t.getStringField(trackMap, "artist")
		universal.Album = t.getStringField(trackMap, "album")

		if duration, ok := trackMap["duration"]; ok {
			if durInt, ok := duration.(int); ok {
				universal.Duration = durInt
			}
		}

		universal.AddedAt = time.Now()
		return universal
	}
}

// spotifyMapToUniversal converts Spotify-specific map data
func (t *MusicTrackTransformer) spotifyMapToUniversal(trackMap map[string]any) UniversalTrack {
	universal := UniversalTrack{
		Title:       t.getStringField(trackMap, "name"),
		Duration:    t.getIntField(trackMap, "duration_ms"),
		ExternalIDs: map[string]string{"spotify": t.getStringField(trackMap, "id")},
		Metadata:    map[string]any{"original_service": "spotify"},
		AddedAt:     time.Now(),
	}

	// Extract artist information
	if artists, ok := trackMap["artists"].([]any); ok && len(artists) > 0 {
		var artistNames []string
		for _, artist := range artists {
			if artistMap, ok := artist.(map[string]any); ok {
				if name := t.getStringField(artistMap, "name"); name != "" {
					artistNames = append(artistNames, name)
				}
			}
		}
		universal.Artist = strings.Join(artistNames, ", ")
	}

	// Extract album information
	if album, ok := trackMap["album"].(map[string]any); ok {
		universal.Album = t.getStringField(album, "name")
	}

	// Extract ISRC
	if externalIDs, ok := trackMap["external_ids"].(map[string]any); ok {
		if isrc := t.getStringField(externalIDs, "isrc"); isrc != "" {
			universal.ISRC = isrc
		}
	}

	return universal
}

// deezerMapToUniversal converts Deezer-specific map data
func (t *MusicTrackTransformer) deezerMapToUniversal(trackMap map[string]any) UniversalTrack {
	universal := UniversalTrack{
		Title:       t.getStringField(trackMap, "title"),
		Duration:    t.getIntField(trackMap, "duration") * 1000, // Convert seconds to ms
		ExternalIDs: map[string]string{"deezer": strconv.FormatInt(int64(t.getIntField(trackMap, "id")), 10)},
		Metadata:    map[string]any{"original_service": "deezer"},
		AddedAt:     time.Now(),
	}

	// Extract artist information
	if artist, ok := trackMap["artist"].(map[string]any); ok {
		universal.Artist = t.getStringField(artist, "name")
	}

	// Extract album information
	if album, ok := trackMap["album"].(map[string]any); ok {
		universal.Album = t.getStringField(album, "title")
	}

	// Handle time_add for added timestamp
	if timeAdd := t.getIntField(trackMap, "time_add"); timeAdd > 0 {
		universal.AddedAt = time.Unix(int64(timeAdd), 0)
	}

	return universal
}

// Helper methods for safe type conversion
func (t *MusicTrackTransformer) getStringField(data map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, exists := data[key]; exists {
			if str, ok := value.(string); ok {
				return str
			}
		}
	}
	return ""
}

func (t *MusicTrackTransformer) getIntField(data map[string]any, key string) int {
	if value, exists := data[key]; exists {
		switch v := value.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

// findBestMatchInternal finds the best matching track using multiple strategies
func (t *MusicTrackTransformer) findBestMatchInternal(sourceTrack UniversalTrack, candidates []UniversalTrack, threshold float64) (*UniversalTrack, float64) {
	var bestMatch *UniversalTrack
	var bestScore float64

	for i := range candidates {
		score := t.calculateMatchScore(sourceTrack, candidates[i])
		if score > bestScore && score >= threshold {
			bestScore = score
			bestMatch = &candidates[i]
		}
	}

	return bestMatch, bestScore
}

// calculateMatchScore uses multiple factors to determine track similarity
func (t *MusicTrackTransformer) calculateMatchScore(track1, track2 UniversalTrack) float64 {
	// ISRC matching - highest priority
	if track1.ISRC != "" && track2.ISRC != "" {
		if track1.ISRC == track2.ISRC {
			t.logger.Printf("Perfect ISRC match: %s", track1.ISRC)
			return 1.0 // Perfect match
		}
		t.logger.Printf("ISRC mismatch: %s != %s", track1.ISRC, track2.ISRC)
		return 0.0 // ISRC mismatch means different tracks
	}

	// Calculate similarity scores for different fields
	titleScore := t.calculateStringSimilarity(
		strings.ToLower(track1.Title),
		strings.ToLower(track2.Title),
	)

	artistScore := t.calculateStringSimilarity(
		strings.ToLower(track1.Artist),
		strings.ToLower(track2.Artist),
	)

	albumScore := t.calculateStringSimilarity(
		strings.ToLower(track1.Album),
		strings.ToLower(track2.Album),
	)

	// Duration similarity (allow 5% variance)
	durationScore := 0.0
	if track1.Duration > 0 && track2.Duration > 0 {
		diff := float64(abs(track1.Duration-track2.Duration)) / float64(max(track1.Duration, track2.Duration))
		if diff <= 0.05 {
			durationScore = 1.0 - diff
		}
	}

	// Weighted scoring: Title and Artist are most important
	score := (titleScore * 0.4) + (artistScore * 0.4) + (albumScore * 0.15) + (durationScore * 0.05)

	t.logger.Printf("Match score: '%s' - '%s' = %.3f (title:%.3f, artist:%.3f, album:%.3f, duration:%.3f)",
		track1.Title, track2.Title, score, titleScore, artistScore, albumScore, durationScore)

	return score
}

// calculateStringSimilarity uses Levenshtein distance for string comparison
func (t *MusicTrackTransformer) calculateStringSimilarity(s1, s2 string) float64 {
	// Normalize strings
	s1 = t.normalizeString(s1)
	s2 = t.normalizeString(s2)

	if s1 == s2 {
		return 1.0
	}

	// Calculate Levenshtein distance
	distance := levenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))

	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - (float64(distance) / float64(maxLen))
}

// normalizeString removes common variations that affect matching
func (t *MusicTrackTransformer) normalizeString(s string) string {
	s = strings.ToLower(s)

	// Remove common parenthetical additions
	s = regexp.MustCompile(`\s*\([^)]*\)`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`\s*\[[^\]]*\]`).ReplaceAllString(s, "")

	// Remove featuring information
	s = regexp.MustCompile(`\s*(feat|ft|featuring)\.?\s+.*`).ReplaceAllString(s, "")

	// Remove common suffixes
	s = regexp.MustCompile(`\s*-\s*(remix|extended|radio edit|clean|explicit).*$`).ReplaceAllString(s, "")

	// Replace special characters with spaces
	s = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(s, " ")

	// Remove extra whitespace
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	return s
}

// GenerateSearchQuery creates an optimized search query for a universal track
func (t *MusicTrackTransformer) GenerateSearchQuery(track UniversalTrack) string {
	// Start with title and main artist
	query := fmt.Sprintf("\"%s\" \"%s\"", track.Title, track.Artist)

	// Add album if available and not too long
	if track.Album != "" && len(track.Album) < 50 {
		query += fmt.Sprintf(" \"%s\"", track.Album)
	}

	return query
}

// NormalizeTrackTitle normalizes a track title for better matching
func (t *MusicTrackTransformer) NormalizeTrackTitle(title string) string {
	return t.normalizeString(title)
}

// NormalizeArtistName normalizes an artist name for better matching
func (t *MusicTrackTransformer) NormalizeArtistName(artist string) string {
	// Handle featuring artists differently
	normalized := t.normalizeString(artist)

	// Split on "feat" variations and take the main artist
	parts := regexp.MustCompile(`\s+(feat|ft|featuring)\s+`).Split(normalized, -1)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}

	return normalized
}

// Utility functions

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}

	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// MatchesSyncType checks if an item type matches the requested sync type for music domain
func (t *MusicTrackTransformer) MatchesSyncType(itemType string, syncType string) bool {
	switch syncType {
	case string(MusicSyncTypeFavorites):
		return itemType == "saved_track" || itemType == "favorite_track"
	case string(MusicSyncTypePlaylists):
		return itemType == "playlist_track"
	case string(MusicSyncTypeRecentlyPlayed):
		return itemType == "recently_played" || itemType == "flow_track"
	default:
		t.logger.Printf("Unknown music sync type '%s', including all items", syncType)
		return true
	}
}
