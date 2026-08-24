.PHONY: help build run clean migrate-create migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build               - Build the application binary"
	@echo "  make run                 - Build and start the application"
	@echo "  make clean               - Remove the build artifacts"
	@echo "  make migrate-create name= - Create a new migration (e.g. make migrate-create name=users_table)"
	@echo "  make migrate-up          - Run database migrations"
	@echo "  make migrate-down        - Rollback database migrations"

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

clean:
	@rm -rf bin

migrate-create:
	@go run ./cmd/migrate create $(name)

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

