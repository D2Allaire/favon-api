package controllers

import (
	"net/http"

	"favon-api/services"

	"fmt"

	"net/url"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// SeriesController Controller for tv show routes
type SeriesController struct {
	redis *redis.Client
	tvdb  *services.TVDB
}

// NewSeriesController creates a new Series Controller instance
func NewSeriesController(redisClient *redis.Client) *SeriesController {
	s := SeriesController{}
	s.redis = redisClient
	s.tvdb = services.NewTVDB(redisClient)
	return &s
}

// Search for a show based on queried name
func (s *SeriesController) Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.tvdb.GetToken()
	queryValues := r.URL.Query()
	data := s.tvdb.Search(url.QueryEscape(queryValues.Get("name")))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data != nil {
		fmt.Fprint(w, string(data))
	}
}

// Get episodes of specified show
func (s *SeriesController) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.tvdb.GetToken()
	id := p.ByName("id")
	series := s.tvdb.Get(id)
	body, _ := json.Marshal(series)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(body))
}
