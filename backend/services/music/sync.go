package music

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"syncer.net/core/services"
	"syncer.net/core/sync"
)

// MusicSyncComponents holds music-specific sync implementations
type MusicSyncComponents struct {
	Transformer *MusicTrackTransformer
	Adder       *MusicCrossServiceAdder
}

// NewMusicSyncComponents creates music-specific sync implementations
func NewMusicSyncComponents() *MusicSyncComponents {
	return &MusicSyncComponents{
		Transformer: NewMusicTrackTransformer(),
		Adder:       NewMusicCrossServiceAdder(),
	}
}

// CreateMusicSyncEngine creates a sync engine configured for music services
func CreateMusicSyncEngine(
	registry *services.ServiceRegistry,
	oauth *services.OAuthManager,
	db *sqlx.DB,
	workers int,
) (*sync.SyncEngine, error) {
	components := NewMusicSyncComponents()

	// Create the sync engine with music-specific components
	engine := sync.NewSyncEngine(
		oauth,
		components.Transformer, // Implements DataTransformer
		components.Adder,       // Implements CrossServiceAdder
		db,
		workers,
	)

	return engine, nil
}

// GetSupportedSyncTypes returns the sync types supported by music services
func GetSupportedSyncTypes() []string {
	return []string{
		string(MusicSyncTypeFavorites),
		string(MusicSyncTypePlaylists),
		string(MusicSyncTypeRecentlyPlayed),
	}
}

// ValidateMusicSyncRequest validates a sync request for music services
func ValidateMusicSyncRequest(servicePairs []string, syncType string) error {
	// Validate that all services are music services
	for _, service := range servicePairs {
		if !IsSupportedMusicService(service) {
			return fmt.Errorf("unsupported music service: %s", service)
		}
	}

	// Validate sync type
	_, valid := GetMusicSyncTypeFromString(syncType)
	if !valid {
		return fmt.Errorf("unsupported sync type for music services: %s", syncType)
	}

	return nil
}
