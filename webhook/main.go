package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultRunnerImageURL = "us-central1-docker.pkg.dev/cr-runner-jasonco/github-actions-runner/image@sha256:8c87e13c36ca3d2d3703bde6a06979bf2daba47b963edbff281cfd4cd468375b"
	defaultJobTimeout     = 60 * time.Minute
	defaultJobCpu         = "1"
	defaultJobMemory      = "512Mi"
	defaultSecretName     = "github-actions-runner"

	hookIDEnvVar            = "HOOK_ID"
	runnerImageURLEnvVar    = "RUNNER_IMAGE_URL"
	jobTimeoutEnvVar        = "JOB_TIMEOUT"
	jobCpuEnvVar            = "JOB_CPU"
	jobMemoryEnvVar         = "JOB_MEMORY"
	gitHubTokenSecretEnvVar = "GITHUB_TOKEN_SECRET"
)

func main() {
	log.SetFlags(0)

	logInfo("Starting server...")

	// Determine port for HTTP service.
	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	config, err := newConfig()
	if err != nil {
		log.Fatalf("Bad config: %v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler{w: w, r: r, config: config}.next()
	})

	// Start HTTP server.
	logInfo("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	project         string
	location        string
	wantHookID      string
	runnerImageURL  string
	jobTimeout      time.Duration
	jobCpu          string
	jobMemory       string
	tokenSecretName string // "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.
}

func newConfig() (config, error) {
	fmt.Printf("ENV VARS:\n%q\n", os.Environ())
	c := config{
		wantHookID:      os.Getenv(hookIDEnvVar),
		runnerImageURL:  defaultRunnerImageURL,
		project:         "cr-runner-jasonco", // TODO
		location:        "us-central1",       // TODO
		jobTimeout:      defaultJobTimeout,
		jobCpu:          defaultJobCpu,    // TODO
		jobMemory:       defaultJobMemory, // TODO
		tokenSecretName: defaultSecretName,
	}
	if sn, ok := os.LookupEnv(gitHubTokenSecretEnvVar); ok {
		c.tokenSecretName = sn
	}
	if ts, ok := os.LookupEnv(jobTimeoutEnvVar); ok {
		var err error
		c.jobTimeout, err = time.ParseDuration(ts)
		if err != nil {
			return config{}, fmt.Errorf("parsing %s=%q: %v", jobTimeoutEnvVar, ts, err)
		}
	}
	if url, ok := os.LookupEnv(runnerImageURLEnvVar); ok {
		c.runnerImageURL = url
	}
	logInfo("Config: %#v", c)
	return c, nil
}
