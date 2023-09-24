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

func (e emailService) SendUpgradeEmail(reservationId string) (bool, error) {
	msg := e.createUpgradeEmailMessage(reservationId)
	return e.workflowService.executeUpgradeEmailWorkflow(reservationId, msg)
}

func (e emailService) createUpgradeEmailMessage(reservationId string) string {
	return fmt.Sprintf("reservation id: \"%s\" can be upgraded!", reservationId)
}
