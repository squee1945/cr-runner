package main

import (
	"net/http"
)

type apphandler struct {
	w      http.ResponseWriter
	r      *http.Request
	config config
}

func (h apphandler) next() {

	// https://.../app/callback?code=8610cb86e134c915eeb0&installation_id=40040599&setup_action=install

	if h.r.Method != http.MethodGet {
		h.clientError("bad method %v", h.r.Method)
		return
	}

	code := h.r.URL.Query().Get("code")
	if code == "" {
		h.clientError("missing code")
		return
	}
	installID := h.r.URL.Query().Get("installation_id")
	if installID == "" {
		h.clientError("missing installation_id")
		return
	}
	setupAction := h.r.URL.Query().Get("setup_action")
	if setupAction != "install" {
		h.clientError("unexpected setup_action %q", setupAction)
		return
	}

	logInfo("Request headers: %v", h.r.Header)

	// body, err := io.ReadAll(h.r.Body)
	// if err != nil {
	// 	h.serverError("reading body: %v", err)
	// }
	// h.r.Body.Close()

	// TODO: Validate signature

	logInfo("POST body: %s", string(body))
}

func (h apphandler) serverError(template string, args ...any) {
	logError("Error: "+template, args...)
	h.w.WriteHeader(http.StatusInternalServerError)
	h.w.Write([]byte("Server error"))
}

func (h apphandler) clientError(template string, args ...any) {
	logWarn("Client error: "+template, args...)
	h.w.WriteHeader(http.StatusBadRequest)
	h.w.Write([]byte("Client error"))
}
