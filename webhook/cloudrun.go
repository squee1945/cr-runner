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
			logInfo("Job %q is already created.", j.config.jobID)
			return nil
		}
		return fmt.Errorf("creating job: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for job operation: %v", err)
	}

	logInfo("Job creation response for %q: %#v", j.config.jobID, resp)
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

	logInfo("Job run response for %q: %#v", j.config.jobID, resp)
	return nil
}

func (j *cloudRunJob) createJobRequest() (*runpb.CreateJobRequest, error) {
	req := &runpb.CreateJobRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#CreateJobRequest.
		Parent: fmt.Sprintf("projects/%s/locations/%s", j.config.project, j.config.location),
		JobId:  j.config.jobID,
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
							Image: j.config.runnerImageURL,
							// Command []string
							Args: []string{
								"./config.sh",
								"--url", j.config.repositoryHtmlURL,
								"--token", "$" + gitHubTokenSecretEnvVar,
								"--ephemeral",
								"--disableupdate",
								"&&",
								"./run.sh",
							},
							Env: []*runpb.EnvVar{
								{
									Name: gitHubTokenSecretEnvVar,
									Values: &runpb.EnvVar_ValueSource{
										ValueSource: &runpb.EnvVarSource{
											SecretKeyRef: &runpb.SecretKeySelector{
												Secret:  j.config.tokenSecretName,
												Version: "latest",
											},
										},
									},
								},
								{
									Name: "MY_PLAINTEXT",
									Values: &runpb.EnvVar_Value{
										Value: "the-value",
									},
								},
							},
							Resources: &runpb.ResourceRequirements{
								Limits:          map[string]string{"cpu": j.config.jobCpu, "memory": j.config.jobMemory},
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
					Timeout:              durationpb.New(j.config.jobTimeout),
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
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", j.config.project, j.config.location, j.config.jobID),
	}, nil
}
