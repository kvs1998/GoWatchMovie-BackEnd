package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request){
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}
	app.logger.Println("id is", id)
	movie, err := app.models.DB.GetOneMovie(id)
	//movie := models.Movie{
	//	ID: id,
	//	Title: "Some Movie",
	//	Description: "desc",
	//	Year: 1234,
	//	ReleaseDate: time.Date(2021,01,01,01,0,0,0, time.Local),
	//	Runtime: 234,
	//	Rating: 4,
	//	CreatedAt: time.Date(2021,01,01,01,0,0,0, time.Local),
	//	UpdatedAt: time.Date(2021,01,01,01,0,0,0, time.Local),
	//}

	err = app.writeJSON(w,http.StatusOK, movie, "LatestMovie")

}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request){
	allmovie, _ := app.models.DB.AllMovies()
	err := app.writeJSON(w,http.StatusOK, allmovie, "AllMovies")
	if err != nil{
		app.logger.Println("error")
		return
	}
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request){
	allgenres, _ := app.models.DB.AllGenres()
	err := app.writeJSON(w,http.StatusOK, allgenres, "AllGenres")
	if err != nil{
		app.logger.Println("error")
		return
	}
}

func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request){	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}
	app.logger.Println("genre id is", id)
	allgenres, _ := app.models.DB.AllMoviesByGenre(id)
	err = app.writeJSON(w,http.StatusOK, allgenres, "MoviesByGenre")
	if err != nil{
		app.logger.Println("error")
		return
	}
}

type MoviePayload struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Year string `json:"year"`
	ReleaseDate string`json:"release_date"`
	Runtime string `json:"runtime"`
	Rating string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request){
	var req MoviePayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.logger.Println(errors.New("Decoder error"))
		app.errorJSON(w, err)
		return
	}
	var movie models.Movie
	movie.ID,_ = strconv.Atoi(req.ID)
	movie.Title = req.Title
	movie.Description = req.Description
	movie.Year,_ = strconv.Atoi(req.Year)
	movie.ReleaseDate, _ = time.Parse("2006-01-02", req.ReleaseDate)
	movie.Rating,_ = strconv.Atoi(req.Rating)
	movie.Runtime,_ = strconv.Atoi(req.Runtime)
	movie.MPAARating = req.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.ID == 0 {
		err = app.models.DB.InsertMovie(movie)
		if err != nil {
			app.logger.Println(errors.New("Insert movie error"))
			app.errorJSON(w, err)
			return
		}
	}

	ok := jsonResp{
		OK: true,
	}
	err = app.writeJSON(w,http.StatusOK, ok, "response")
	if err != nil {
		app.logger.Println(errors.New("api error"))
		app.errorJSON(w, err)
		return
	}
}


func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request){
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	log.Println("DELETE", id)
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}
	err = app.models.DB.DeleteMovie(id)
	if err != nil {
		app.logger.Println(errors.New("Delete movie error"))
		app.errorJSON(w, err)
		return
	}
	ok := jsonResp{
		OK: true,
	}
	err = app.writeJSON(w,http.StatusOK, ok, "response")
	if err != nil {
		app.logger.Println(errors.New("api error"))
		app.errorJSON(w, err)
		return
	}
}