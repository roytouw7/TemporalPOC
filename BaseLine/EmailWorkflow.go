package email

import (
	"context"
	"fmt"
	"log"
	"time"

	"TemporalTemplatePatternPOC/Mocks"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type WorkflowClient interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

type emailWorkflowService struct {
	client WorkflowClient
}

func NewEmailWorkflowService(client WorkflowClient) EmailWorkflowService {
	return &emailWorkflowService{
		client: client,
	}
}

type EmailWorkflowService interface {
	executeUpgradeEmailWorkflow(reservationId, msg string) (bool, error)
}

func (e emailWorkflowService) executeUpgradeEmailWorkflow(reservationId, msg string) (bool, error) {
	workflowId := fmt.Sprintf("greetings-workflow-%s", reservationId)
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: "greeting-tasks",
	}

	we, err := e.client.ExecuteWorkflow(context.Background(), options, UpgradeEmailWorkflow, msg)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var success bool
	err = we.Get(context.Background(), &success)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", success)

	return success, nil
}

func UpgradeEmailWorkflow(ctx workflow.Context, emailMessage string) (bool, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var success bool
	err := workflow.ExecuteActivity(ctx, Mocks.SendEmail, emailMessage).Get(ctx, &success)
	if err != nil {
		return false, err
	}

	return success, nil
}
