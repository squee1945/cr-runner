package main

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func readSecret(ctx context.Context, name string) ([]byte, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating client: %v", err)
	}
	defer client.Close()

	if !strings.HasPrefix(name, "projects/") {
		name = fmt.Sprintf("projects/%s/secrets/%s", h.config.Project, name)
	}
	parts := strings.Split(name, "/")
	if len(parts) < 6 {
		name += "/versions/latest"
	}
	logInfo("Accessing secret %q", name)

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{Name: name}
	response, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("accessing secret: %v", err)
	}
	return response.Payload.Data, nil

}
