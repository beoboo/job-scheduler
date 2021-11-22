package client

import (
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
)

type GrcpJobClient struct {
	client     protocol.JobSchedulerClient
	ctx        context.Context
	connection *grpc.ClientConn
}

func NewGrcpClient(addr string, opts ...grpc.DialOption) *GrcpJobClient {
	log.Printf("Creating GRPC client connecting to \"%s\"", addr)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := &GrcpJobClient{
		client:     protocol.NewJobSchedulerClient(conn),
		ctx:        context.Background(),
		connection: conn,
	}

	return client
}

func (c *GrcpJobClient) Close() {
	err := c.connection.Close()
	if err != nil {
		log.Printf("Unable to closeBody gRPC channel. %v", err)
	}
}

func (c *GrcpJobClient) Start(executable string, args string) (string, error) {
	request := protocol.Job{
		Executable: executable,
		Args:       args,
	}

	js, err := c.client.Start(c.ctx, &request)
	if err != nil {
		return "", err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Status, nil
}

func (c *GrcpJobClient) Stop(id string) (string, error) {
	request := protocol.JobId{
		Id: id,
	}

	js, err := c.client.Stop(c.ctx, &request)
	if err != nil {
		return "", err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Status, nil
}

func (c *GrcpJobClient) Status(id string) (string, error) {
	request := protocol.JobId{
		Id: id,
	}

	js, err := c.client.Status(c.ctx, &request)
	if err != nil {
		return "", err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Status, nil
}

func (c *GrcpJobClient) Output(id string) (string, error) {
	request := protocol.JobId{
		Id: id,
	}

	stream, err := c.client.Output(c.ctx, &request)
	if err != nil {
		return "", err
	}

	count := 0

	for {
		o, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error retrieving output stream, %v", err)
		}
		log.Printf("[%s][%s] %s", o.Time, o.Channel, o.Text)

		count += 1
	}

	return fmt.Sprintf("Output length: %d lines", count), nil
}
