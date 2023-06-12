package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler{w: w, r: r, config: config}.next()
	})
	logInfo("Listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	Port              string        `env:"PORT,default=8080"`
	Project           string        `env:"PROJECT,required"`
	Location          string        `env:"LOCATION,required"`
	HookID            string        `env:"HOOK_ID"` // Will validate against GitHub header, if provided.
	RunnerImageURL    string        `env:"RUNNER_IMAGE_URL,required"`
	JobID             string        `env:"JOB_ID,default=gha-runner"`
	JobTimeout        time.Duration `env:"JOB_TIMEOUT,default=10m"`
	JobCpu            string        `env:"JOB_CPU,default=1"`
	JobMemory         string        `env:"JOB_MEMORY,default=1Gi"`
	TokenSecretName   string        `env:"GITHUB_TOKEN_SECRET,required"` // "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.
	RepositoryHtmlURL string        `env:"REPOSITORY_URL,required"`
}

func newConfig(ctx context.Context) (config, error) {
	fmt.Printf("ENV VARS:\n%q\n", os.Environ())
	var c config
	if err := envconfig.Process(ctx, &c); err != nil {
		return config{}, fmt.Errorf("processing envconfig: %v", err)
	}
	c.JobID += "k" // TODO: remove me after testing.
	logInfo("Config: %#v", c)
	return c, nil
}
