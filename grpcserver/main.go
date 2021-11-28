package main

import (
	"flag"
	"fmt"
	"github.com/beoboo/job-scheduler/grpcserver/server"
	"log"
)

func main() {
	enableMTLS := flag.Bool("mtls", false, "Enable mTLS")
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", -1, "Server port")

	flag.Parse()

	addr := buildAddr(*host, *port, *enableMTLS)

	log.Fatal(server.ServeGrpc(addr, *enableMTLS))
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
