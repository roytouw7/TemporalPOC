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
	executeUpgradeEmailWorkflow(reservationId string) (bool, error)
}

// executeUpgradeEmailWorkflow Temporal specific method for executing the workflow, no business logic, should be fine
func (e emailWorkflowService) executeUpgradeEmailWorkflow(reservationId string) (bool, error) {
	workflowId := fmt.Sprintf("greetings-workflow-%s", reservationId)
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: "greeting-tasks",
	}

	we, err := e.client.ExecuteWorkflow(context.Background(), options, UpgradeEmailWorkflow, reservationId)
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

// TODO we can leave it as is, and accept this as how business logic meets temporal logic (downside, testability, flexibility but easy)
// TODO we can ask for a function signature and pass the business logic in (flexible, easy)
// TODO we can use the strategy pattern and pass a strategy (bit overhead, flexible, testable)
// TODO we can use the strategy pattern and supply it using dependency injection (bit overhead, more complex, user does not have to suply business logic)
// TODO we can use the template method pattern (most complicated and overhead)

// UpgradeEmailWorkflow the workflow, only location where Temporal specific code and business logic should be mixed
// also used by the worker
func UpgradeEmailWorkflow(ctx workflow.Context, reservationId string) (bool, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var room string
	err := workflow.ExecuteActivity(ctx, Mocks.GetRoomToUpgrade, reservationId).Get(ctx, &room)
	if err != nil {
		return false, err
	}

	// mixed in business logic
	emailMessage := CreateUpgradeEmailMessage(reservationId, room)

	var success bool
	err = workflow.ExecuteActivity(ctx, Mocks.SendEmail, emailMessage).Get(ctx, &success)
	if err != nil {
		return false, err
	}
	if !success {
		return success, nil
	}

	return success, nil
}
