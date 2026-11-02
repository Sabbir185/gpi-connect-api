.PHONY: help build run worker keygen clean migrate-create migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build               - Build the application binary"
	@echo "  make run                 - Build and start the application"
	@echo "  make worker              - Build and start the background worker (emails, push)"
	@echo "  make keygen              - Print a random secret for JWT_SECRET / OTP_SECRET"
	@echo "  make clean               - Remove the build artifacts"
	@echo "  make migrate-create name= - Create a new migration (e.g. make migrate-create name=users_table)"
	@echo "  make migrate-up          - Run database migrations"
	@echo "  make migrate-down        - Rollback database migrations"

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

worker:
	@go build -o bin/worker ./cmd/worker
	@./bin/worker

keygen:
	@go run ./cmd/tools/keygen

clean:
	@rm -rf bin

migrate-create:
	@go run ./cmd/migrate create $(name)

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

