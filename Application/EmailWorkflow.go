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

var globalFactory HandlerProvider

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

func NewEmailWorkflowService(client WorkflowClient) EmailWorkflowService {
	globalFactory = NewHandlerProvider(&getRoomToUpgradeHandler{}, &createUpgradeEmailMessageHandler{}, &sendEmailHandler{})

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

type handlerProvider struct {
	roomUpgradeHandler handler
	createEmailHandler handler
	sendHandler        handler
}

type HandlerProvider interface {
	provideHandler(reservationId string) (*receiver, handler)
}

func NewHandlerProvider(roomUpgradeHandler handler, createEmailHandler handler, sendHandler handler) HandlerProvider {
	return &handlerProvider{
		roomUpgradeHandler: roomUpgradeHandler,
		createEmailHandler: createEmailHandler,
		sendHandler:        sendHandler,
	}
}

func (p *handlerProvider) provideHandler(reservationId string) (*receiver, handler) {
	r := &receiver{reservationId: reservationId}

	roomUpgradeHandler := p.roomUpgradeHandler
	createMailHandler := p.createEmailHandler
	sendHandler := p.sendHandler

	roomUpgradeHandler.setNext(createMailHandler).setNext(sendHandler)

	return r, roomUpgradeHandler
}

// UpgradeEmailWorkflowV3 using the Chain of Responsibility design pattern
func UpgradeEmailWorkflowV3(ctx workflow.Context, reservationId string) (bool, error) {
	//r, handler := globalFactory(reservationId)
	if globalFactory == nil {
		return false, fmt.Errorf("factory not initiated")
	}
	r, handler := globalFactory.provideHandler(reservationId)

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

type getRoomToUpgradeHandler struct {
	baseHandler
}

func (h *getRoomToUpgradeHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.GetRoomToUpgrade, r.reservationId).Get(*r.ctx, &r.room)
	}
	h.next.execute(r)
}

type createUpgradeEmailMessageHandler struct {
	baseHandler
}

// Now this piece of business logic can be separated of any Temporal knowledge and unit tested completely isolated
func (h *createUpgradeEmailMessageHandler) execute(r *receiver) {
	if r.error == nil {
		r.email = CreateUpgradeEmailMessage(r.reservationId, r.room)
	}
	h.next.execute(r)
}

type sendEmailHandler struct {
	baseHandler
}

func (h *sendEmailHandler) execute(r *receiver) {
	if r.error == nil {
		r.error = workflow.ExecuteActivity(*r.ctx, Mocks.SendEmail, r.email).Get(*r.ctx, &r.success)
	}
}
