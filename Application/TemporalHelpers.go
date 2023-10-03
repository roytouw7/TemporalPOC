package email

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

type MockClient struct{}

// ExecuteWorkflow noop for passing the interface
func (c *MockClient) ExecuteWorkflow(_ context.Context, _ client.StartWorkflowOptions, _ interface{}, _ ...interface{}) (client.WorkflowRun, error) {
	we := &WorkFlowRunStruct{}
	return we, nil
}

type WorkFlowRunStruct struct{}

func (we *WorkFlowRunStruct) Get(ctx context.Context, result interface{}) error {
	successPtr, ok := result.(*bool)
	if !ok {
		return fmt.Errorf("result is not a pointer to a bool")
	}

	*successPtr = true
	return nil
}

func (we *WorkFlowRunStruct) GetID() string {
	return ""
}

func (we *WorkFlowRunStruct) GetRunID() string {
	return ""
}

func (we *WorkFlowRunStruct) GetWithOptions(ctx context.Context, valuePtr interface{}, options client.WorkflowRunGetOptions) error {
	return nil
}
