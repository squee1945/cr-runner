package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	HookID            string        `env:"HOOK_ID"` // Will validate against GitHub header, if provided.
	RunnerImageURL    string        `env:"RUNNER_IMAGE_URL,required"`
	JobID             string        `env:"JOB_ID,default=runner"`
	JobTimeout        time.Duration `env:"JOB_TIMEOUT,default=10m"`
	JobCpu            string        `env:"JOB_CPU,default=1"`
	JobMemory         string        `env:"JOB_MEMORY,default=1Gi"`
	TokenSecretName   string        `env:"GITHUB_TOKEN_SECRET,required"` // "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.
	RepositoryHtmlURL string        `env:"REPOSITORY_URL,required"`

	Project  string
	Location string
}

func newConfig(ctx context.Context) (config, error) {
	var c config
	if err := envconfig.Process(ctx, &c); err != nil {
		return config{}, fmt.Errorf("processing envconfig: %v", err)
	}

	project, err := projectID(ctx)
	if err != nil {
		return config{}, fmt.Errorf("fetching project ID from metadata server: %v", err)
	}
	c.Project = project
	loc, err := location(ctx)
	if err != nil {
		return config{}, fmt.Errorf("fetching location from metadata server: %v", err)
	}
	c.Location = loc

	// Hash the config and use it as the suffix to the JobID.
	// The ensures a new job is created when the env var settings change.
	b, err := json.Marshal(c)
	if err != nil {
		return config{}, fmt.Errorf("marshalling config: %v", err)
	}
	h := sha1.New()
	if _, err := io.WriteString(h, string(b)); err != nil {
		return config{}, fmt.Errorf("writing to hash: %v", err)
	}
	c.JobID += fmt.Sprintf("-%x", h.Sum(nil))[:11]

	logInfo("Config: %#v", c)
	return c, nil
}

func projectID(ctx context.Context) (string, error) {
	return metadataQuery(ctx, "/project/project-id")
}

func location(ctx context.Context) (string, error) {
	return metadataQuery(ctx, "/instance/region")
}

func metadataQuery(ctx context.Context, path string) (string, error) {
	url := "http://metadata.google.internal/computeMetadata/v1" + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating metadata request: %v", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching metadata: %v", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading metadata response: %v", err)
	}

	return strings.TrimSpace(string(b)), nil
}
