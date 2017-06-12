package controllers

import (
	"favon-api/services"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// AuthController Controller for authentication routes
type AuthController struct {
	redis *redis.Client
	tvdb  *services.TVDB
}

// NewAuthController creates a new Auth Controller instance
func NewAuthController(redisClient *redis.Client) *AuthController {
	a := AuthController{}
	a.redis = redisClient
	a.tvdb = services.NewTVDB(redisClient)
	return &a
}

// GetTVDBToken Get a valid TVDB token for API calls.
func (a *AuthController) GetTVDBToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.tvdb.GetToken()
	toJSON, _ := json.Marshal(a.tvdb.Token)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(toJSON))
}
