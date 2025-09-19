package services

import (
	"log"

	"syncer.net/core/services"
	"syncer.net/services/music/deezer"
	"syncer.net/services/music/spotify"
)

// InitializeServices registers all available service providers with the registry
// This is the ONLY place where specific service implementations are imported
// and connected to the generic service registry
func InitializeServices(registry *services.ServiceRegistry, logger *log.Logger) error {
	if logger == nil {
		logger = log.New(log.Writer(), "[ServiceInit] ", log.LstdFlags)
	}

	logger.Println("Initializing and registering service providers...")

	// Register music services
	if err := registerMusicServices(registry, logger); err != nil {
		return err
	}

	// Future: Register calendar services
	// if err := registerCalendarServices(registry, logger); err != nil {
	//     return err
	// }

	logger.Printf("Successfully registered %d services", registry.GetServiceCount())
	return nil
}

// registerMusicServices registers all music streaming service providers
func registerMusicServices(registry *services.ServiceRegistry, logger *log.Logger) error {
	// Register Spotify service
	spotifyService := spotify.NewSpotifyService()
	if err := registry.Register(spotifyService); err != nil {
		logger.Printf("Failed to register Spotify service: %v", err)
		return err
	}
	logger.Printf("Registered Spotify music service")

	// Register Deezer service
	deezerService := deezer.NewDeezerService()
	if err := registry.Register(deezerService); err != nil {
		logger.Printf("Failed to register Deezer service: %v", err)
		return err
	}
	logger.Printf("Registered Deezer music service")

	return nil
}

// GetAllRegisteredServices returns information about all registered services
// This can be used for debugging and service discovery
func GetAllRegisteredServices(registry *services.ServiceRegistry) []services.ServiceInfo {
	return registry.ListServices()
}

// GetMusicServices returns only music streaming services
func GetMusicServices(registry *services.ServiceRegistry) []services.ServiceInfo {
	return registry.GetServicesByCategory(services.CategoryMusic)
}
