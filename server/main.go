package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/beoboo/job-worker-service/server/http_handler"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Set up a /hello resource handler
	handler := http_handler.NewHttpProcessHandler()

	//http.Handle("/start", handler)
	http.HandleFunc("/start", handler.Start)
	http.HandleFunc("/stop", handler.Stop)
	//http.Handle("/stop", http_handler.NewHttpProcessHandler())
	//http.Handle("/start", http_handler.NewHttpBasicHandler(http_handler.NewHttpProcessHandler()))

	mtls := flag.Bool("mtls", false, "Enable mTLS")
	port := flag.Int("port", -1, "Server port")

	flag.Parse()

	addr := buildAddr(*port, *mtls)

	if *mtls {
		// Create a CA certificate pool and add cert.pem to it
		caCert, err := ioutil.ReadFile("../certs/cert.pem")
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Create the TLS Config with the CA pool and enable Client certificate validation
		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}

		// Create a Server instance to listen on port 8443 with the TLS config
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
		}

		// Listen to HTTPS connections with the server certificate and wait
		log.Fatal(server.ListenAndServeTLS("../certs/cert.pem", "../certs/key.pem"))
	} else {
		server := &http.Server{
			Addr: addr,
		}

		// Listen to HTTP connections with no certificates and wait
		log.Printf("Listening on %s", addr)
		log.Fatal(server.ListenAndServe())
	}
}

func buildAddr(port int, mtls bool) string {
	if port != -1 {
		return fmt.Sprintf(":%d", port)
	}

	if mtls {
		return ":8443"
	}

	return ":8080"
}
