package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/beoboo/job-scheduler/pkg/config"
	"github.com/beoboo/job-scheduler/pkg/protocol"
	"github.com/beoboo/job-scheduler/pkg/secret"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"io"
	"io/ioutil"
	"log"
)

type GrcpClient struct {
	client     protocol.JobSchedulerClient
	ctx        context.Context
	connection *grpc.ClientConn
}

func NewGrcpClient(addr string, enableMTLS bool) *GrcpClient {
	log.Printf("Creating GRPC client connecting to \"%s\"", addr)
	conn, err := grpc.Dial(addr, options(enableMTLS)...)
	if err != nil {
		log.Fatal(err)
	}

	client := &GrcpClient{
		client: protocol.NewJobSchedulerClient(conn),
		ctx:    context.Background(),

		connection: conn,
	}

	return client
}

func options(enableMTLS bool) []grpc.DialOption {
	if enableMTLS {
		perRPC := oauth.NewOauthAccess(createToken())

		return []grpc.DialOption{
			grpc.WithPerRPCCredentials(perRPC),
			grpc.WithTransportCredentials(loadCredentials()),
		}
	} else {
		return []grpc.DialOption{grpc.WithInsecure()}
	}
}

func createToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: secret.INCREDIBLY_SECURE,
	}
}

func loadCredentials() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair(config.CLIENT_CERT, config.CLIENT_KEY)
	if err != nil {
		log.Fatalf("Could not load client certificate/key failed: %v\n", err)
	}

	caCrt, err := ioutil.ReadFile(config.CA_CRT)
	if err != nil {
		log.Fatalf("Could not read CA file: %v\n", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCrt) {
		log.Fatalln("Could not add CA cert to the pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      caPool,
	}

	return credentials.NewTLS(tlsConfig)
}

func (c *GrcpClient) Close() {
	err := c.connection.Close()
	if err != nil {
		log.Fatalf("Unable to closeBody gRPC channel. %v\n", err)
	}
}

func (c *GrcpClient) Start(executable string, args string) (string, error) {
	request := protocol.Job{
		Executable: executable,
		Args:       args,
	}

	js, err := c.client.Start(c.ctx, &request)
	if err != nil {
		return "", err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Id, nil
}

func (c *GrcpClient) Stop(id string) (string, int32, error) {
	request := protocol.JobId{
		Id: id,
	}

	js, err := c.client.Stop(c.ctx, &request)
	if err != nil {
		return "", -1, err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Status.String(), js.ExitCode, nil
}

func (c *GrcpClient) Status(id string) (string, int32, error) {
	request := protocol.JobId{
		Id: id,
	}

	js, err := c.client.Status(c.ctx, &request)
	if err != nil {
		return "", -1, err
	}

	log.Printf("Job \"%s\" status: %s", js.Id, js.Status)
	return js.Status.String(), js.ExitCode, nil
}

func (c *GrcpClient) Output(id string) (string, error) {
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
			log.Fatalf("Error retrieving output stream, %v\n", err)
		}
		log.Printf("[%s][%s] %s", o.Time.AsTime().Format("15:04:05.000"), o.Type, o.Text)

		count += 1
	}

	return fmt.Sprintf("Output length: %d lines\n", count), nil
}
