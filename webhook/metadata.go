package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func projectID(ctx context.Context) (string, error) {
	return metadataQuery(ctx, "/project/project-id")
}

func location(ctx context.Context) (string, error) {
	full, err := metadataQuery(ctx, "/instance/region") // full is like "projects/659154930685/regions/us-central1"
	if err != nil {
		return "", nil
	}
	parts := strings.Split(full, "/")
	return parts[len(parts)-1], nil
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
