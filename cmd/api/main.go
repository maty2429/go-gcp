package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"task-api/internal/handler"
	"task-api/internal/repository"
)

func main() {
	ctx := context.Background()

	projectID := "project-4bf3e3b5-7b1b-4eb3-8a3"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error al conectar a Firestore: %v", err)
	}
	defer client.Close()

	httpAddr := getHTTPAddr()
	taskRepo := repository.NewTaskFirestoreRepository(client)

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
