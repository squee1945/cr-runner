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

const (

// defaultRunnerImageURL = "us-central1-docker.pkg.dev/cr-runner-jasonco/github-actions-runner/image@sha256:8c87e13c36ca3d2d3703bde6a06979bf2daba47b963edbff281cfd4cd468375b"
// defaultJobTimeout     = 10 * time.Second
// defaultJobCpu         = "1"
// defaultJobMemory      = "512Mi"
// defaultSecretName     = "github-actions-runner"

// portEnvVar              = "PORT"
// hookIDEnvVar            = "HOOK_ID"
// runnerImageURLEnvVar    = "RUNNER_IMAGE_URL"
// jobTimeoutEnvVar        = "JOB_TIMEOUT"
// jobCpuEnvVar            = "JOB_CPU"
// jobMemoryEnvVar         = "JOB_MEMORY"
// gitHubTokenSecretEnvVar = "GITHUB_TOKEN_SECRET"
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
	JobID             string        `env:"JOB_ID,required"`
	JobTimeout        time.Duration `env:"JOB_TIMEOUT,default=30s"`
	JobCpu            string        `env:"JOB_CPU,default=1"`
	JobMemory         string        `env:"JOB_MEMORY,default=512Mi"`
	TokenSecretName   string        `env:"GITHUB_TOKEN_SECRET,required"` // "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.
	RepositoryHtmlURL string        `env:"REPOSITORY_URL,required"`
}

func newConfig(ctx context.Context) (config, error) {
	fmt.Printf("ENV VARS:\n%q\n", os.Environ())
	var c config
	if err := envconfig.Process(ctx, &c); err != nil {
		return config{}, fmt.Errorf("processing envconfig: %v", err)
	}
	c.JobID += "i"
	// if c.RepositoryHtmlURL == "" {
	// 	c.RepositoryHtmlURL = "https://github.com/squee1945/self-hosted-runner" // TODO
	// }
	// c.JobID = "github-runner-" + jobVersion // TODO
	// c := config{
	// 	port:              "8080",
	// 	wantHookID:        os.Getenv(hookIDEnvVar),
	// 	jobID:             "github-runner-" + jobVersion,
	// 	runnerImageURL:    defaultRunnerImageURL,
	// 	project:           "cr-runner-jasonco", // TODO
	// 	location:          "us-central1",       // TODO
	// 	jobTimeout:        defaultJobTimeout,
	// 	jobCpu:            defaultJobCpu,    // TODO
	// 	jobMemory:         defaultJobMemory, // TODO
	// 	tokenSecretName:   defaultSecretName,
	// 	repositoryHtmlURL: "https://github.com/squee1945/self-hosted-runner", // TODO
	// }
	// if p, ok := os.LookupEnv(portEnvVar); ok {
	// 	c.port = p
	// }
	// if sn, ok := os.LookupEnv(gitHubTokenSecretEnvVar); ok {
	// 	c.tokenSecretName = sn
	// }
	// if ts, ok := os.LookupEnv(jobTimeoutEnvVar); ok {
	// 	var err error
	// 	c.jobTimeout, err = time.ParseDuration(ts)
	// 	if err != nil {
	// 		return config{}, fmt.Errorf("parsing %s=%q: %v", jobTimeoutEnvVar, ts, err)
	// 	}
	// }
	// if url, ok := os.LookupEnv(runnerImageURLEnvVar); ok {
	// 	c.runnerImageURL = url
	// }
	logInfo("Config: %#v", c)
	return c, nil
}
