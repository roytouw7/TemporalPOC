package main

import (
	"log"

	email "TemporalTemplatePatternPOC/Application"
	"TemporalTemplatePatternPOC/Mocks"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "greeting-tasks", worker.Options{})

	w.RegisterWorkflow(email.UpgradeEmailWorkflow)

	w.RegisterActivity(Mocks.SendEmail)
	w.RegisterActivity(Mocks.GetRoomToUpgrade)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
