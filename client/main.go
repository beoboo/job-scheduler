package main

import (
	"flag"
	"fmt"
	"github.com/beoboo/job-worker-service/client/controller"
	"github.com/beoboo/job-worker-service/client/net"
	"log"
	"os"
	"strconv"
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
	executable := remaining[0]
	args := remaining[1:]

	client := net.NewHttpClient(*host, *port, *mtls, *basicAuth)
	ctrl := controller.New(client)

	switch executable {
	case "start":
		start(ctrl, args)
	case "stop":
		stop(ctrl, args)
	default:
		log.Fatalf("Unknown \"%s\" executable\n", executable)
	}
}

func start(ctrl *controller.Controller, args []string) {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: start executable [args]")
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
		log.Fatalf("Usage: stop pid")
	}

	pid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Pid must be an int")
	}

	fmt.Printf("Stopping job \"%d\"\n", pid)
	output, err := ctrl.Stop(pid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
