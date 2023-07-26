package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtTimeout = 5 * time.Minute
)

type apphandler struct {
	w      http.ResponseWriter
	r      *http.Request
	config config
}

func (h apphandler) next() {
	if h.r.Method != http.MethodGet {
		h.clientError("bad method %v", h.r.Method)
		return
	}

	installationID := h.r.URL.Query().Get("installation-id")
	if installationID == "" {
		h.clientError("missing installation-id")
		return
	}

	applicationID := h.r.URL.Query().Get("application-id")
	if applicationID == "" {
		h.clientError("missing application-id")
		return
	}

	jwt, err := h.generateJWT(applicationID)
	if err != nil {
		h.serverError("generating JWT: %v", err)
		return
	}
	logInfo("JWT: %q", jwt)

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		h.serverError("creating http request: %v", err)
		return
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		h.serverError("calling access_tokens API (status %d): %v", res.StatusCode, err)
		return
	}

	token, err := io.ReadAll(res.Body)
	if err != nil {
		h.serverError("reading access_tokens response body: %v", err)
		return
	}

	h.w.Write(token)
}

func (h apphandler) generateJWT(applicationID string) (string, error) {
	if h.config.AppPrivateKeyName == "" {
		return "", errors.New("missing GitHub app private key, did you set $GITHUB_APP_PRIVATE_KEY https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/managing-private-keys-for-github-apps#generating-private-keys")
	}

	raw, err := readSecret(h.r.Context(), h.config, h.config.AppPrivateKeyName)
	if err != nil {
		return "", fmt.Errorf("reading private key pem: %v", err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return "", errors.New("no block found in pem")
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parsing private key: %v", err)
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Add(-1 * time.Minute).Unix(), // Guard against clock drift.
		"exp": now.Add(jwtTimeout).Unix(),
		"iss": applicationID,
		"alg": "RS256",
	})

	return token.SignedString(pk)
}

func (h apphandler) serverError(template string, args ...any) {
	logError("Error: "+template, args...)
	h.w.WriteHeader(http.StatusInternalServerError)
	h.w.Write([]byte("Server error"))
}

func (h apphandler) clientError(template string, args ...any) {
	msg := "Client error: " + fmt.Sprintf(template, args...)
	logWarn(msg)
	h.w.WriteHeader(http.StatusBadRequest)
	h.w.Write([]byte(msg))
}
