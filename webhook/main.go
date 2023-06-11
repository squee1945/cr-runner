package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kr/pretty"
)

const (
	eventHeader    = "X-GitHub-Event"
	deliveryHeader = "X-GitHub-Delivery"
	sigHeader      = "X-Hub-Signature"
	sig256Header   = "X-Hub-Signature-256"
)

func main() {
	log.SetFlags(0)

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
		clientError(w, "bad method %v", r.Method)
		return
	}

	ev, err := parseEvent(r.Body)
	if err != nil {
		serverError(w, "parsing event: %v", err)
		return
	}

	pretty.Print(ev)
	logInfo("Headers:\n%v\n", r.Header)
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
	logWarn("Client error: "+fmt, args...)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Client error"))
}
