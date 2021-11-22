package main

import (
	"flag"
	"fmt"
	"github.com/beoboo/job-worker-service/client/client"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
)

func main() {
	enableMTLS := flag.Bool("mtls", false, "Enable mTLS")
	enableGRPC := flag.Bool("grpc", false, "Enable GRPC")
	basicAuth := flag.Bool("basic-auth", false, "Enable basic auth")
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", -1, "Remote port")

	flag.Parse()
	remaining := flag.Args()

	if len(remaining) == 0 {
		log.Fatalf("Usage: start|stop")
	}

	command := remaining[0]
	args := remaining[1:]

	addr := buildAddr(*host, *port, *enableMTLS)

	var clnt client.Client
	if *enableGRPC {
		clnt = client.NewGrcpClient(addr, grpc.WithInsecure())
	} else {
		clnt = client.NewHttpClient(addr, *enableMTLS, *basicAuth)
	}

	defer clnt.Close()

	switch command {
	case "start":
		start(clnt, args)
	case "stop":
		stop(clnt, args)
	case "status":
		status(clnt, args)
	case "output":
		output(clnt, args)
	default:
		log.Fatalf("Unknown \"%s\" command\n", command)
	}
}

func buildAddr(host string, port int, enableMTLS bool) string {
	return fmt.Sprintf("%s:%d", host, buildPort(port, enableMTLS))
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

func start(clnt client.Client, args []string) {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: start EXECUTABLE [ARGS]")
	}

	executable := args[0]
	params := strings.Join(args[1:], " ")

	if len(params) > 0 {
		fmt.Printf("Starting \"%s %s\"\n", executable, params)
	} else {
		fmt.Printf("Starting \"%s\"\n", executable)
	}

	output, err := clnt.Start(executable, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func stop(clnt client.Client, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: stop ID")
	}

	id := args[0]

	fmt.Printf("Stopping job \"%s\"\n", id)
	output, err := clnt.Stop(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func status(clnt client.Client, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: status ID")
	}

	id := args[0]

	fmt.Printf("Checking status for the job \"%s\"\n", id)
	output, err := clnt.Status(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func output(clnt client.Client, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: output ID")
	}

	id := args[0]

	fmt.Printf("Retrieving output for the job \"%s\"\n", id)
	output, err := clnt.Output(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
