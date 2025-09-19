package services

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
)

// ServiceRegistry manages all available service providers
type ServiceRegistry struct {
	services map[string]ServiceProvider
	mu       sync.RWMutex
	db       *sqlx.DB
	logger   *log.Logger
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(db *sqlx.DB, logger *log.Logger) *ServiceRegistry {
	if logger == nil {
		logger = log.New(log.Writer(), "[ServiceRegistry] ", log.LstdFlags)
	}

	registry := &ServiceRegistry{
		services: make(map[string]ServiceProvider),
		db:       db,
		logger:   logger,
	}

	return registry
}

// Register adds a new service provider to the registry
func (r *ServiceRegistry) Register(service ServiceProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := service.Name()
	if _, exists := r.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	r.services[name] = service
	r.logger.Printf("Registered service: %s (category: %s)", name, service.Category())

	return nil
}

// GetService retrieves a service provider by name
func (r *ServiceRegistry) GetService(name string) (ServiceProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// ListServices returns information about all registered services
func (r *ServiceRegistry) ListServices() []ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []ServiceInfo
	for _, service := range r.services {

		services = append(services, ServiceInfo{
			Name:        service.Name(),
			DisplayName: service.DisplayName(),
			Category:    service.Category(),
			Scopes:      service.RequiredScopes(),
			Available:   service.HealthCheck() == nil,
		})
	}

	return services
}

// GetServicesByCategory returns all services in a specific category
func (r *ServiceRegistry) GetServicesByCategory(category ServiceCategory) []ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []ServiceInfo
	for _, service := range r.services {
		if service.Category() == category {

			services = append(services, ServiceInfo{
				Name:        service.Name(),
				DisplayName: service.DisplayName(),
				Category:    service.Category(),
				Scopes:      service.RequiredScopes(),
				Available:   service.HealthCheck() == nil,
			})
		}
	}

	return services
}

// IsServiceAvailable checks if a service is registered and available
func (r *ServiceRegistry) IsServiceAvailable(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.services[name]
	return exists && r.services[name].HealthCheck() == nil
}

// GetServiceCount returns the number of registered services
func (r *ServiceRegistry) GetServiceCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.services)
}

// Unregister removes a service from the registry
func (r *ServiceRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[name]; !exists {
		return fmt.Errorf("service %s not found", name)
	}

	delete(r.services, name)
	r.logger.Printf("Unregistered service: %s", name)

	return nil
}
