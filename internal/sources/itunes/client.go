package itunes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/khatibomar/gomania/internal/sources"
)

var _ sources.Client = (*Client)(nil)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type SearchResponse struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

type Result struct {
	TrackID          int    `json:"trackId"`
	TrackName        string `json:"trackName"`
	ArtistName       string `json:"artistName"`
	CollectionName   string `json:"collectionName"`
	TrackViewURL     string `json:"trackViewUrl"`
	ArtworkURL100    string `json:"artworkUrl100"`
	ArtworkURL600    string `json:"artworkUrl600"`
	ReleaseDate      string `json:"releaseDate"`
	TrackTimeMillis  int    `json:"trackTimeMillis"`
	Country          string `json:"country"`
	PrimaryGenreName string `json:"primaryGenreName"`
	Description      string `json:"description"`
	ShortDescription string `json:"shortDescription"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://itunes.apple.com/search",
	}
}

func (c *Client) SearchPodcasts(term string, limit int) ([]sources.Podcast, error) {
	if limit == 0 {
		limit = 50
	}

	params := url.Values{}
	params.Set("term", term)
	params.Set("media", "podcast")
	params.Set("limit", strconv.Itoa(limit))

	searchURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes API returned status: %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert iTunes results to common Podcast format
	podcasts := make([]sources.Podcast, 0, len(searchResp.Results))
	for _, result := range searchResp.Results {
		podcast := sources.Podcast{
			ID:          strconv.Itoa(result.TrackID),
			Title:       result.TrackName,
			Description: c.getDescription(result),
			Host:        result.ArtistName,
			Genre:       result.PrimaryGenreName,
			Country:     result.Country,
			Duration:    result.ToDuration(),
			PublishedAt: c.parsePublishedAt(result),
			ArtworkURL:  result.ArtworkURL600,
			ExternalURL: result.TrackViewURL,
			SourceName:  "itunes",
			ExternalID:  strconv.Itoa(result.TrackID),
		}
		podcasts = append(podcasts, podcast)
	}

	return podcasts, nil
}

// ToDuration converts the track time in milliseconds to seconds.
func (r *Result) ToDuration() int {
	return r.TrackTimeMillis / 1000
}

func (c *Client) GetSourceName() string {
	return "itunes"
}

func (c *Client) getDescription(result Result) string {
	if result.Description != "" {
		return result.Description
	}
	if result.ShortDescription != "" {
		return result.ShortDescription
	}
	return "Podcast by %s" + result.ArtistName
}

func (c *Client) parsePublishedAt(result Result) *time.Time {
	if result.ReleaseDate == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", result.ReleaseDate)
	if err != nil {
		return nil
	}
	return &t
}

func (r *Result) ToPublishedAt() (*time.Time, error) {
	if r.ReleaseDate == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", r.ReleaseDate)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
