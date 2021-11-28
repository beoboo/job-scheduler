package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-scheduler/pkg/config"
	http2 "github.com/beoboo/job-scheduler/pkg/protocol/http"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewHttpClient(addr string, enableMTLS bool, basicAuth bool) *HttpClient {
	scheme := "HTTP"
	if enableMTLS {
		scheme = "HTTPS"
	}

	log.Printf("Creating %s client connecting to \"%s\"", scheme, addr)
	wrapped, _ := buildClient(enableMTLS)

	hc := &HttpClient{
		baseUrl: buildBaseUrl(addr, enableMTLS, basicAuth),
		client:  wrapped,
	}

	return hc
}

func buildBaseUrl(addr string, enableMTLS, basicAuth bool) string {
	credentials := ""
	if basicAuth {
		credentials = "user:test@"
	}
	if enableMTLS {
		return fmt.Sprintf("https://%s%s", credentials, addr)
	}

	return fmt.Sprintf("http://%s%s", credentials, addr)
}

func buildClient(mtls bool) (*http.Client, error) {
	if mtls {
		// Create a CA certificate pool and add cert.pem to it
		caCert, err := ioutil.ReadFile(config.SERVER_CERT)
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
	data := http2.StartRequestData{
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
	data := http2.StopRequestData{
		Id: id,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("stop", content)
}

func (c *HttpClient) Status(id string) (string, error) {
	data := http2.StatusRequestData{
		Id: id,
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return c.post("status", content)
}

func (c *HttpClient) Output(id string) (string, error) {
	data := http2.OutputRequestData{
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

	defer closeBody(response)

	fmt.Println("Status:", response.Status)
	fmt.Println("Headers:", response.Header)

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body[:]), err
}

func (c *HttpClient) Close() {

}

func closeBody(r *http.Response) {
	err := r.Body.Close()
	if err != nil {
		log.Fatal("Cannot closeBody stream")
	}
}
