package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/kr/pretty"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

const (
	eventHeader  = "X-GitHub-Event"   // X-Github-Event: workflow_job
	hookIDHeader = "X-Github-Hook-Id" // X-Github-Hook-Id: 419040544
	sig256Header = "X-Hub-Signature-256"
	// deliveryHeader   = "X-GitHub-Delivery"                      // X-Github-Delivery: 2e4a3250-08a3-11ee-8cc8-00632da95790
	// targetIDHeader   = "X-Github-Hook-Installation-Target-Id"   // X-Github-Hook-Installation-Target-Id: 652279005
	// targetTypeHeader = "X-Github-Hook-Installation-Target-Type" // X-Github-Hook-Installation-Target-Type: repository
	// sigHeader        = "X-Hub-Signature"

	eventWorkFlowJob     = "workflow_job"
	eventAppInstallation = "installation"
)

var (
	validEvents = map[string]bool{
		eventWorkFlowJob:     true,
		eventAppInstallation: true,
	}
)

type handler struct {
	w      http.ResponseWriter
	r      *http.Request
	config config
}

func (h handler) next() {
	if h.r.Method != http.MethodPost {
		h.clientError("bad method %v", h.r.Method)
		return
	}
	if eh := h.r.Header.Get(hookIDHeader); h.config.HookID != "" && eh != h.config.HookID {
		h.clientError("incorrect %s got:%q want:%q", hookIDHeader, eh, h.config.HookID)
		return
	}

	eh := h.r.Header.Get(eventHeader)
	if !validEvents[eh] {
		h.clientError("unexpected event type %q", eh)
		return
	}

	body, err := io.ReadAll(h.r.Body)
	if err != nil {
		h.serverError("reading body: %v", err)
	}
	h.r.Body.Close()

	if err := h.validateSignature(body); err != nil {
		h.clientError("validating signature: %v", err)
		return
	}

	ev, err := parseEvent(body)
	if err != nil {
		h.serverError("parsing event: %v", err)
		return
	}

	switch eh {
	case eventWorkFlowJob:
		h.handleWorkFlowJob(ev)
	case eventAppInstallation:
		h.handleAppInstallation(ev)
	default:
		h.serverError("Unhandled event: %q", eh)
	}
}

func (h *handler) handleWorkFlowJob(ev *event) {
	if ev.Action != actionQueued {
		logInfo("Event action %q not %q. Ignoring.", ev.Action, actionQueued)
		return
	}

	logInfo("Processing event:\n%s\n", pretty.Sprint(ev))

	crJob := cloudRunJob{config: h.config}
	if err := crJob.runJob(h.r.Context()); err != nil {
		h.serverError("running job %q: %v", h.config.JobID, err)
		return
	}
}

func (h *handler) handleAppInstallation(ev *event) {
	if ev.Action != actionCreated {
		logInfo("Event action %q not %q. Ignoring.", ev.Action, actionCreated)
		return
	}

	logInfo("Received installation event: %#v", *ev)
}

func (h handler) validateSignature(body []byte) error {
	if h.config.SignatureSecretName == "" {
		// Signature validation not configured.
		return nil
	}

	// See https://docs.github.com/en/webhooks-and-events/webhooks/securing-your-webhooks
	messageMAC := h.r.Header.Get(sig256Header)
	if messageMAC == "" {
		return errors.New("$GITHUB_SIGNATURE_SECRET is set, but webhook message did not have signature. Did you configure the `Secret` in the GitHub webhook?")
	}

	signatureSecret, err := h.readSecret(h.config.SignatureSecretName)
	if err != nil {
		return fmt.Errorf("reading $GITHUB_SIGNATURE_SECRET secret: %v", err)
	}

	mac := hmac.New(sha256.New, signatureSecret)
	mac.Write(body)
	expectedMACBytes := mac.Sum(nil)

	expectedMAC := "sha256=" + hex.EncodeToString(expectedMACBytes)

	if messageMAC != expectedMAC {
		return errors.New("signatures do not match")
	}
	return nil
}

func (h handler) serverError(template string, args ...any) {
	logError("Error: "+template, args...)
	h.w.WriteHeader(http.StatusInternalServerError)
	h.w.Write([]byte("Server error"))
}

func (h handler) clientError(template string, args ...any) {
	logWarn("Client error: "+template, args...)
	h.w.WriteHeader(http.StatusBadRequest)
	h.w.Write([]byte("Client error"))
}

func (h handler) readSecret(name string) ([]byte, error) {
	ctx := h.r.Context()

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
