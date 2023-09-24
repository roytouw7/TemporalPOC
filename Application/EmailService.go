package email

import (
	"fmt"
)

type emailService struct {
	workflowService EmailWorkflowService
}

type EmailService interface {
	SendUpgradeEmail(reservationId string) (bool, error)
}

func NewEmailService(workflowService EmailWorkflowService) EmailService {
	return &emailService{
		workflowService: workflowService,
	}
}

// SendUpgradeEmail the only thing the email service knows,
// it gets injected a dependency which can execute a workflow requiring some data it(emailService) can provide
func (e emailService) SendUpgradeEmail(reservationId string) (bool, error) {
	return e.workflowService.executeUpgradeEmailWorkflow(reservationId)
}

func CreateUpgradeEmailMessage(reservationId, room string) string {
	return fmt.Sprintf("reservation id: \"%s\" can be upgraded to room %s!", reservationId, room)
}
