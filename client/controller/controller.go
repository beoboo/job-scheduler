package controller

import "github.com/beoboo/job-worker-service/client/net"

type Controller struct {
	client net.Client
}

func (c *Controller) Start(command string, args string) (string, error) {
	return c.client.Start(command, args)
}

func (c *Controller) Stop(pid int) (string, error) {
	return c.client.Stop(pid)
}

func New(client net.Client) *Controller {
	return &Controller{
		client: client,
	}
}
