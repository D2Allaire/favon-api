package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type TVDBToken struct {
	Token   string `json:"token"`
	Expires time.Time
}

func (t *TVDBToken) Load(redis *redis.Client) {
	data, err := redis.Get("tvdb").Result()
	// Error hopefully only occurs when key is not set. In that case, set it.
	if err != nil {
		t.requestToken()
		t.save(redis)
		return
	}
	// If no error occured, unserialize TVDBToken.
	err = json.Unmarshal([]byte(data), t)
	if err != nil {
		panic(err)
	}
	// Check if token is expiring soon, if so request new one.
	if t.Expires.Sub(time.Now()).Hours() < float64(1) {
		t.requestToken()
		t.save(redis)
	}
}

func (t *TVDBToken) requestToken() {
	key := os.Getenv("TVDB_KEY")
	cred := make(map[string]string)
	cred["apikey"] = key
	body, _ := json.Marshal(cred)
	req, err := http.NewRequest("POST", "https://api.thetvdb.com/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	rbody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rbody, t)
	t.Expires = time.Now().Add(24 * time.Hour)
}

func (t *TVDBToken) save(redis *redis.Client) {
	serialized, _ := json.Marshal(t)
	err := redis.Set("tvdb", serialized, 0).Err()
	if err != nil {
		panic(err)
	}
}
