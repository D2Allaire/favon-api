package web

import (
	"log"
	"net/http"

	"favon-api/web/controllers"

	"github.com/julienschmidt/httprouter"
)

// Init creates the server and loads application routes
func Init() {
	router := httprouter.New()
	routes(router)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func routes(router *httprouter.Router) {
	// Auth routes
	router.GET("/auth/tvdb", controllers.AuthTVDBToken)
}
