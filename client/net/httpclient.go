package net

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewHttpClient(host string, port int, mtls bool, basicAuth bool) *HttpClient {
	wrapped, _ := buildClient(mtls)

	hc := &HttpClient{
		baseUrl: buildBaseUrl(host, buildPort(port, mtls), mtls, basicAuth),
		client:  wrapped,
	}

	return hc
}

func buildPort(port int, mtls bool) int {
	if port != -1 {
		return port
	}

	if mtls {
		return 8443
	}

	return 8080
}

func buildBaseUrl(host string, port int, mtls, basicAuth bool) string {
	credentials := ""
	if basicAuth {
		credentials = "user:test@"
	}
	if mtls {
		return fmt.Sprintf("https://%s%s:%d", credentials, host, port)
	}

	return fmt.Sprintf("http://%s%s:%d", credentials, host, port)
}

func buildClient(mtls bool) (*http.Client, error) {
	if mtls {
		// Create a CA certificate pool and add cert.pem to it
		caCert, err := ioutil.ReadFile("../certs/cert.pem")
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Read the key pair to create certificate
		cert, err := tls.LoadX509KeyPair("../certs/cert.pem", "../certs/key.pem")
		if err != nil {
			return nil, err
		}
		// Create an HTTPS client and supply the created CA pool
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		}

		return client, err
	} else {
		client := &http.Client{}

		return client, nil
	}

}

func (c *HttpClient) Start(executable string, args string) (string, error) {
	data := protocol.StartRequestData{
		Executable: executable,
		Args:       args,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("start", content)
}

func (c *HttpClient) Stop(id string) (string, error) {
	data := protocol.StopRequestData{
		Id: id,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("stop", content)
}

func (c *HttpClient) Status(id string) (string, error) {
	data := protocol.StatusRequestData{
		Id: id,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("status", content)
}

func (c *HttpClient) Output(id string) (string, error) {
	data := protocol.OutputRequestData{
		Id: id,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("output", content)
}

func (c *HttpClient) post(endpoint string, data []byte) (string, error) {
	url := fmt.Sprintf("%s/%s", c.baseUrl, endpoint)

	log.Printf("Endpoint: %s, data: %s", endpoint, data)

	response, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	defer c.close(response)

	fmt.Println("Status:", response.Status)
	fmt.Println("Headers:", response.Header)

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body[:]), err
}

func (c *HttpClient) close(r *http.Response) {
	err := r.Body.Close()
	if err != nil {
		log.Fatal("Cannot close stream")
	}
}
