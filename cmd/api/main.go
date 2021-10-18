package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// version of application
const version = "1.0.0"

type config struct {
	port int
	env string
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
}

func main() {
	var cfg config

	// read input from command line
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application env development|Production")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger:logger,
	}

	//start web server
	//err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
	srv := &http.Server {
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadHeaderTimeout: 10* time.Second,
		WriteTimeout: 30*time.Second,
	}
	//ListenAndServe is built in method
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
