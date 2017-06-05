package web

import (
	"log"
	"net/http"

	"favon-api/web/controllers"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// Init creates the server and loads application routes
func Init(redisClient *redis.Client) {
	router := httprouter.New()
	routes(router, redisClient)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func routes(router *httprouter.Router, redisClient *redis.Client) {
	// Auth routes
	authController := controllers.NewAuthController(redisClient)
	router.GET("/auth/tvdb", authController.GetTVDBToken)
}
