package email

import (
	"fmt"
	"testing"

	"TemporalTemplatePatternPOC/Mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

type TestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	mockClient WorkflowClient
	env        *testsuite.TestWorkflowEnvironment

	service EmailWorkflowService
}

func (test *TestSuite) SetupTest() {
	test.env = test.NewTestWorkflowEnvironment()
	test.mockClient = &MockClient{}
	test.service = NewEmailWorkflowService(test.mockClient)
}

func (test *TestSuite) AfterTest() {
	test.env.AssertExpectations(test.T())
}

// TestRun runs the test suite
func TestRun(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TestSuite))
}

// TestUpgradeEmailWorkflow tests the workflow happy flow
func (test *TestSuite) TestUpgradeEmailWorkflow() {
	test.env.OnActivity(Mocks.GetRoomToUpgrade, mock.Anything).Return("Mocked Room", nil)
	test.env.OnActivity(Mocks.SendEmail, mock.Anything).Return(true, nil)

	id := uuid.NewString()

	test.env.ExecuteWorkflow(UpgradeEmailWorkflowV3, id)

	test.True(test.env.IsWorkflowCompleted())
	test.NoError(test.env.GetWorkflowError())
}

type mockCreateUpgradeEmailMessageHandler struct {
	baseHandler
}

func (h *mockCreateUpgradeEmailMessageHandler) execute(r *receiver) {
	r.error = fmt.Errorf("expected error")
	h.next.execute(r)
}

// TestUpgradeEmailWorkflow_mockedCreateEmailHandler tests the failure of the createUpgradeEmailMessageHandler
func (test *TestSuite) TestUpgradeEmailWorkflow_mockedCreateEmailHandler() {
	test.env.OnActivity(Mocks.GetRoomToUpgrade, mock.Anything).Return("Mocked Room", nil)

	var mockUpgradeEmailWorkflow = func(ctx workflow.Context, reservationId string) (bool, error) {
		roomUpgradeHandler := &getRoomToUpgradeHandler{}
		createMailHandler := &mockCreateUpgradeEmailMessageHandler{} // mock
		sendHandler := &sendEmailHandler{}

		roomUpgradeHandler.
			setNext(createMailHandler).
			setNext(sendHandler)

		r := createReceiver(reservationId)
		r.ctx = &ctx

		return handleUpgrade(roomUpgradeHandler, r)
	}

	id := uuid.NewString()

	test.env.ExecuteWorkflow(mockUpgradeEmailWorkflow, id)

	test.True(test.env.IsWorkflowCompleted())
	test.Contains(test.env.GetWorkflowError().Error(), "expected error")
}

// TestExecuteUpgradeEmailWorkflow Tests the happy flow of executing the workflow
func (test *TestSuite) TestExecuteUpgradeEmailWorkflow() {
	success, err := test.service.executeUpgradeEmailWorkflow(uuid.NewString())

	test.True(success)
	test.NoError(err)
}
