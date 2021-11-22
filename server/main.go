package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"github.com/beoboo/job-worker-service/server/server"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	enableMTLS := flag.Bool("mtls", false, "Enable mTLS")
	enableGRPC := flag.Bool("grpc", false, "Enable GRPC")
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", -1, "Server port")

	flag.Parse()

	addr := buildAddr(*host, *port, *enableMTLS)

	if *enableGRPC {
		log.Fatal(buildGRPCServer(addr, *enableMTLS))
	} else if *enableMTLS {
		log.Fatal(buildHTTPServer(addr, *enableMTLS))
	}
}

func buildGRPCServer(addr string, enableMTLS bool) error {
	log.Printf("Creating GRPC server on \"%s\"", addr)
	srv := grpc.NewServer()

	protocol.RegisterJobSchedulerServer(srv, server.NewGrpcJobService(enableMTLS))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start listener: %v\n", err)
	}

	defer func() {
		err = listener.Close()
		if err != nil {
			log.Printf("Failed to close net listener: %v\n", err)
		}
	}()

	return srv.Serve(listener)
}

// Listen to HTTP(S) connections, with the server certificate if MTLS is enabled
func buildHTTPServer(addr string, enableMTLS bool) error {
	scheme := buildScheme(enableMTLS)

	log.Printf("Creating %s server on \"%s\"", scheme, addr)

	// Create a Server instance to listen on port 8443 with the TLS config
	srv := &http.Server{
		Addr:      addr,
		TLSConfig: buildTlsConfig(enableMTLS),
	}

	_ = server.NewHttpJobService()

	if enableMTLS {
		return srv.ListenAndServeTLS("../certs/cert.pem", "../certs/key.pem")
	} else {
		return srv.ListenAndServe()
	}
}

func buildScheme(enableMTLS bool) string {
	if enableMTLS {
		return "HTTPS"
	}
	return "HTTP"
}

func buildTlsConfig(enableMtls bool) *tls.Config {
	if !enableMtls {
		return nil
	}

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("../certs/cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	return &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
}

func buildAddr(host string, port int, enableMTLS bool) string {
	port = buildPort(port, enableMTLS)
	return fmt.Sprintf("%s:%d", host, port)
}

func buildPort(port int, enableMTLS bool) int {
	if port != -1 {
		return port
	}

	if enableMTLS {
		return 8443
	}

	return 8080
}
