package main

import (
	"fmt"
	"net/http"

	"github.com/kr/pretty"

	run "cloud.google.com/go/run/apiv2"
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

// type server struct {}

// func (s server) handler(config config) func (http.ResponseWriter,  *http.Request){
// 	return func (w http.ResponseWriter, r *http.Request) {
// 		handler{w: w, r: r, config: config}.next()
// 	}
// }

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
	if eh := h.r.Header.Get(hookIDHeader); h.config.wantHookID != "" && eh != h.config.wantHookID {
		h.clientError("incorrect %s got:%q want:%q", hookIDHeader, eh, h.config.wantHookID)
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

	// Start Cloud Run Job.
	// if err := h.startCRJob(ev); err != nil {
	// 	h.serverError("creating Cloud Run job: %v", err)
	// 	return
	// }
}

func (h handler) startCRJob(ev *event) error {
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	ctx := h.r.Context()
	c, err := run.NewJobsClient(ctx)
	if err != nil {
		return fmt.Errorf("creating Cloud Run client: %v", err)
	}
	defer c.Close()

	crJob := cloudRunJob{ev: ev, config: h.config}
	req, err := crJob.createJobRequest()
	if err != nil {
		return fmt.Errorf("creating job request: %v", err)
	}

	op, err := c.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("creating Cloud Run job: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for Cloud Run operation: %v", err)
	}

	logInfo("Cloud Run Job API response for %q: %#v", crJob.jobID(), resp)
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
