package server

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/beoboo/job-scheduler/server/service"
	"io/ioutil"
	"log"
	"net/http"
)

// ServeHttp listens to HTTP(S) connections, with the server certificate if MTLS is enabled
func ServeHttp(addr string, enableMTLS bool) error {
	scheme := buildScheme(enableMTLS)

	log.Printf("Creating %s server on \"%s\"", scheme, addr)

	// Create a Server instance to listen on port 8443 with the TLS config
	srv := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig(enableMTLS),
	}

	_ = service.NewHttpJobService()

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

func tlsConfig(enableMTLS bool) *tls.Config {
	if !enableMTLS {
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
