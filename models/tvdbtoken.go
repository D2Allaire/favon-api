package models

import (
	"encoding/json"

	"time"

	"github.com/go-redis/redis"
)

// TVDBToken is the Auth token for TVDB
type TVDBToken struct {
	Token string `json:"token"`
}

// Load an existing TVDB token from database or request new one
func (t *TVDBToken) Load(client *redis.Client) error {
	data, err := client.Get("tvdb").Result()

	// If redis returns nil the key has expired. Request new token.
	if err != nil {
		return err
	}
	// If no error occured, unserialize TVDBToken.
	err = json.Unmarshal([]byte(data), t)
	if err != nil {
		panic(err)
	}
	return nil
}

// Save token to redis
func (t *TVDBToken) Save(client *redis.Client) {
	serialized, _ := json.Marshal(t)
	err := client.Set("tvdb", serialized, time.Duration(24*time.Hour)).Err()
	if err != nil {
		panic(err)
	}
}
