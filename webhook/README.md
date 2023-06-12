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
- `$RUNNER_IMAGE_URL` - The Artifact Registry URL for the runner image (see below).
- `$GITHUB_TOKEN_SECRET` - The name of a Secret Manager secret holding your GitHub Personal Access Token (see below). **DO NOT PUT THE TOKEN ITSELF IN THIS ENV VAR!**
- `$REPOSITORY_URL` - The GitHub URL for the repository running actions (e.g., "https://github.com/joeschmoe/my-repo")
- `$HOOK_ID` - The Hook ID for the webhook POSTing to the Cloud Run Service (e.g., "123456"); will only be validated if provided (see below).
- `$JOB_ID` (default "gha-runner") - The name of the Cloud Run job. If you change the definition of the Job in the code, you must update this value to something unique.
- `$JOB_TIMEOUT` (default "10m") - The allowed time for the action to execute.
- `$JOB_CPU` (default "1") - The CPUs allocated for the job. See https://cloud.google.com/run/docs/configuring/cpu
- `$JOB_MEMORY` (default "1Gi") - The RAM allocated for the job. See https://cloud.google.com/run/docs/configuring/memory-limits

## Setting up the `$RUNNER_IMAGE_URL`

You must make a copy of the GitHub runner image so that Cloud Run has access to it.
To do this, you must navigate to https://console.cloud.google.com/artifacts and create
a new Docker repository in the same region as your Cloud Run service. Then you can use
`docker` to make a copy of the image:

```
PROJECT=[your project, e.g., "my-project"]
LOCATION=[your location, e.g., "us-central1"]
REPOSITORY=[your newly created repository, e.g., "my-repository"]

docker pull ghcr.io/actions/actions-runner:latest

docker tag ghcr.io/actions/actions-runner:latest ${LOCATION}-docker.pkg.dev/${PROJECT}/${REPOSITORY}/actions-runner

docker push ${LOCATION}-docker.pkg.dev/${PROJECT}/${REPOSITORY}/actions-runner
```

Note the sha256 returned from the `docker push` and use this when setting the $RUNNER_IMAGE_URL env var, e.g., 
`RUNNER_IMAGE_URL=us-central1-docker.pkg.dev/some-project/some-repo/actions-runner@sha256:ABCDEF123456`

## Setting up the `$GITHUB_TOKEN_SECRET`

This app uses a "classic" Git Hub personal access token for authenticating from Cloud Run to Git Hub.

IMPORTANT: your personal access token must never be placed in an environment variable, or in any source code.

To create a personal access token, navigate to https://github.com/settings/profile, then select 
Developer Settings. Choose "Personal Access Token" -> "Tokens (classic)".

Create a new token and grant the following scopes:
**TODO:** Are these all required?

- admin:org
- repo
- workflow

Copy the generated token, then head to Google Secret Manager https://console.cloud.google.com/security/secret-manager
In the same project as the Cloud Run service, create a new secret e.g., `gha-runner`
and paste your personal access token into the secret value. You can leave all other
values as defaults.

Next you must grant Cloud Run the ability to read the secret. To do this, grant the `Secret Manager Secret Accessor` role to the 
Cloud Run service account which looks like `[your-project-number]-compute@developer.gserviceaccount.com`.

Now you can use the secret *name* in the $GITHUB_TOKEN_SECRET env var when deploying to Cloud Run:

```
GITHUB_TOKEN_SECRET=gha-runner
```

It is possible to create the secret in a different project, in which case, use the longer secret name:

```
GITHUB_TOKEN_SECRET=projects/[your-other-project]/secrets/gha-runner
```

## Building and deploying this application

Cloud Run needs a container image. The easiest way to get one is to use https://cloud.google.com/docs/buildpacks/build-application
to build this application. Choose an IMAGE_URL in Artifact Registry you created above,
e.g., `${LOCATION}-docker.pkg.dev/${PROJECT}/${REPOSITORY}/actions-manager`

You can then create a new service with this container and the environment variables mentioned above.
The service only needs minimial resources, and can scale to zero.


## Setting up the GitHub webhook and `$HOOK_ID`

First build and deploy this application to your Cloud Run service, giving you a run.app service URL 
(e.g., `https://my-service-wa5vxuhyra-uc.a.run.app`).

Nexxt set up the GitHub Actions for your repo at https://github.com/[user]/[repo]/actions.

IMPORTANT! You must edit the generated `.github/workflows/go.yml` file and set the
`runs-on` field to `self-hosted`:

```
runs-on: self-hosted
```

Finally, for the same repo, now set up the webhook at head to https://github.com/[user]/[repo]/settings.
Select Webhooks, and create a new webhook:

- **Payload URL** Your run.app URL from above.
- **Content type** Choose `application/json`.
- **Secret** Not yet supported by this app. **TODO** Support secrets.
- **Which events would you like to trigger this webhook?** "Let me select individual events." "Workflow jobs"

Leave everything else as defaults.

The HookID is an integer that will appear in the URL and you can use this to configure your `$HOOK_ID`
environment variable (will require another Cloud Run deployment; hook ID validation is optional).
