package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	// Required env vars.
	RepositoryURL   string `env:"REPOSITORY_URL,required"`
	RunnerImageURL  string `env:"RUNNER_IMAGE_URL,required"`
	TokenSecretName string `env:"GITHUB_TOKEN_SECRET,required"` // "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.

	// Optional env vars.
	HookID              string        `env:"HOOK_ID"`                 // Will validate against GitHub header, if provided.
	SignatureSecretName string        `env:"GITHUB_SIGNATURE_SECRET"` // Will validate against GitHub signatures, if provided. "{secret_name}" for same project, "projects/{project}/secrets/{secret_name}" for different project.
	JobID               string        `env:"JOB_ID,default=runner"`
	JobTimeout          time.Duration `env:"JOB_TIMEOUT,default=10m"`
	JobCpu              string        `env:"JOB_CPU,default=1"`
	JobMemory           string        `env:"JOB_MEMORY,default=1Gi"`
	Port                string        `env:"PORT,default=8080"`
	AppClientSecretName string        `env:"GITHUB_APP_CLIENT_SECRET"`

	// Pulled from metadata.
	Project  string
	Location string

	// JobVersion is used to ensure a unique hash for the JobID when the Cloud Run Job definition changes.
	JobVersion string
}

func newConfig(ctx context.Context) (config, error) {
	c := config{
		JobVersion: jobVersion,
	}
	if err := envconfig.Process(ctx, &c); err != nil {
		return config{}, fmt.Errorf("processing envconfig: %v", err)
	}

	var err error
	if c.Project, err = projectID(ctx); err != nil {
		return config{}, fmt.Errorf("fetching project ID from metadata server: %v", err)
	}
	if c.Location, err = location(ctx); err != nil {
		return config{}, fmt.Errorf("fetching location from metadata server: %v", err)
	}

	// Hash the config and use it as the suffix to the JobID.
	// The ensures a new job is created when the env var settings change.
	b, err := json.Marshal(c)
	if err != nil {
		return config{}, fmt.Errorf("marshalling config: %v", err)
	}
	h := md5.New()
	if _, err := io.WriteString(h, string(b)); err != nil {
		return config{}, fmt.Errorf("writing to hash: %v", err)
	}
	c.JobID += fmt.Sprintf("-%x", h.Sum(nil))

	logInfo("Config: %#v", c)
	return c, nil
}
