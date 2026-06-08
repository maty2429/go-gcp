package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"task-api/internal/db"
	"task-api/internal/handler"
	"task-api/internal/repository"
)

func main() {
	ctx := context.Background()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	httpAddr := getHTTPAddr()

	var taskRepo *repository.TaskRepository
	if databaseURL == "" {
		log.Printf("DATABASE_URL not set, starting without database")
	} else {
		pool, err := db.Connect(ctx, databaseURL)
		if err != nil {
			log.Printf("database unavailable, starting without database: %v", err)
		} else {
			defer pool.Close()

			if err := db.InitSchema(ctx, pool); err != nil {
				log.Printf("could not initialize schema, starting without database: %v", err)
			} else {
				taskRepo = repository.NewTaskRepository(pool)
			}
		}
	}

	taskHandler := handler.NewTaskHandler(taskRepo)

	mux := http.NewServeMux()
	taskHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("task-api listening on %s", httpAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getHTTPAddr() string {
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		return ":" + port
	}

	return getEnv("HTTP_ADDR", ":8080")
}
