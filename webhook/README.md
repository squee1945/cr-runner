# Example to use Cloud Run Jobs as ephemeral GitHub Actions Runners

This webhook application can be deployed as a regular Cloud Run service.
The URL for this service (e.g., https://my-service-wa5vxuhyra-uc.a.run.app)
can we used as the webhook to receive Actions events. When a new Action is
required, this webhook will be invoked, which will cause a Cloud Run Job
to start. The Cloud Run Job uses a GitHub runner image which connects
to GitHub on a long-polling HTTP connection. GitHub then sends the job to
the runner. The runner is configured as --ephemeral meaning it will process
only one job, then exit.

## Environment variables for the Cloud Run service

The following env vars must be configured on the Cloud Run service:
- `$PROJECT` - The project for the Cloud Run service (e.g., "my-project").
- `$LOCATION` - The location for the Cloud Run service (e.g., "us-central1").
- `$HOOK_ID` (optional) - The Hook ID for the webhook POSTing to the Cloud Run Service (e.g., "123456"); will only be validated if provided.
- `$RUNNER_IMAGE_URL` - The Artifact Registry URL for the runner image (see below).
- `$JOB_ID` (default "runner") - The name of the Cloud Run job. If you change the definition of the Job in the code, you must update this value to something unique.
- `$JOB_TIMEOUT` (default "10m") - The allowed time for the action to execute.
- `$JOB_CPU` (default "1") - The CPUs allocated for the job. See https://cloud.google.com/run/docs/configuring/cpu
- `$JOB_MEMORY` (default "512Mi") - The RAM allocated for the job. See https://cloud.google.com/run/docs/configuring/memory-limits
- `$GITHUB_TOKEN_SECRET` - The name of a Secret Manager secret holding your GitHub Personal Access Token (see below). **DO NOT PUT THE TOKEN ITSELF IN THIS ENV VAR!**
- `$REPOSITORY_URL` - The GitHub URL for the repository running actions (e.g., "https://github.com/joeschmoe/my-repo")

