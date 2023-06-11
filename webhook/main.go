package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kr/pretty"
)

const (
	eventHeader      = "X-GitHub-Event"                         // X-Github-Event: workflow_job
	deliveryHeader   = "X-GitHub-Delivery"                      // X-Github-Delivery: 2e4a3250-08a3-11ee-8cc8-00632da95790
	hookIDHeader     = "X-Github-Hook-Id"                       // X-Github-Hook-Id: 419040544
	targetIDHeader   = "X-Github-Hook-Installation-Target-Id"   // X-Github-Hook-Installation-Target-Id: 652279005
	targetTypeHeader = "X-Github-Hook-Installation-Target-Type" // X-Github-Hook-Installation-Target-Type: repository
	sigHeader        = "X-Hub-Signature"
	sig256Header     = "X-Hub-Signature-256"
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

	logInfo("Received event:\n%s\n", pretty.Sprint(ev))
	logInfo("Headers:\n%v\n", r.Header)
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

func logInfo(template string, args ...any) {
	logStructured("INFO", template, args...)
}

func logWarn(template string, args ...any) {
	logStructured("WARN", template, args...)
}

func logError(template string, args ...any) {
	logStructured("ERROR", template, args...)
}

type structuredLog struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func logStructured(severity string, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	sl := structuredLog{
		Severity: severity,
		Message:  msg,
	}
	content, err := json.Marshal(sl)
	if err != nil {
		fmt.Printf("Failed to log (message below): %v\n", err)
		fmt.Println(msg)
		return
	}
	fmt.Println(string(content))
}
