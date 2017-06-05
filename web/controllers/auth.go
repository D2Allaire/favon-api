package controllers

import (
	"fmt"
	"net/http"

	"favon-api/services"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	redis *redis.Client
}

func NewAuthController(redisClient *redis.Client) *AuthController {
	return &AuthController{
		redis: redisClient,
	}
}

func AuthIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! \n")
}

// GetTVDBToken Get a valid TVDB token for API calls.
func (c *AuthController) GetTVDBToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := services.TVDBToken{}
	token.Load(c.redis)
	toJSON, _ := json.Marshal(token)
	fmt.Fprint(w, string(toJSON))
}
