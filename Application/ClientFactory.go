package email

import (
	"go.temporal.io/sdk/client"
)

type ClosableWorkflowClient interface {
	WorkflowClient
	Close()
}

func ClientFactory() (ClosableWorkflowClient, func(), error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, nil, err
	}

	return c, c.Close, nil
}
