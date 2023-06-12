package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"

	run "cloud.google.com/go/run/apiv2"
)

const (
	tokenSecretEnvVar = "TOKEN_SECRET"
)

type cloudRunJob struct {
	config config
}

func (j *cloudRunJob) ensureJob(ctx context.Context) error {
	c, err := run.NewJobsClient(ctx)
	if err != nil {
		return fmt.Errorf("creating Cloud Run client: %v", err)
	}
	defer c.Close()

	// crJob := cloudRunJob{ev: ev, config: h.config}
	req, err := j.createJobRequest()
	if err != nil {
		return fmt.Errorf("creating job request: %v", err)
	}

	op, err := c.CreateJob(ctx, req)
	if err != nil {
		// If we already have a job by this name, we're done.
		var aerr *apierror.APIError
		if errors.As(err, &aerr) && aerr.GRPCStatus().Code() == codes.AlreadyExists {
			logInfo("Job %q is already created.", j.config.JobID)
			return nil
		}
		return fmt.Errorf("creating job: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for job operation: %v", err)
	}

	logInfo("Job creation response for %q: %#v", j.config.JobID, resp)
	return nil
}

func (j *cloudRunJob) runJob(ctx context.Context, ev *event) error {
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := run.NewJobsClient(ctx)
	if err != nil {
		return fmt.Errorf("creating Cloud Run client: %v", err)
	}
	defer c.Close()

	req, err := j.runJobRequest(ev)
	if err != nil {
		return fmt.Errorf("creating job request: %v", err)
	}

	op, err := c.RunJob(ctx, req)
	if err != nil {
		return fmt.Errorf("running job: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for job operation: %v", err)
	}

	logInfo("Job run response for %q: %#v", j.config.JobID, resp)
	return nil
}

func (j *cloudRunJob) createJobRequest() (*runpb.CreateJobRequest, error) {
	req := &runpb.CreateJobRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#CreateJobRequest.
		Parent: fmt.Sprintf("projects/%s/locations/%s", j.config.Project, j.config.Location),
		JobId:  j.config.JobID,
		Job: &runpb.Job{
			// Labels map[string]string, // TODO
			// Annotations map[string]string // TODO
			// BinaryAuthorization *BinaryAuthorization
			Template: &runpb.ExecutionTemplate{
				// Labels map[string]string // TODO
				// Annotations map[string]string // TODO
				Parallelism: 0, // 0 allows maximum parallelism for the jobs
				TaskCount:   1,
				Template: &runpb.TaskTemplate{
					Containers: []*runpb.Container{
						{
							Name:  "job",
							Image: j.config.RunnerImageURL,
							// Command []string
							// Args: []string{
							// 	"./config.sh",
							// 	"--url", j.config.RepositoryHtmlURL,
							// 	"--token", "$" + tokenSecretEnvVar,
							// 	"--ephemeral",
							// 	"--disableupdate",
							//  "--unattended",
							// 	"&&",
							// 	"./run.sh",
							// },
							// Args: []string{
							// 	"/bin/bash",
							// 	"-c",
							// 	fmt.Sprintf("./config.sh --check --url %q --pat $%s && cat /home/runner/_diag/*.log", j.config.RepositoryHtmlURL, tokenSecretEnvVar),
							// },
							Args: []string{
								"nslookup", "www.google.com",
							},
							// Args: []string{
							// 	"/bin/bash",
							// 	"-c",
							// 	`echo "config.sh --help" && ./config.sh --help && echo "run.sh --help" && ./run.sh --help && echo "Token $TOKEN_SECRET"`,
							// },
							Env: []*runpb.EnvVar{
								{
									Name: tokenSecretEnvVar,
									Values: &runpb.EnvVar_ValueSource{
										ValueSource: &runpb.EnvVarSource{
											SecretKeyRef: &runpb.SecretKeySelector{
												Secret:  j.config.TokenSecretName,
												Version: "latest",
											},
										},
									},
								},
								// {
								// 	Name: "MY_PLAINTEXT",
								// 	Values: &runpb.EnvVar_Value{
								// 		Value: "the-value",
								// 	},
								// },
							},
							Resources: &runpb.ResourceRequirements{
								Limits:          map[string]string{"cpu": j.config.JobCpu, "memory": j.config.JobMemory},
								CpuIdle:         false,
								StartupCpuBoost: true,
							},
							// Ports []*ContainerPort - leave unspecified to get random port
							// VolumeMounts []*VolumeMount
							// WorkingDir
							// LivenessProbe *Probe
							// StartupProbe *Probe
						},
					},
					// Volumes: []*runpb.Volume,
					Retries: &runpb.TaskTemplate_MaxRetries{
						MaxRetries: 0,
					},
					Timeout:              durationpb.New(j.config.JobTimeout),
					ExecutionEnvironment: runpb.ExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN2,
					// ServiceAccount,
					// EncryptionKey
					// VpcAccess
				},
			},
		},
	}
	return req, nil
}

func (j *cloudRunJob) runJobRequest(ev *event) (*runpb.RunJobRequest, error) {
	return &runpb.RunJobRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#RunJobRequest.
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", j.config.Project, j.config.Location, j.config.JobID),
	}, nil
}

/*
Commands:
./config.sh Configures the runner
./config.sh remove Unconfigures the runner
./run.sh Runs the runner interactively. Does not require any options.
Options:
--help Prints the help for each command
--version Prints the runner version
--commit Prints the runner commit
--check Check the runner's network connectivity with GitHub server
Config Options:
--unattended Disable interactive prompts for missing arguments. Defaults will be used for missing options
--url string Repository to add the runner to. Required if unattended
--token string Registration token. Required if unattended
--name string Name of the runner to configure (default localhost)
--runnergroup string Name of the runner group to add this runner to (defaults to the default runner group)
--labels string Extra labels in addition to the default: 'self-hosted,Linux,X64'
--local Removes the runner config files from your local machine. Used as an option to the remove command
--work string Relative runner work directory (default _work)
--replace Replace any existing runner with the same name (default false)
--pat GitHub personal access token with repo scope. Used for checking network connectivity when executing `./run.sh --check`
--disableupdate Disable self-hosted runner automatic update to the latest released version`
--ephemeral Configure the runner to only take one job and then let the service un-configure the runner after the job finishes (default false)

Examples:
Check GitHub server network connectivity:
./run.sh --check --url <url> --pat <pat>
Configure a runner non-interactively:
./config.sh --unattended --url <url> --token <token>
Configure a runner non-interactively, replacing any existing runner with the same name:
./config.sh --unattended --url <url> --token <token> --replace [--name <name>]
Configure a runner non-interactively with three extra labels:
./config.sh --unattended --url <url> --token <token> --labels L1,L2,L3
*/
