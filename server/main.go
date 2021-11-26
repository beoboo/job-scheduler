package main

import (
	"flag"
	"fmt"
	"github.com/beoboo/job-scheduler/server/server"
	"log"
)

func main() {
	enableMTLS := flag.Bool("mtls", false, "Enable mTLS")
	enableGRPC := flag.Bool("grpc", false, "Enable GRPC")
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", -1, "Server port")

	flag.Parse()

	addr := buildAddr(*host, *port, *enableMTLS)

	if *enableGRPC {
		log.Fatal(server.ServeGrpc(addr, *enableMTLS))
	} else {
		log.Fatal(server.ServeHttp(addr, *enableMTLS))
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
