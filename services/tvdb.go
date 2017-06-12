package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"favon-api/models"

	"favon-api/utils"

	"github.com/go-redis/redis"
)

// TVDB is the model for the TVDB connection
type TVDB struct {
	client *redis.Client
	Token  *models.TVDBToken
}

// NewTVDB creates a new TVDB service instance
func NewTVDB(client *redis.Client) *TVDB {
	tvdb := TVDB{}
	tvdb.client = client
	tvdb.Token = &models.TVDBToken{}
	return &tvdb
}

// GetToken returns a valid TVDB token
func (tvdb *TVDB) GetToken() {
	err := tvdb.Token.Load(tvdb.client)
	// Key not found (expired), request new token and store
	if err == redis.Nil {
		tvdb.requestToken()
		tvdb.Token.Save(tvdb.client)
	} else if err != nil {
		panic(err)
	}
}

// Request a new Token from the API
func (tvdb *TVDB) requestToken() {
	key := os.Getenv("TVDB_KEY")
	cred := make(map[string]string)
	cred["apikey"] = key
	body, _ := json.Marshal(cred)
	req, err := http.NewRequest("POST", "https://api.thetvdb.com/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		utils.ErrorLog.Println("TVDB requestTokenRequest", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLog.Println("TVDB requestToken: ", err)
		return
	}

	defer resp.Body.Close()

	rbody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rbody, tvdb.Token)
}

// Search for series based on passed query name
func (tvdb *TVDB) Search(query string) []byte {
	req, err := http.NewRequest("GET", "https://api.thetvdb.com/search/series?name="+query, nil)
	if err != nil {
		utils.ErrorLog.Println("TVDB SearchRequest: ", err)
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+tvdb.Token.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLog.Println("TVDB Search: ", err)
		return nil
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

// Get episode info for specified series by ID and paginated
func (tvdb *TVDB) Get(id string) *models.Series {
	series := new(models.Series)
	err := series.Get(id, tvdb.client)
	if err == nil {
		return series
	}

	start := 1
	page := &start

	// Request series information
	tvdb.getSeriesInfo(id, series)
	utils.InfoLog.Printf("Requesting %s", series.Name)

	// Get all episodes by page until all pages have been requested
	for page != nil && *page > 0 {
		response := tvdb.getEpisodes(id, page)
		page = response.Links.Next
		series.Episodes = append(series.Episodes, response.Data...)
	}

	if series.IsFinished() {
		series.Save(tvdb.client)
	}

	return series
}

func (tvdb *TVDB) getSeriesInfo(id string, series *models.Series) {
	req, err := http.NewRequest("GET", "https://api.thetvdb.com/series/"+id, nil)
	if err != nil {
		utils.ErrorLog.Println("TVDB getSeriesInfoRequest: ", err)
	}
	req.Header.Set("Authorization", "Bearer "+tvdb.Token.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLog.Println("TVDB getSeriesInfo: ", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := new(models.SeriesAPIResponse)
	json.Unmarshal(body, &response)
	*series = response.Series
}

func (tvdb *TVDB) getEpisodes(id string, page *int) *models.EpisodesAPIResponse {
	url := fmt.Sprint("https://api.thetvdb.com/series/", id, "/episodes?page=", *page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.ErrorLog.Println("TVDB getEpisodesRequest: ", err)
	}
	req.Header.Set("Authorization", "Bearer "+tvdb.Token.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLog.Println("TVDB getEpisodes: ", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := new(models.EpisodesAPIResponse)
	json.Unmarshal(body, &response)
	return response
}
