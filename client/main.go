package main

import (
	"flag"
	"fmt"
	"github.com/beoboo/job-worker-service/client/controller"
	"github.com/beoboo/job-worker-service/client/net"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: start|stop")
	}

	mtls := flag.Bool("mtls", false, "Enable mTLS")
	basicAuth := flag.Bool("basic-auth", false, "Enable basic auth")
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", -1, "Remote port")

	flag.Parse()
	remaining := flag.Args()
	command := remaining[0]
	args := remaining[1:]

	client := net.NewHttpClient(*host, *port, *mtls, *basicAuth)
	ctrl := controller.New(client)

	switch command {
	case "start":
		start(ctrl, args)
	case "stop":
		stop(ctrl, args)
	case "status":
		status(ctrl, args)
	case "output":
		output(ctrl, args)
	default:
		log.Fatalf("Unknown \"%s\" command\n", command)
	}
}

func start(ctrl *controller.Controller, args []string) {
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

	output, err := ctrl.Start(executable, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func stop(ctrl *controller.Controller, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: stop ID")
	}

	id := os.Args[2]

	fmt.Printf("Stopping job #\"%d\"\n", id)
	output, err := ctrl.Stop(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func status(ctrl *controller.Controller, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: status ID")
	}

	id := os.Args[2]

	fmt.Printf("Checking status for the job #\"%d\"\n", id)
	output, err := ctrl.Status(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func output(ctrl *controller.Controller, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: output ID")
	}

	id := os.Args[2]

	fmt.Printf("Retrieving output for the job #\"%d\"\n", id)
	output, err := ctrl.Output(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
