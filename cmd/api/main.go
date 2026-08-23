package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sabbir185/gpc/config"
)

func main() {
	cnf := config.LoadConfig()
	log.Println(cnf.App.Name)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":" + cnf.App.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Starting server on port %s", cnf.App.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe error: %v", err.Error())
	}
}
