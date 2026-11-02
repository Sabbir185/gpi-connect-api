package main

import (
	"context"
	"log"

	"github.com/Sabbir185/gpi/config"
	"github.com/Sabbir185/gpi/infra/db"
	"github.com/Sabbir185/gpi/infra/email"
	"github.com/Sabbir185/gpi/infra/firebase"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/internal/worker"
)

const concurrency = 10

func main() {
	cnf := config.LoadConfig()

	db, err := db.Connect(cnf.DB.Url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	firebaseClient, err := firebase.NewClient(context.Background(), cnf.Firebase.CredentialsBase64)
	if err != nil {
		log.Fatalf("Failed to init firebase: %v", err)
	}

	mailer := email.NewBrevoClient(cnf.Brevo.APIKey, cnf.Brevo.SenderEmail, cnf.Brevo.SenderName)
	if !mailer.Configured() {
		log.Println("BREVO_API_KEY or BREVO_SENDER_EMAIL is not set, welcome emails will be skipped")
	}

	// Run blocks until SIGINT/SIGTERM and lets in-flight tasks finish.
	log.Println("Starting worker")
	w := worker.New(cnf.DB.RedisUrl, concurrency, mailer, firebaseClient.Messaging, device.NewDeviceRepository(db))
	if err := w.Run(); err != nil {
		log.Fatalf("Worker stopped: %v", err)
	}
}
