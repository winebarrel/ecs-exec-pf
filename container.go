package ecsexecpf

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ecsTaskDescriber abstracts the DescribeTasks method for dependency injection and testing.
type ecsTaskDescriber interface {
	DescribeTasks(context.Context, *ecs.DescribeTasksInput, ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
}

// GetContainerId returns the ECS container runtime ID for the specified cluster, task, and container name.
// It uses the AWS SDK ECS client created from the provided config.
// Returns an error if the task or container cannot be found, or if the container does not have a runtime ID.
func GetContainerId(cfg aws.Config, cluster, taskId, containerName string) (string, error) {
	return getContainerIdWithDescriber(cluster, taskId, containerName, ecs.NewFromConfig(cfg))
}

// getContainerIdWithDescriber looks up the container runtime ID for the given cluster, task, and container name
// using the provided ecsTaskDescriber. This function is intended for testing, allowing a mock describer to be injected.
// Returns an error if the task or container cannot be found, or if the container does not have a runtime ID.
func getContainerIdWithDescriber(cluster, taskId, containerName string, taskDescriber ecsTaskDescriber) (string, error) {

	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   []string{taskId},
	}

	output, err := taskDescriber.DescribeTasks(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("failed to call DescribeTasks: %s/%s: %w", cluster, taskId, err)
	}
	if len(output.Tasks) == 0 {
		return "", fmt.Errorf("task not found: %s/%s", cluster, taskId)
	}
	if len(output.Tasks) > 1 {
		return "", errors.New("DescribeTasks returned multiple tasks, this should never happen when passing a taskId")
	}

	task := output.Tasks[0]
	containers := task.Containers

	if len(containers) == 0 {
		return "", fmt.Errorf("task contains no running containers: %s/%s/%s", cluster, taskId, containerName)
	}

	if containerName == "" {
		if len(containers) == 1 {
			if containers[0].RuntimeId != nil {
				return *containers[0].RuntimeId, nil
			}
			return "", errors.New("no containers with a runtimeId found")
		}
		return "", fmt.Errorf("for tasks containing multiple containers, you must specify a container name: %s/%s", cluster, taskId)
	}

	for _, c := range containers {
		if c.Name != nil && *c.Name == containerName {
			if c.RuntimeId != nil {
				return *c.RuntimeId, nil
			}
			return "", errors.New("no containers with a runtimeId found")
		}
	}

	return "", fmt.Errorf("container not found: %s/%s/%s", cluster, taskId, containerName)
}
