package controller

import "github.com/beoboo/job-worker-service/client/net"

type Controller struct {
	client net.Client
}

func (c *Controller) Start(executable string, args string) (string, error) {
	return c.client.Start(executable, args)
}

func (c *Controller) Stop(id string) (string, error) {
	return c.client.Stop(id)
}

func (c *Controller) Status(id string) (string, error) {
	return c.client.Status(id)
}

func (c *Controller) Output(id string) (string, error) {
	return c.client.Output(id)
}

func New(client net.Client) *Controller {
	return &Controller{
		client: client,
	}
}
