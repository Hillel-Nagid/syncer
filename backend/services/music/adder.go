package music

import (
	"context"
	"fmt"
	"log"

	"syncer.net/core/services"
	"syncer.net/core/sync"
)

// MusicCrossServiceAdder handles adding music tracks to different services
// Implements the CrossServiceAdder interface from core/sync
type MusicCrossServiceAdder struct {
	logger *log.Logger
}

// NewMusicCrossServiceAdder creates a new music cross-service adder
func NewMusicCrossServiceAdder() *MusicCrossServiceAdder {
	return &MusicCrossServiceAdder{
		logger: log.New(log.Writer(), "[MusicAdder] ", log.LstdFlags),
	}
}

// AddItemToService adds a universal music track to a target service
func (a *MusicCrossServiceAdder) AddItemToService(
	ctx context.Context,
	targetService services.ServiceProvider,
	tokens *services.OAuthTokens,
	universalItem sync.UniversalItem,
	options any,
) error {
	// Convert to UniversalTrack
	track, ok := universalItem.(UniversalTrack)
	if !ok {
		return fmt.Errorf("item is not a UniversalTrack, got %T", universalItem)
	}

	// Convert options to MusicSyncOptions
	var musicOptions MusicSyncOptions
	if opts, ok := options.(MusicSyncOptions); ok {
		musicOptions = opts
	} else {
		// Use default options
		musicOptions = MusicSyncOptions{
			MatchThreshold: 0.8,
			DryRun:         false,
			ConflictPolicy: ConflictPolicySkip,
		}
	}

	serviceName := targetService.Name()

	// Check if we already have the track ID for this service
	if existingID, exists := track.ExternalIDs[serviceName]; exists && existingID != "" {
		a.logger.Printf("Track already exists in %s with ID: %s", serviceName, existingID)
		return nil // Consider this a success since the track already exists
	}

	if musicOptions.DryRun {
		a.logger.Printf("DRY RUN: Would add track '%s' by '%s' to %s", track.Title, track.Artist, serviceName)
		return nil
	}

	// For now, return a placeholder error indicating this needs service-specific implementation
	// This is where we would implement search and add functionality for each service
	a.logger.Printf("Adding track '%s' by '%s' to %s - service-specific implementation needed",
		track.Title, track.Artist, serviceName)

	return fmt.Errorf("cross-service track addition for %s not yet fully implemented - needs service-specific search and add logic", serviceName)
}

// SearchAndAddTrack searches for a track on the target service and adds it
func (a *MusicCrossServiceAdder) SearchAndAddTrack(
	ctx context.Context,
	targetService services.ServiceProvider,
	tokens *services.OAuthTokens,
	track UniversalTrack,
	options MusicSyncOptions,
) error {
	serviceName := targetService.Name()

	if !IsSupportedMusicService(serviceName) {
		return fmt.Errorf("unsupported music service: %s", serviceName)
	}

	// This would be implemented with service-specific search and add logic
	// For now, we'll return a placeholder
	a.logger.Printf("Would search for track '%s' by '%s' on %s and add it",
		track.Title, track.Artist, serviceName)

	return fmt.Errorf("search and add functionality not yet implemented for %s", serviceName)
}

// ValidateTrackForService checks if a track can be added to the target service
func (a *MusicCrossServiceAdder) ValidateTrackForService(track UniversalTrack, serviceName string) error {
	if !IsSupportedMusicService(serviceName) {
		return fmt.Errorf("unsupported music service: %s", serviceName)
	}

	if track.Title == "" {
		return fmt.Errorf("track title is required")
	}

	if track.Artist == "" {
		return fmt.Errorf("track artist is required")
	}

	return nil
}
