package email

import (
	"context"
	"fmt"
	"testing"

	"TemporalTemplatePatternPOC/Mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
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

func (test *TestSuite) AfterTest(suiteName, testName string) {
	test.env.AssertExpectations(test.T())
}

// TestRun runs the test suite
func TestRun(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TestSuite))
}

type MockClient struct{}

// ExecuteWorkflow noop for passing the interface
func (c *MockClient) ExecuteWorkflow(_ context.Context, _ client.StartWorkflowOptions, _ interface{}, _ ...interface{}) (client.WorkflowRun, error) {
	return nil, nil
}

func (test *TestSuite) TestExecuteUpgradeEmailWorkflow() {
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

// TestExecuteUpgradeEmailWorkflow_mockedCreateEmailHandler tests the failure of the createUpgradeEmailMessageHandler
func (test *TestSuite) TestExecuteUpgradeEmailWorkflow_mockedCreateEmailHandler() {
	test.env.OnActivity(Mocks.GetRoomToUpgrade, mock.Anything).Return("Mocked Room", nil)

	var mockUpgradeEmailWorkflow = func(ctx workflow.Context, reservationId string) (bool, error) {
		roomUpgradeHandler := &getRoomToUpgradeHandler{}
		createMailHandler := &mockCreateUpgradeEmailMessageHandler{} // mock
		sendHandler := &sendEmailHandler{}

		roomUpgradeHandler.
			setNext(createMailHandler).
			setNext(sendHandler)

		r := createReceiver(ctx, reservationId)

		return handleUpgrade(roomUpgradeHandler, r)
	}

	id := uuid.NewString()

	test.env.ExecuteWorkflow(mockUpgradeEmailWorkflow, id)

	test.True(test.env.IsWorkflowCompleted())
	test.Contains(test.env.GetWorkflowError().Error(), "expected error")
}
