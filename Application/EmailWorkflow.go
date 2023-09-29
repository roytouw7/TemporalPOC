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

var globalHandlerProvider HandlerProvider

// WorkflowClient limit our knowledge of Temporal
type WorkflowClient interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

type EmailWorkflowService interface {
	executeUpgradeEmailWorkflow(reservationId string) (bool, error)
}

type emailWorkflowService struct {
	client WorkflowClient
}

func NewEmailWorkflowService(client WorkflowClient, factory HandlerProvider) EmailWorkflowService {
	globalHandlerProvider = factory // TODO I'm not happy with exposing the injected provider to the global scope to use it in the workflow again... But it works?

	return &emailWorkflowService{
		client: client,
	}
}

// executeUpgradeEmailWorkflow Temporal specific method for executing the workflow, no business logic, should be fine
func (e emailWorkflowService) executeUpgradeEmailWorkflow(reservationId string) (bool, error) {
	workflowId := fmt.Sprintf("greetings-workflow-%s", reservationId)
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: "greeting-tasks",
	}

	we, err := e.client.ExecuteWorkflow(context.Background(), options, UpgradeEmailWorkflowV3, reservationId)
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

type HandlerProvider interface {
	createChainOfCommand(reservationId string) (*receiver, handler)
}

type handlerProvider struct {
	getRoomToUpgradeHandler          handler
	createUpgradeEmailMessageHandler handler
	sendEmailHandler                 handler
}

func NewHandlerProvider(getRoomToUpgradeHandler handler, createUpgradeEmailMessageHandler handler, sendEmailHandler handler) HandlerProvider {
	return &handlerProvider{
		getRoomToUpgradeHandler:          getRoomToUpgradeHandler,
		createUpgradeEmailMessageHandler: createUpgradeEmailMessageHandler,
		sendEmailHandler:                 sendEmailHandler,
	}
}

// createChainOfCommand create orchestration of handlers
func (f *handlerProvider) createChainOfCommand(reservationId string) (*receiver, handler) {
	r := &receiver{reservationId: reservationId}

	roomUpgradeHandler := f.getRoomToUpgradeHandler
	createMailHandler := f.createUpgradeEmailMessageHandler
	sendHandler := f.sendEmailHandler

	roomUpgradeHandler.
		setNext(createMailHandler).
		setNext(sendHandler)

	return r, roomUpgradeHandler
}

// UpgradeEmailWorkflowV3 using the Chain of Responsibility design pattern
// sadly we can not simply pass the handlers in as a workflow function can not accept functions do to them not being serializable
func UpgradeEmailWorkflowV3(ctx workflow.Context, reservationId string) (bool, error) {
	r, handler := globalHandlerProvider.createChainOfCommand(reservationId)

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	r.ctx = &ctx

	handler.execute(r)

	if r.error != nil {
		return false, r.error
	}

	return r.success, nil
}

// Orchestration code below, a handler for every part of the process

type receiver struct {
	ctx           *workflow.Context
	reservationId string
	room          string
	email         string
	success       bool
	error         error
}

type GetRoomToUpgradeHandler struct {
	baseHandler
}

func (h *GetRoomToUpgradeHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.GetRoomToUpgrade, r.reservationId).Get(*r.ctx, &r.room)
	}
	h.next.execute(r)
}

type CreateUpgradeEmailMessageHandler struct {
	baseHandler
}

// Now this piece of business logic can be separated of any Temporal knowledge and unit tested completely isolated
func (h *CreateUpgradeEmailMessageHandler) execute(r *receiver) {
	if r.error == nil {
		r.email = CreateUpgradeEmailMessage(r.reservationId, r.room)
	}
	h.next.execute(r)
}

type SendEmailHandler struct {
	baseHandler
}

func (h *SendEmailHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.SendEmail, r.email).Get(*r.ctx, &r.success)
	}
}
