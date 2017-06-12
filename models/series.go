package models

import (
	"encoding/json"

	"fmt"

	"github.com/go-redis/redis"
)

// Episode model for tv show
type Episode struct {
	Episode uint8  `json:"airedEpisodeNumber"`
	Season  uint8  `json:"airedSeason"`
	Title   string `json:"episodeName"`
}

// EpisodesAPIResponse models the JSON response for TVDB /episodes request
type EpisodesAPIResponse struct {
	Links struct {
		Next *int `json:"next"`
	} `json:"links"`
	Data []Episode `json:"data"`
}

// SeriesAPIResponse models the JSON response TVDB show request
type SeriesAPIResponse struct {
	Series Series `json:"data"`
}

// Series model for tv show
type Series struct {
	ID       int      `json:"id"`
	Name     string   `json:"seriesName"`
	Aliases  []string `json:"aliases"`
	Status   string   `json:"status"`
	Year     string   `json:"firstAired"`
	Episodes []Episode
}

// IsFinished checks if a show has finished airing
func (s *Series) IsFinished() bool {
	return s.Status != "Continuing"
}

// Get a series from database
func (s *Series) Get(id string, client *redis.Client) error {
	data, err := client.HGet("series", id).Result()

	// If redis returns nil the series has not been saved in our DB
	if err != nil {
		return err
	}
	// If no error occured, unserialize series
	err = json.Unmarshal([]byte(data), s)
	if err != nil {
		panic(err)
	}
	return nil
}

// Save a series to redis
func (s *Series) Save(client *redis.Client) {
	serialized, _ := json.Marshal(s)
	id := fmt.Sprintf("%d", s.ID)
	err := client.HSet("series", id, serialized).Err()
	if err != nil {
		panic(err)
	}
}
