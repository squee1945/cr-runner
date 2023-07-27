package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

	repo := h.r.URL.Query().Get("repo")
	if repo == "" {
		h.clientError(`missing repo (e.g., "owner/repo")`)
		return
	}
	// TODO: validate repo

	installationID, err := h.parseIntParam("installation-id")
	if err != nil {
		h.clientError(err.Error())
		return
	}
	applicationID, err := h.parseIntParam("application-id")
	if err != nil {
		h.clientError(err.Error())
		return
	}

	pk, err := h.privateKey()
	if err != nil {
		h.serverError("fetching private key: %v", err)
		return
	}

	r := NewRegistration(applicationID, installationID, repo, pk)
	token, err := r.Token(h.r.Context())
	if err != nil {
		h.serverError("generating registration token: %v", err)
		return
	}

	logInfo("Token: " + string(token))
	h.w.Write([]byte("See logs"))
}

func (h apphandler) privateKey() (*rsa.PrivateKey, error) {
	if h.config.AppPrivateKeyName == "" {
		return nil, errors.New("missing GitHub app private key, did you set $GITHUB_APP_PRIVATE_KEY https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/managing-private-keys-for-github-apps#generating-private-keys")
	}
	raw, err := readSecret(h.r.Context(), h.config, h.config.AppPrivateKeyName)
	if err != nil {
		return nil, fmt.Errorf("reading private key pem: %v", err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("no block found in pem")
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %v", err)
	}
	return pk, nil
}

func (h apphandler) parseIntParam(name string) (int64, error) {
	s := h.r.URL.Query().Get(name)
	if s == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	i, err := strconv.ParseInt(s, 10)
	if err != nil {
		return 0, fmt.Errorf("%s must be integer: %v", name, err)
	}
	if i <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}
	return i, nil
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
