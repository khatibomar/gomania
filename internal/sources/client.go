package sources

import "time"

type Podcast struct {
	ID          string
	Title       string
	Description string
	Host        string
	Genre       string
	Country     string
	Duration    int // in seconds
	PublishedAt *time.Time
	ArtworkURL  string
	ExternalURL string
	SourceName  string // "itunes", "spotify", etc.
	ExternalID  string // external platform's ID
}

type Client interface {
	SearchPodcasts(term string, limit int) ([]Podcast, error)
	GetSourceName() string
}
