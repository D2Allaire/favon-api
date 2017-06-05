package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AuthIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! \n")
}

// AuthTVDBToken Get a valid TVDB token for API calls.
func AuthTVDBToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "TVDB Token!\n")
}
