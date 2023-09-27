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
	client  WorkflowClient
	factory ChainOfCommandFactory
}

func NewEmailWorkflowService(client WorkflowClient, factory ChainOfCommandFactory) EmailWorkflowService {
	return &emailWorkflowService{
		client:  client,
		factory: factory,
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

	we, err := e.client.ExecuteWorkflow(context.Background(), options, UpgradeEmailWorkflowV3, e.factory, reservationId)
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

type ChainOfCommandFactory interface {
	create(reservationId string) (*receiver, []handler)
}

// chainOfCommandFactory made a factory object to allow for dependency injection, should allow for tailored unit testing by injecting a mock factory
type chainOfCommandFactory struct{}

func (f chainOfCommandFactory) create(reservationId string) (*receiver, []handler) {
	r := &receiver{reservationId: reservationId}

	roomUpgradeHandler := &getRoomToUpgradeHandler{}
	createMailHandler := &createUpgradeEmailMessageHandler{}
	sendHandler := &sendEmailHandler{}

	roomUpgradeHandler.setNext(createMailHandler)
	createMailHandler.setNext(sendHandler)

	handlers := []handler{
		roomUpgradeHandler,
		createMailHandler,
		sendHandler,
	}

	return r, handlers
}

// UpgradeEmailWorkflowV3 using the Chain of Responsibility design pattern and passed in handlers
// sadly we can not simply pass the handlers in as a workflow function can not accept functions do to them not being serializable
func UpgradeEmailWorkflowV3(ctx workflow.Context, factory ChainOfCommandFactory, reservationId string) (bool, error) {
	r, handlers := factory.create(reservationId)

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	r.ctx = &ctx

	handlers[len(handlers)-1].execute(r)

	if r.error != nil {
		return false, r.error
	}

	return r.success, nil
}

type receiver struct {
	ctx           *workflow.Context
	reservationId string
	room          string
	email         string
	success       bool
	error         error
}

type getRoomToUpgradeHandler struct {
	baseHandler
}

func (h *getRoomToUpgradeHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.GetRoomToUpgrade, r.reservationId).Get(*r.ctx, &r.room)
	}
}

type createUpgradeEmailMessageHandler struct {
	baseHandler
}

func (h *createUpgradeEmailMessageHandler) execute(r *receiver) {
	if r.error == nil {
		r.email = CreateUpgradeEmailMessage(r.reservationId, r.room)
	}
}

type sendEmailHandler struct {
	baseHandler
}

func (h *sendEmailHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.SendEmail, r.email).Get(*r.ctx, &r.success)
	}
}
