package music

import (
	"time"

	"syncer.net/core/services"
)

// UniversalTrack represents a track in a platform-agnostic format for cross-service sync
type UniversalTrack struct {
	Title       string            `json:"title"`
	Artist      string            `json:"artist"`
	Album       string            `json:"album"`
	Duration    int               `json:"duration_ms"`
	ISRC        string            `json:"isrc,omitempty"`         // International Standard Recording Code
	ExternalIDs map[string]string `json:"external_ids,omitempty"` // Original service IDs
	Metadata    map[string]any    `json:"metadata,omitempty"`     // Additional platform-specific data
	AddedAt     time.Time         `json:"added_at"`
	Action      services.SyncAction
}

func (u UniversalTrack) GetItemAction() services.SyncAction {
	return u.Action
}
func (u UniversalTrack) GetItemIdentifier() string {
	return u.Title
}

func (u UniversalTrack) GetItemType() string {
	return "track"
}

// UniversalPlaylist represents a playlist in a platform-agnostic format
type UniversalPlaylist struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Public      bool              `json:"public"`
	Tracks      []UniversalTrack  `json:"tracks"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
}

// MusicSyncType defines what type of music data to sync
type MusicSyncType string

const (
	MusicSyncTypeFavorites      MusicSyncType = "favorites"       // Liked/saved tracks
	MusicSyncTypePlaylists      MusicSyncType = "playlists"       // User playlists
	MusicSyncTypeRecentlyPlayed MusicSyncType = "recently_played" // Listen history
)

// MusicSyncOptions defines options for music sync operations
type MusicSyncOptions struct {
	MatchThreshold float64        `json:"match_threshold"` // Minimum match confidence
	DryRun         bool           `json:"dry_run"`         // Preview mode
	ConflictPolicy ConflictPolicy `json:"conflict_policy"` // How to handle conflicts
}

// ConflictPolicy defines how to handle sync conflicts
type ConflictPolicy string

const (
	ConflictPolicySkip      ConflictPolicy = "skip"      // Skip conflicting items
	ConflictPolicyOverwrite ConflictPolicy = "overwrite" // Replace existing items
	ConflictPolicyMerge     ConflictPolicy = "merge"     // Merge metadata
)

// SupportedMusicServices defines the supported music streaming services
var SupportedMusicServices = []string{"spotify", "deezer"}

// IsSupportedMusicService checks if a service is supported for music sync
func IsSupportedMusicService(serviceName string) bool {
	for _, service := range SupportedMusicServices {
		if service == serviceName {
			return true
		}
	}
	return false
}

// GetMusicSyncTypeFromString converts string to MusicSyncType
func GetMusicSyncTypeFromString(s string) (MusicSyncType, bool) {
	switch s {
	case string(MusicSyncTypeFavorites):
		return MusicSyncTypeFavorites, true
	case string(MusicSyncTypePlaylists):
		return MusicSyncTypePlaylists, true
	case string(MusicSyncTypeRecentlyPlayed):
		return MusicSyncTypeRecentlyPlayed, true
	default:
		return "", false
	}
}

func GetMusicSyncTypeDescription() map[string]string {

	return map[string]string{
		string(MusicSyncTypeFavorites):      "Liked/saved tracks",
		string(MusicSyncTypePlaylists):      "User playlists",
		string(MusicSyncTypeRecentlyPlayed): "Recently played tracks",
	}
}
