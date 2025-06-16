package sources

import (
	"context"
	"fmt"
)

// Manager handles multiple external sources for podcast content
type Manager struct {
	clients map[string]Client
}

// NewManager creates a new sources manager
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]Client),
	}
}

// RegisterClient adds a new external source client
func (m *Manager) RegisterClient(client Client) {
	m.clients[client.GetSourceName()] = client
}

// GetClient returns a specific client by source name
func (m *Manager) GetClient(sourceName string) (Client, bool) {
	client, exists := m.clients[sourceName]
	return client, exists
}

// SearchAllSources searches across all registered sources
func (m *Manager) SearchAllSources(ctx context.Context, term string, limit int) (map[string][]Podcast, error) {
	results := make(map[string][]Podcast)

	for sourceName, client := range m.clients {
		podcasts, err := client.SearchPodcasts(ctx, term, limit)
		if err != nil {
			// Log error but continue with other sources
			continue
		}
		results[sourceName] = podcasts
	}

	return results, nil
}

// SearchBySource searches a specific source
func (m *Manager) SearchBySource(ctx context.Context, sourceName, term string, limit int) ([]Podcast, error) {
	client, exists := m.GetClient(sourceName)
	if !exists {
		return nil, fmt.Errorf("source '%s' not found", sourceName)
	}

	return client.SearchPodcasts(ctx, term, limit)
}

// GetAvailableSources returns list of registered source names
func (m *Manager) GetAvailableSources() []string {
	sources := make([]string, 0, len(m.clients))
	for name := range m.clients {
		sources = append(sources, name)
	}
	return sources
}
