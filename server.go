package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"job-worker-service/process"
	"job-worker-service/runner"
	"log"
	"net/http"
)

//func stream(w http.ResponseWriter, out io.ReadCloser) {
//	scanner := bufio.NewScanner(out)
//	scanner.Split(bufio.ScanWords)
//	for scanner.Scan() {
//		m := scanner.Text()
//		io.WriteString(w, m)
//		fmt.Println(m)
//	}
//}
//
//func startHandler(w http.ResponseWriter, r *http.Request) {
//	args := "5 1"
//	cmd := exec.Command("../test.sh", strings.Split(args, " ")...)
//
//	fmt.Println("Starting job...")
//	stdout, _ := cmd.StdoutPipe()
//	stderr, _ := cmd.StderrPipe()
//
//	cmd.Start()
//	cmd.Process.Release()
//
//	stream(w, stdout)
//	stream(w, stderr)
//
//	cmd.Wait()
//
//	// Write "Hello, world!" to the response body
//	io.WriteString(w, "Job completed!\n")
//	fmt.Println("Job completed...")
//}

type HttpHandler struct {
	runner *runner.Runner
}

func New() *HttpHandler {
	factory := process.ProcessFactoryImpl{}

	return &HttpHandler{
		runner: runner.New(&factory),
	}
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("startHandler start")

	h.runner.Start("./test.sh", "5 .5")

	fmt.Printf("startHandler end (%d)\n", len(h.runner.Processes))
}

func main() {
	// Set up a /hello resource handler
	http.Handle("/start", New())

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("certs/cert.pem")
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
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"))
}
