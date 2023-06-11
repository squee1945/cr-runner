package main

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	actionCompleted  = "completed"
	actionInProgress = "in_progress"
	actionQueued     = "queued"
	actionWaiting    = "waiting"

	jobStatusQueued     = "queued"
	jobStatusInProgress = "in_progress"
	jobStatusCompleted  = "completed"
	jobStatusWaiting    = "waiting"

	stepConclusionFailure   = "failure"
	stepConclusionSkipped   = "skipped"
	stepConclusionSuccess   = "success"
	stepConclusionCancelled = "cancelled"
	stepConclusionNull      = "null"

	stepStatusCompleted  = "completed"
	stepStatusInProgress = "in_progress"
	stepStatusQueued     = "queued"
	stepStatusPending    = "pending"
	stepStatusWaiting    = "waiting"
)

func parseEvent(r io.Reader) (*event, error) {
	var e event
	if err := json.NewDecoder(r).Decode(&e); err != nil {
		return nil, fmt.Errorf("unmarshalling json: %v", err)
	}
	return &e, nil
}

type event struct {
	Action       string            `json:"action"`
	Sender       gitHubUser        `json:"sender"`
	Repository   eventRepository   `json:"repository"`
	Organization eventOrganization `json:"organization"`
	WorkflowJob  eventWorkflowJob  `json:"workflow_job"`
	Deployment   eventDeployment   `json:"deployment"`
	// Installation eventInstallation `json:"installation"`
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
type eventRepository struct {
	ID       int    `json:"id"`        // "id": 1296269,
	NodeID   string `json:"node_id"`   // "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
	Name     string `json:"name"`      // "name": "Hello-World",
	FullName string `json:"full_name"` // "full_name": "octocat/Hello-World",
}

// https://docs.github.com/en/rest/orgs/orgs?apiVersion=2022-11-28#get-an-organization
type eventOrganization struct {
	ID     int    `json:"id"`      // "id": 1,
	NodeID string `json:"node_id"` // "node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
	URL    string `json:"url"`     // "url": "https://api.github.com/orgs/github",
}

// https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads?actionType=waiting#workflow_job
type eventWorkflowJob struct {
	CheckRunURL     string              `json:"check_run_url"`
	CompletedAt     string              `json:"completed_at"` // "2011-01-26T19:06:43Z"
	Conclusion      string              `json:"conclusion"`
	CreatedAt       string              `json:"created_at"`
	HeadSHA         string              `json:"head_sha"`
	HtmlURL         string              `json:"html_url"`
	ID              int                 `json:"id"`
	Labels          []string            `json:"labels"`
	Name            string              `json:"name"`
	NodeID          string              `json:"node_id"`
	RunAttempt      int                 `json:"run_attempt"`
	RunID           int64               `json:"run_id"`
	RunURL          string              `json:"run_url"`
	RunnerGroupID   int                 `json:"runner_group_id"`
	RunnerGroupName string              `json:"runner_group_name"`
	RunnerID        int                 `json:"runner_id"`
	RunnerName      string              `json:"runner_name"`
	StartedAt       string              `json:"started_at"`
	HeadBranch      string              `json:"head_branch"`
	WorkflowName    string              `json:"workflow_name"`
	Status          string              `json:"status"` // One of: queued, in_progress, completed, waiting
	Steps           []eventWorkflowStep `json:"steps"`
	URL             string              `json:"url"`
}

// https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads?actionType=waiting#workflow_job
type eventWorkflowStep struct {
	CompletedAt string `json:"completed_at"`
	Conclusion  string `json:"conclusion"` // One of: failure, skipped, success, cancelled, null
	Name        string `json:"name"`
	Number      int    `json:"number"`
	StartedAt   string `json:"started_at"`
	Status      string `json:"status"` // One of: completed, in_progress, queued, pending, waiting
}

type eventDeployment struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
	SHA    string `json:"sha"`
	Ref    string `json:"ref"`
	Task   string `json:"task"`
	// Payload               eventDeploymentPayload `json:"payload"` // TODO may be a string
	OriginalEnvironment   string     `json:"original_environment"`
	Envioronment          string     `json:"environment"`
	Description           string     `json:"description"`
	Creator               gitHubUser `json:"creator"`
	CreatedAt             string     `json:"created_at"`
	UpdatedAt             string     `json:"updated_at"`
	StatusesURL           string     `json:"statuses_url"`
	RepositoryURL         string     `json:"repository_url"`
	TransientEnvironment  bool       `json:"transient_environment"`
	ProductionEnvironment bool       `json:"production_environment"`
	// PerformedViaGitHubApp any `json:"performed_via_github_app"`
}

// type eventDeploymentPayload struct {
// }

type gitHubUser struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HtmlURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredAt         string `json:"starred_at"`
}
