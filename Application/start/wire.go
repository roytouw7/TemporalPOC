//go:build wireinject

package main

import (
	email "TemporalTemplatePatternPOC/Application"

	"github.com/google/wire"
)

func CastToWorkflowClient(c email.ClosableWorkflowClient) email.WorkflowClient {
	return c
}

func handlerProviderConcrete() email.HandlerProvider {
	return email.NewHandlerProvider(new(email.GetRoomToUpgradeHandler), new(email.CreateUpgradeEmailMessageHandler), new(email.SendEmailHandler))
}

func initialiseApplication() (_ email.ReservationService, _ func(), err error) {
	wire.Build(
		email.NewReservationService,
		email.NewEmailService,
		email.NewEmailWorkflowService,
		email.ClientFactory,
		CastToWorkflowClient,
		handlerProviderConcrete,
	)

	return
}
