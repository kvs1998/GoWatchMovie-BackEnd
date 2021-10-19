package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// chi router has pattern matching
//returns pointer to the julien router
func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	return router
}