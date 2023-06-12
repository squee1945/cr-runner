package main

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

type cloudRunJob struct {
	config config
	ev     *event
}

func (j *cloudRunJob) jobID() string {
	return strconv.Itoa(j.ev.WorkflowJob.ID)
}

func (j *cloudRunJob) createJobRequest() (*runpb.CreateJobRequest, error) {
	req := &runpb.CreateJobRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/run/apiv2/runpb#CreateJobRequest.
		Parent: fmt.Sprintf("projects/%s/locations/%s", j.config.project, j.config.location),
		JobId:  j.jobID(),
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
								"--url", j.ev.Repository.HtmlURL,
								"--token", "$" + gitHubTokenSecretEnvVar,
								"--ephemeral",
								"--disableupdate",
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
