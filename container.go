package ecsexecpf

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func GetContainerId(cfg aws.Config, cluster string, taskId string, containerName string) (string, error) {
	svc := ecs.NewFromConfig(cfg)

	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   []string{taskId},
	}

	output, err := svc.DescribeTasks(context.Background(), input)

	if err != nil {
		return "", fmt.Errorf("Faild to call DescribeTasks: %s/%s", cluster, taskId)
	}

	if len(output.Tasks) == 0 {
		return "", fmt.Errorf("Task not found: %s/%s", cluster, taskId)
	}

	task := output.Tasks[0]

	if len(task.Containers) == 0 {
		return "", fmt.Errorf("Container not found: %s/%s/%s", cluster, taskId, containerName)
	}

	container := task.Containers[0]

	if containerName != "" {
		for _, c := range task.Containers {
			if c.Name == &containerName {
				container = c
			}
		}

		return "", fmt.Errorf("Container not found: %s/%s/%s", cluster, taskId, containerName)
	}

	return *container.RuntimeId, nil
}
