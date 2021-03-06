package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// chi router has pattern matching
//returns pointer to the julien router
func (app *application) routes() http.Handler{
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMoviesByGenre)
	router.HandlerFunc(http.MethodPost, "/v1/admin/insertmovie", app.InsertMovie)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deletemovie/:id", app.DeleteMovie)
	return app.enableCORS(router)
}