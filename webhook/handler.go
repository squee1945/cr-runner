package main

import (
	"net/http"

	"github.com/kr/pretty"
)

const (
	eventHeader      = "X-GitHub-Event"                         // X-Github-Event: workflow_job
	deliveryHeader   = "X-GitHub-Delivery"                      // X-Github-Delivery: 2e4a3250-08a3-11ee-8cc8-00632da95790
	hookIDHeader     = "X-Github-Hook-Id"                       // X-Github-Hook-Id: 419040544
	targetIDHeader   = "X-Github-Hook-Installation-Target-Id"   // X-Github-Hook-Installation-Target-Id: 652279005
	targetTypeHeader = "X-Github-Hook-Installation-Target-Type" // X-Github-Hook-Installation-Target-Type: repository
	sigHeader        = "X-Hub-Signature"
	sig256Header     = "X-Hub-Signature-256"

	eventWorkFlowJobHeader = "workflow_job"
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
	if eh := h.r.Header.Get(eventHeader); eh != eventWorkFlowJobHeader {
		h.clientError("unexpected event type %q", eh)
		return
	}
	if eh := h.r.Header.Get(hookIDHeader); h.config.HookID != "" && eh != h.config.HookID {
		h.clientError("incorrect %s got:%q want:%q", hookIDHeader, eh, h.config.HookID)
		return
	}
	// TODO: Check signatures.

	ev, err := parseEvent(h.r.Body)
	if err != nil {
		h.serverError("parsing event: %v", err)
		return
	}

	logInfo("Received event:\n%s\n", pretty.Sprint(ev))
	logInfo("Headers:\n%v\n", h.r.Header)

	if ev.Action != actionQueued {
		logInfo("Event action %q not %q. Ignoring.", ev.Action, actionQueued)
		return
	}

	crJob := cloudRunJob{config: h.config}
	if err := crJob.runJob(h.r.Context(), ev); err != nil {
		h.serverError("running job %q: %v", h.config.JobID, err)
		return
	}
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
