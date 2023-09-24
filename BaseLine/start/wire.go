//go:build wireinject

package main

import (
	email "TemporalTemplatePatternPOC/BaseLine"

	"github.com/google/wire"
)

func CastToWorkflowClient(c email.ClosableWorkflowClient) email.WorkflowClient {
	return c
}

func initialiseApplication() (_ email.ReservationService, _ func(), err error) {
	wire.Build(
		email.NewReservationService,
		email.NewEmailService,
		email.NewEmailWorkflowService,
		email.ClientFactory,
		CastToWorkflowClient,
	)

	return
}
