package server

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/beoboo/job-scheduler/grpcserver/service"
	"github.com/beoboo/job-scheduler/pkg/config"
	"github.com/beoboo/job-scheduler/pkg/protocol"
	"github.com/beoboo/job-scheduler/pkg/secret"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func ServeGrpc(addr string, enableMTLS bool) error {
	log.Printf("Creating gRPC server on \"%s\" (MTLS: %v)", addr, enableMTLS)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start listener: %v\n", err)
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Failed to close listener: %v\n", err)
		}
	}()

	srv := grpc.NewServer(options(enableMTLS)...)

	protocol.RegisterJobSchedulerServer(srv, service.NewGrpcJobService())

	return srv.Serve(listener)
}

func options(enableMTLS bool) []grpc.ServerOption {
	if !enableMTLS {
		return []grpc.ServerOption{
			grpc.EmptyServerOption{},
		}
	}

	data, err := ioutil.ReadFile(config.CA_CRT)
	if err != nil {
		log.Fatalf("Can't read CA file: %v\n", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(data) {
		log.Fatalln("Can't add CA cert")
	}

	cert, err := tls.LoadX509KeyPair(config.SERVER_CERT, config.SERVER_KEY)
	if err != nil {
		log.Fatalf("Cannot load server credentials: %v\n", err)
	}

	tlsConfig := &tls.Config{
		ClientCAs:  caPool,
		ClientAuth: tls.RequireAndVerifyClientCert,

		Certificates: []tls.Certificate{cert},
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
		},
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
	}

	return []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.UnaryInterceptor(authorize),
	}
}

func authorize(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if !validateToken(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token ")
	}

	return handler(ctx, req)
}

func validateToken(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	return token == secret.INCREDIBLY_SECURE
}
