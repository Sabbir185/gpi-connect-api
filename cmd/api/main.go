package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/infra/db"
	"github.com/Sabbir185/gpi/infra/redis"
	"github.com/Sabbir185/gpi/internal"
)

func main() {
	cnf := config.LoadConfig()

	db, err := db.Connect(cnf.DB.Url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database successfully")

	redisClient, err := redis.Connect(cnf.DB.RedisUrl)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Connected to redis successfully")

	// Register all the routes
	multiplexer := http.NewServeMux()
	mux := internal.RegisterRoutes(multiplexer)

	server := &http.Server{
		Addr:         ":" + cnf.App.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// Graceful Shutdown, error handling or listen and serve
	go func() {
		log.Printf("Starting server on port %s", cnf.App.Port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Server is shutting down with signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:\n%v\n", err.Error())
	}

	log.Println("Server exited properly")
}
