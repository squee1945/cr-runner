package main

import (
	"log"
	"net/http"
	"os"
)

const (
	eventHeader    = "X-GitHub-Event"
	deliveryHeader = "X-GitHub-Delivery"
	sigHeader      = "X-Hub-Signature"
	sig256Header   = "X-Hub-Signature-256"
)

func main() {
	logInfo("Starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logInfo("Defaulting to port %s", port)
	}

	// Start HTTP server.
	logInfo("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logWarn("Bad method %v", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ev, err := parseEvent(r.Body)
	if err != nil {
		serverError(w, "parsing event: %v", err)
		return
	}

	logInfo("Event received:\n%v\n", ev)
}

func logInfo(fmt string, args ...any) {
	log.Printf("I "+fmt+"\n", args...)
}

func logWarn(fmt string, args ...any) {
	log.Printf("W "+fmt+"\n", args...)
}

func logError(fmt string, args ...any) {
	log.Printf("E "+fmt+"\n", args...)
}

func serverError(w http.ResponseWriter, fmt string, args ...any) {
	logError("Error: "+fmt, args...)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Server error"))
}

func clientError(w http.ResponseWriter, fmt string, args ...any) {
	logInfo("Client error: "+fmt, args...)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Client error"))
}
