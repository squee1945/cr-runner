package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Registration struct {
	applicationID  int64
	installationID int64
	repo           string
	pk             *rsa.PrivateKey
}

func NewRegistration(applicationID, installationID int64, repo string, pk *rsa.PrivateKey) Registration {
	return Registration{
		applicationID:  applicationID,
		installationID: installationID,
		repo:           repo,
		pk:             pk,
	}
}

func (r Registration) Token(ctx context.Context) (string, error) {
	jwt, err := r.generateJWT(ctx)
	if err != nil {
		return "", fmt.Errorf("generating JWT: %v", err)
	}
	appToken, err := r.appAccessToken(ctx, jwt)
	if err != nil {
		return "", fmt.Errorf("generating app access token: %v", err)
	}
	regToken, err := r.appRegistrationToken(ctx, appToken)
	if err != nil {
		return "", fmt.Errorf("generating app registration token: %v", err)
	}
	return regToken, nil
}

func (r Registration) generateJWT(ctx context.Context) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Add(-1 * time.Minute).Unix(), // Guard against clock drift.
		"exp": now.Add(jwtTimeout).Unix(),
		"iss": r.applicationID,
		"alg": "RS256",
	})
	return token.SignedString(r.pk)
}

func (r Registration) appAccessToken(ctx context.Context, jwt string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", r.installationID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating http request: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling access_tokens API (status %d): %v", res.StatusCode, err)
	}

	tokenBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading access_tokens response body: %v", err)
	}

	type ghToken struct {
		Token              string            `json:"token"`
		ExpiresAt          time.Time         `json:"expires_at"`
		Permissions        map[string]string `json:"permissions"`
		RepositorySelected string            `json:"repository_selection"`
	}
	var ght ghToken
	if err := json.Unmarshal(tokenBytes, &ght); err != nil {
		return "", fmt.Errorf("unmarshalling response: %v", err)
	}
	if ght.Token == "" {
		return "", errors.New("token was empty")
	}
	return ght.Token, nil
}

func (r Registration) appRegistrationToken(ctx context.Context, appAccessToken string) (string, error) {
	// https://docs.github.com/en/rest/actions/self-hosted-runners?apiVersion=2022-11-28#create-a-registration-token-for-a-repository
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runners/registration-token", r.repo)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating http request: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+appAccessToken)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling registration-token API (status %d): %v", res.StatusCode, err)
	}

	tokenBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading registration-token response body: %v", err)
	}

	type ghToken struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	var ght ghToken
	if err := json.Unmarshal(tokenBytes, &ght); err != nil {
		return "", fmt.Errorf("unmarshalling response: %v", err)
	}
	if ght.Token == "" {
		return "", errors.New("token was empty")
	}
	return ght.Token, nil
}
