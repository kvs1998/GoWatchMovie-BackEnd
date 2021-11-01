package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

// version of application
const version = "1.0.0"

const (
	host     = "localhost"
	port     = 6000
	user     = "cal"
	password = "newport"
	dbname   = "movies"
)

type config struct {
	port int
	env string
	db struct{
		//connection string for database
		dsn string
	}
}

type AppStatus struct {
	Status string `json:"status"`
	Environment string `json:"environment"`
	Version string `json:"version"`
}

//hold application configuration
type application struct {
	config config
	logger *log.Logger //pointer to built in package
	models models.Models
}

func main() {
	var cfg config

	// read input from command line
	//"postgres://cal:newport@localhost:5432/movies?sslmode=disable"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application env development|Production")
	flag.StringVar(&cfg.db.dsn, "dsn", psqlInfo, "Postgres connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Println("ERROR")
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	//start web server
	//err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
	srv := &http.Server {
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(), //handler to invoke, http.DefaultServeMux if nil
		IdleTimeout: time.Minute,
		ReadHeaderTimeout: 10* time.Second,
		WriteTimeout: 30*time.Second,
	}
	//ListenAndServe is built in method
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		log.Println(" sql open  ERROR")
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	err = db.PingContext(ctx)
	defer cancel()
	if err != nil {
		log.Println(" ping error")
		return nil, err
	}
	return db, nil
}