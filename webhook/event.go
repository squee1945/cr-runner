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
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading all: %v", err)
	}
	// logInfo("Raw event:\n%s\n", string(b))
	var e event
	if err := json.Unmarshal(b, &e); err != nil {
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
	// Deployment   eventDeployment   `json:"deployment"`
	// Installation eventInstallation `json:"installation"`
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
/*
  "repository": {
    "id": 652279005,
    "node_id": "R_kgDOJuD83Q",
    "name": "self-hosted-runner",
    "full_name": "squee1945/self-hosted-runner",
    "private": false,
    "owner": { [githubUser] },
    "html_url": "https://github.com/squee1945/self-hosted-runner",
    "description": null,
    "fork": false,
    "url": "https://api.github.com/repos/squee1945/self-hosted-runner",
    "forks_url": "https://api.github.com/repos/squee1945/self-hosted-runner/forks",
    "keys_url": "https://api.github.com/repos/squee1945/self-hosted-runner/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/squee1945/self-hosted-runner/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/squee1945/self-hosted-runner/teams",
    "hooks_url": "https://api.github.com/repos/squee1945/self-hosted-runner/hooks",
    "issue_events_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues/events{/number}",
    "events_url": "https://api.github.com/repos/squee1945/self-hosted-runner/events",
    "assignees_url": "https://api.github.com/repos/squee1945/self-hosted-runner/assignees{/user}",
    "branches_url": "https://api.github.com/repos/squee1945/self-hosted-runner/branches{/branch}",
    "tags_url": "https://api.github.com/repos/squee1945/self-hosted-runner/tags",
    "blobs_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/squee1945/self-hosted-runner/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/squee1945/self-hosted-runner/languages",
    "stargazers_url": "https://api.github.com/repos/squee1945/self-hosted-runner/stargazers",
    "contributors_url": "https://api.github.com/repos/squee1945/self-hosted-runner/contributors",
    "subscribers_url": "https://api.github.com/repos/squee1945/self-hosted-runner/subscribers",
    "subscription_url": "https://api.github.com/repos/squee1945/self-hosted-runner/subscription",
    "commits_url": "https://api.github.com/repos/squee1945/self-hosted-runner/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/squee1945/self-hosted-runner/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/squee1945/self-hosted-runner/contents/{+path}",
    "compare_url": "https://api.github.com/repos/squee1945/self-hosted-runner/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/squee1945/self-hosted-runner/merges",
    "archive_url": "https://api.github.com/repos/squee1945/self-hosted-runner/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/squee1945/self-hosted-runner/downloads",
    "issues_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues{/number}",
    "pulls_url": "https://api.github.com/repos/squee1945/self-hosted-runner/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/squee1945/self-hosted-runner/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/squee1945/self-hosted-runner/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/squee1945/self-hosted-runner/labels{/name}",
    "releases_url": "https://api.github.com/repos/squee1945/self-hosted-runner/releases{/id}",
    "deployments_url": "https://api.github.com/repos/squee1945/self-hosted-runner/deployments",
    "created_at": "2023-06-11T16:44:18Z",
    "updated_at": "2023-06-11T16:58:53Z",
    "pushed_at": "2023-06-11T22:07:56Z",
    "git_url": "git://github.com/squee1945/self-hosted-runner.git",
    "ssh_url": "git@github.com:squee1945/self-hosted-runner.git",
    "clone_url": "https://github.com/squee1945/self-hosted-runner.git",
    "svn_url": "https://github.com/squee1945/self-hosted-runner",
    "homepage": null,
    "size": 2,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": "Go",
    "has_issues": true,
    "has_projects": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "has_discussions": false,
    "forks_count": 0,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 0,
    "license": null,
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [],
    "visibility": "public",
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "main"
  },
*/
type eventRepository struct {
	ID       int    `json:"id"`        // "id": 1296269,
	NodeID   string `json:"node_id"`   // "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
	Name     string `json:"name"`      // "name": "Hello-World",
	FullName string `json:"full_name"` // "full_name": "octocat/Hello-World",
	HtmlURL  string `json:"html_url"`  // "html_url": "https://github.com/squee1945/self-hosted-runner"
}

// https://docs.github.com/en/rest/orgs/orgs?apiVersion=2022-11-28#get-an-organization
type eventOrganization struct {
	ID     int    `json:"id"`      // "id": 1,
	NodeID string `json:"node_id"` // "node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
	URL    string `json:"url"`     // "url": "https://api.github.com/orgs/github",
}

// https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads?actionType=waiting#workflow_job
/*
{
    "id": 14171575672,
    "run_id": 5237732760,
    "workflow_name": "Go",
    "head_branch": "main",
    "run_url": "https://api.github.com/repos/squee1945/self-hosted-runner/actions/runs/5237732760",
    "run_attempt": 1,
    "node_id": "CR_kwDOJuD83c8AAAADTLEVeA",
    "head_sha": "720bd4aae6521bfce13488b41b93919c64088420",
    "url": "https://api.github.com/repos/squee1945/self-hosted-runner/actions/jobs/14171575672",
    "html_url": "https://github.com/squee1945/self-hosted-runner/actions/runs/5237732760/jobs/9456107741",
    "status": "queued",
    "conclusion": null,
    "created_at": "2023-06-11T22:08:00Z",
    "started_at": "2023-06-11T22:07:59Z",
    "completed_at": null,
    "name": "build",
    "steps": [],
    "check_run_url": "https://api.github.com/repos/squee1945/self-hosted-runner/check-runs/14171575672",
    "labels": [
      "ubuntu-latest"
    ],
    "runner_id": null,
    "runner_name": null,
    "runner_group_id": null,
    "runner_group_name": null
  }*/
type eventWorkflowJob struct {
	ID              int                 `json:"id"`
	RunID           int64               `json:"run_id"`
	WorkflowName    string              `json:"workflow_name"`
	HeadBranch      string              `json:"head_branch"`
	RunURL          string              `json:"run_url"`
	RunAttempt      int                 `json:"run_attempt"`
	NodeID          string              `json:"node_id"`
	HeadSHA         string              `json:"head_sha"`
	URL             string              `json:"url"`
	HtmlURL         string              `json:"html_url"`
	Status          string              `json:"status"` // One of: queued, in_progress, completed, waiting
	Conclusion      string              `json:"conclusion"`
	CreatedAt       string              `json:"created_at"`
	StartedAt       string              `json:"started_at"`
	CompletedAt     string              `json:"completed_at"` // "2011-01-26T19:06:43Z"
	Name            string              `json:"name"`
	Steps           []eventWorkflowStep `json:"steps"`
	CheckRunURL     string              `json:"check_run_url"`
	Labels          []string            `json:"labels"`
	RunnerID        int                 `json:"runner_id"`
	RunnerName      string              `json:"runner_name"`
	RunnerGroupID   int                 `json:"runner_group_id"`
	RunnerGroupName string              `json:"runner_group_name"`
}

// https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads?actionType=waiting#workflow_job
/*
{
	{CompletedAt:"2023-06-11T22:34:45.000Z", Conclusion:"success", Name:"Set up job", Number:1, StartedAt:"2023-06-11T22:34:43.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:46.000Z", Conclusion:"success", Name:"Run actions/checkout@v3", Number:2, StartedAt:"2023-06-11T22:34:45.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:53.000Z", Conclusion:"success", Name:"Set up Go", Number:3, StartedAt:"2023-06-11T22:34:47.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:54.000Z", Conclusion:"success", Name:"Build", Number:4, StartedAt:"2023-06-11T22:34:53.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:55.000Z", Conclusion:"success", Name:"Test", Number:5, StartedAt:"2023-06-11T22:34:54.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:56.000Z", Conclusion:"success", Name:"Post Set up Go", Number:9, StartedAt:"2023-06-11T22:34:56.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:56.000Z", Conclusion:"success", Name:"Post Run actions/checkout@v3", Number:10, StartedAt:"2023-06-11T22:34:56.000Z", Status:"completed"},
	{CompletedAt:"2023-06-11T22:34:56.000Z", Conclusion:"success", Name:"Complete job", Number:11, StartedAt:"2023-06-11T22:34:55.000Z", Status:"completed"},
}
*/
type eventWorkflowStep struct {
	CompletedAt string `json:"completed_at"`
	Conclusion  string `json:"conclusion"` // One of: failure, skipped, success, cancelled, null
	Name        string `json:"name"`
	Number      int    `json:"number"`
	StartedAt   string `json:"started_at"`
	Status      string `json:"status"` // One of: completed, in_progress, queued, pending, waiting
}

// type eventDeployment struct {
// 	URL    string `json:"url"`
// 	ID     int    `json:"id"`
// 	NodeID string `json:"node_id"`
// 	SHA    string `json:"sha"`
// 	Ref    string `json:"ref"`
// 	Task   string `json:"task"`
// 	// Payload               eventDeploymentPayload `json:"payload"` // TODO may be a string
// 	OriginalEnvironment   string     `json:"original_environment"`
// 	Envioronment          string     `json:"environment"`
// 	Description           string     `json:"description"`
// 	Creator               gitHubUser `json:"creator"`
// 	CreatedAt             string     `json:"created_at"`
// 	UpdatedAt             string     `json:"updated_at"`
// 	StatusesURL           string     `json:"statuses_url"`
// 	RepositoryURL         string     `json:"repository_url"`
// 	TransientEnvironment  bool       `json:"transient_environment"`
// 	ProductionEnvironment bool       `json:"production_environment"`
// 	// PerformedViaGitHubApp any `json:"performed_via_github_app"`
// }

// type eventDeploymentPayload struct {
// }

/*
"login": "squee1945",
"id": 1146523,
"node_id": "MDQ6VXNlcjExNDY1MjM=",
"avatar_url": "https://avatars.githubusercontent.com/u/1146523?v=4",
"gravatar_id": "",
"url": "https://api.github.com/users/squee1945",
"html_url": "https://github.com/squee1945",
"followers_url": "https://api.github.com/users/squee1945/followers",
"following_url": "https://api.github.com/users/squee1945/following{/other_user}",
"gists_url": "https://api.github.com/users/squee1945/gists{/gist_id}",
"starred_url": "https://api.github.com/users/squee1945/starred{/owner}{/repo}",
"subscriptions_url": "https://api.github.com/users/squee1945/subscriptions",
"organizations_url": "https://api.github.com/users/squee1945/orgs",
"repos_url": "https://api.github.com/users/squee1945/repos",
"events_url": "https://api.github.com/users/squee1945/events{/privacy}",
"received_events_url": "https://api.github.com/users/squee1945/received_events",
"type": "User",
"site_admin": false
*/
type gitHubUser struct {
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

	// The following items do not seem to be emitted by GitHub.
	Name      string `json:"name"`
	Email     string `json:"email"`
	StarredAt string `json:"starred_at"`
}

/*
{
  "action": "queued",
  "workflow_job": {
    "id": 14171575672,
    "run_id": 5237732760,
    "workflow_name": "Go",
    "head_branch": "main",
    "run_url": "https://api.github.com/repos/squee1945/self-hosted-runner/actions/runs/5237732760",
    "run_attempt": 1,
    "node_id": "CR_kwDOJuD83c8AAAADTLEVeA",
    "head_sha": "720bd4aae6521bfce13488b41b93919c64088420",
    "url": "https://api.github.com/repos/squee1945/self-hosted-runner/actions/jobs/14171575672",
    "html_url": "https://github.com/squee1945/self-hosted-runner/actions/runs/5237732760/jobs/9456107741",
    "status": "queued",
    "conclusion": null,
    "created_at": "2023-06-11T22:08:00Z",
    "started_at": "2023-06-11T22:07:59Z",
    "completed_at": null,
    "name": "build",
    "steps": [],
    "check_run_url": "https://api.github.com/repos/squee1945/self-hosted-runner/check-runs/14171575672",
    "labels": [
      "ubuntu-latest"
    ],
    "runner_id": null,
    "runner_name": null,
    "runner_group_id": null,
    "runner_group_name": null
  },
  "repository": {
    "id": 652279005,
    "node_id": "R_kgDOJuD83Q",
    "name": "self-hosted-runner",
    "full_name": "squee1945/self-hosted-runner",
    "private": false,
    "owner": {
      "login": "squee1945",
      "id": 1146523,
      "node_id": "MDQ6VXNlcjExNDY1MjM=",
      "avatar_url": "https://avatars.githubusercontent.com/u/1146523?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/squee1945",
      "html_url": "https://github.com/squee1945",
      "followers_url": "https://api.github.com/users/squee1945/followers",
      "following_url": "https://api.github.com/users/squee1945/following{/other_user}",
      "gists_url": "https://api.github.com/users/squee1945/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/squee1945/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/squee1945/subscriptions",
      "organizations_url": "https://api.github.com/users/squee1945/orgs",
      "repos_url": "https://api.github.com/users/squee1945/repos",
      "events_url": "https://api.github.com/users/squee1945/events{/privacy}",
      "received_events_url": "https://api.github.com/users/squee1945/received_events",
      "type": "User",
      "site_admin": false
    },
    "html_url": "https://github.com/squee1945/self-hosted-runner",
    "description": null,
    "fork": false,
    "url": "https://api.github.com/repos/squee1945/self-hosted-runner",
    "forks_url": "https://api.github.com/repos/squee1945/self-hosted-runner/forks",
    "keys_url": "https://api.github.com/repos/squee1945/self-hosted-runner/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/squee1945/self-hosted-runner/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/squee1945/self-hosted-runner/teams",
    "hooks_url": "https://api.github.com/repos/squee1945/self-hosted-runner/hooks",
    "issue_events_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues/events{/number}",
    "events_url": "https://api.github.com/repos/squee1945/self-hosted-runner/events",
    "assignees_url": "https://api.github.com/repos/squee1945/self-hosted-runner/assignees{/user}",
    "branches_url": "https://api.github.com/repos/squee1945/self-hosted-runner/branches{/branch}",
    "tags_url": "https://api.github.com/repos/squee1945/self-hosted-runner/tags",
    "blobs_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/squee1945/self-hosted-runner/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/squee1945/self-hosted-runner/languages",
    "stargazers_url": "https://api.github.com/repos/squee1945/self-hosted-runner/stargazers",
    "contributors_url": "https://api.github.com/repos/squee1945/self-hosted-runner/contributors",
    "subscribers_url": "https://api.github.com/repos/squee1945/self-hosted-runner/subscribers",
    "subscription_url": "https://api.github.com/repos/squee1945/self-hosted-runner/subscription",
    "commits_url": "https://api.github.com/repos/squee1945/self-hosted-runner/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/squee1945/self-hosted-runner/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/squee1945/self-hosted-runner/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/squee1945/self-hosted-runner/contents/{+path}",
    "compare_url": "https://api.github.com/repos/squee1945/self-hosted-runner/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/squee1945/self-hosted-runner/merges",
    "archive_url": "https://api.github.com/repos/squee1945/self-hosted-runner/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/squee1945/self-hosted-runner/downloads",
    "issues_url": "https://api.github.com/repos/squee1945/self-hosted-runner/issues{/number}",
    "pulls_url": "https://api.github.com/repos/squee1945/self-hosted-runner/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/squee1945/self-hosted-runner/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/squee1945/self-hosted-runner/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/squee1945/self-hosted-runner/labels{/name}",
    "releases_url": "https://api.github.com/repos/squee1945/self-hosted-runner/releases{/id}",
    "deployments_url": "https://api.github.com/repos/squee1945/self-hosted-runner/deployments",
    "created_at": "2023-06-11T16:44:18Z",
    "updated_at": "2023-06-11T16:58:53Z",
    "pushed_at": "2023-06-11T22:07:56Z",
    "git_url": "git://github.com/squee1945/self-hosted-runner.git",
    "ssh_url": "git@github.com:squee1945/self-hosted-runner.git",
    "clone_url": "https://github.com/squee1945/self-hosted-runner.git",
    "svn_url": "https://github.com/squee1945/self-hosted-runner",
    "homepage": null,
    "size": 2,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": "Go",
    "has_issues": true,
    "has_projects": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "has_discussions": false,
    "forks_count": 0,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 0,
    "license": null,
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [],
    "visibility": "public",
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "main"
  },
  "sender": {
    "login": "squee1945",
    "id": 1146523,
    "node_id": "MDQ6VXNlcjExNDY1MjM=",
    "avatar_url": "https://avatars.githubusercontent.com/u/1146523?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/squee1945",
    "html_url": "https://github.com/squee1945",
    "followers_url": "https://api.github.com/users/squee1945/followers",
    "following_url": "https://api.github.com/users/squee1945/following{/other_user}",
    "gists_url": "https://api.github.com/users/squee1945/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/squee1945/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/squee1945/subscriptions",
    "organizations_url": "https://api.github.com/users/squee1945/orgs",
    "repos_url": "https://api.github.com/users/squee1945/repos",
    "events_url": "https://api.github.com/users/squee1945/events{/privacy}",
    "received_events_url": "https://api.github.com/users/squee1945/received_events",
    "type": "User",
    "site_admin": false
  }
}
*/
