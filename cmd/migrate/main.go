package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sabbir185/geopunch/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <create <name> | up | down | force <version>>", os.Args[0])
	}

	cnf := config.LoadConfig()

	// Create sql file
	if os.Args[1] == "create" {
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s create <migration_name>", os.Args[0])
		}
		name := strings.Join(os.Args[2:], "_")
		createMigration(cnf.DB.MigrationsPath, name)
		return
	}

	// Initialize migration
	m, err := migrate.New("file://"+cnf.DB.MigrationsPath, cnf.DB.Url)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("database close error: %v", dbErr)
		}
	}()

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("m.Up error: %v", err)
		}
		fmt.Println("Migration up: completed successfully")

	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) && !errors.Is(err, migrate.ErrNilVersion) {
			log.Fatalf("m.Down error: %v", err)
		}
		fmt.Println("Migration down: completed successfully")

	case "force":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s force <version>", os.Args[0])
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("m.Force error: %v", err)
		}
		fmt.Printf("Migration force: forced to version %d successfully\n", version)

	default:
		log.Fatalf("Invalid command: %s. Use create, up, down, or force", os.Args[1])
	}
}

// Migration file create
func createMigration(dir string, name string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	maxSeq := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) > 0 {
			if seq, err := strconv.Atoi(parts[0]); err == nil && seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	nextSeq := maxSeq + 1
	cleanName := strings.TrimSpace(strings.ToLower(name))
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "-", "_")

	if cleanName == "" {
		log.Fatalf("Migration name cannot be empty")
	}

	upFile := filepath.Join(dir, fmt.Sprintf("%06d_%s.up.sql", nextSeq, cleanName))
	downFile := filepath.Join(dir, fmt.Sprintf("%06d_%s.down.sql", nextSeq, cleanName))

	if err := os.WriteFile(upFile, []byte("-- Write your UP migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create %s: %v", upFile, err)
	}
	if err := os.WriteFile(downFile, []byte("-- Write your DOWN migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create %s: %v", downFile, err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
}
