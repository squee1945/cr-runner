package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(0)

	logInfo("Starting server...")

	config, err := newConfig(context.Background())
	if err != nil {
		log.Fatalf("Bad config: %v", err)
	}

	// Ensure we have the Cloud Run Job created.
	job := cloudRunJob{config: config}
	if err := job.ensureJob(context.Background()); err != nil {
		log.Fatalf("Failed to create Cloud Run job %q: %v", config.JobID, err)
	}

	// Start HTTP server.
	http.HandleFunc("/app/token", func(w http.ResponseWriter, r *http.Request) {
		apphandler{w: w, r: r, config: config}.next()
	})
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		handler{w: w, r: r, config: config}.next()
	})
	logInfo("Listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
