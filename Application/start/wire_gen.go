// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"TemporalTemplatePatternPOC/Application"
)

// Injectors from wire.go:

func initialiseApplication() (email.ReservationService, func(), error) {
	closableWorkflowClient, cleanup, err := email.ClientFactory()
	if err != nil {
		return nil, nil, err
	}
	workflowClient := CastToWorkflowClient(closableWorkflowClient)
	handlerProvider := handlerProviderConcrete()
	emailWorkflowService := email.NewEmailWorkflowService(workflowClient, handlerProvider)
	emailService := email.NewEmailService(emailWorkflowService)
	reservationService := email.NewReservationService(emailService)
	return reservationService, func() {
		cleanup()
	}, nil
}

// wire.go:

func CastToWorkflowClient(c email.ClosableWorkflowClient) email.WorkflowClient {
	return c
}

func handlerProviderConcrete() email.HandlerProvider {
	return email.NewHandlerProvider(new(email.GetRoomToUpgradeHandler), new(email.CreateUpgradeEmailMessageHandler), new(email.SendEmailHandler))
}
