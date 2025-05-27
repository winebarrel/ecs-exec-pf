package ecsexecpf

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockTaskDescriber struct {
	Tasks []types.Task
}

func (mtd *MockTaskDescriber) DescribeTasks(_ context.Context, _ *ecs.DescribeTasksInput, _ ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	output := &ecs.DescribeTasksOutput{}
	output.Tasks = mtd.Tasks

	return output, nil
}

var mockTaskOneContainer = []types.Task{{
	Containers: []types.Container{{
		Name:      aws.String("container_ONE"),
		RuntimeId: aws.String("containerRuntimeID_ONE"),
	}}},
}

var mockTaskMultipleContainers = []types.Task{{
	Containers: []types.Container{
		{Name: aws.String("container_ONE"), RuntimeId: aws.String("containerRuntimeID_ONE")},
		{Name: aws.String("container_TWO"), RuntimeId: aws.String("containerRuntimeID_TWO")},
		{Name: aws.String("container_THREE"), RuntimeId: aws.String("containerRuntimeID_THREE")}},
}}

func newMockTaskDescriber(tasks []types.Task) *MockTaskDescriber {
	return &MockTaskDescriber{Tasks: tasks}
}

// TestGetContainerId tests GetContainerId and getContainerIdWithDescriber for various task/container scenarios.
func TestGetContainerId(t *testing.T) {
	tests := []struct {
		name           string
		tasks          []types.Task
		containerName  string
		expectError    string
		expectedResult string
	}{
		{
			name:          "No tasks found",
			tasks:         nil,
			containerName: "testContainerName",
			expectError:   "task not found",
		},
		{
			name:          "Task has no containers",
			tasks:         []types.Task{{Containers: []types.Container{}}},
			containerName: "testContainerName",
			expectError:   "task contains no running containers",
		},
		{
			name:           "One container, no name",
			tasks:          mockTaskOneContainer,
			containerName:  "",
			expectedResult: "containerRuntimeID_ONE",
		},
		{
			name:          "Container name does not exist (one container)",
			tasks:         mockTaskOneContainer,
			containerName: "containerNameThatDoesntExist",
			expectError:   "container not found",
		},
		{
			name:          "Container name does not exist (multiple containers)",
			tasks:         mockTaskMultipleContainers,
			containerName: "containerNameThatDoesntExist",
			expectError:   "container not found",
		},
		{
			name:          "Multiple containers, no name",
			tasks:         mockTaskMultipleContainers,
			containerName: "",
			expectError:   "for tasks containing multiple containers, you must specify a container name",
		},
		{
			name:           "Find container by name (one container)",
			tasks:          mockTaskOneContainer,
			containerName:  "container_ONE",
			expectedResult: "containerRuntimeID_ONE",
		},
		{
			name:           "Find container by name (multiple containers)",
			tasks:          mockTaskMultipleContainers,
			containerName:  "container_THREE",
			expectedResult: "containerRuntimeID_THREE",
		},
		{
			name:          "First container has no runtimeId",
			tasks:         []types.Task{{Containers: []types.Container{{Name: aws.String("containerName"), RuntimeId: nil}}}},
			containerName: "",
			expectError:   "no containers with a runtimeId found",
		},
		{
			name:          "Named container has no runtimeId",
			tasks:         []types.Task{{Containers: []types.Container{{Name: aws.String("containerName"), RuntimeId: nil}}}},
			containerName: "containerName",
			expectError:   "no containers with a runtimeId found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTaskDescriber := newMockTaskDescriber(tc.tasks)
			actual, err := getContainerIdWithDescriber("testCluster", "testTaskID", tc.containerName, mockTaskDescriber)
			if tc.expectError != "" {
				require.ErrorContains(t, err, tc.expectError)
				assert.Empty(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, actual)
			}
		})
	}
}

// TestGetContainerId_DescribeTasksReturnsMultipleTasks_ReturnsError tests that an error is returned
// when DescribeTasks returns more than one task for a given task ID.
func TestGetContainerId_DescribeTasksReturnsMultipleTasks_ReturnsError(t *testing.T) {
	mockTaskDescriber := newMockTaskDescriber([]types.Task{{}, {}})
	_, err := getContainerIdWithDescriber("testCluster", "testTaskID", "", mockTaskDescriber)
	require.ErrorContains(t, err, "DescribeTasks returned multiple tasks, this should never happen")
}
