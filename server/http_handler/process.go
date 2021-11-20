package http_handler

import (
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"github.com/beoboo/job-worker-service/server/process"
	"github.com/beoboo/job-worker-service/server/runner"
	"log"
	"net/http"
)

type HttpProcessHandler struct {
	runner *runner.Runner
}

func NewHttpProcessHandler() *HttpProcessHandler {
	factory := process.ProcessFactoryImpl{}

	return &HttpProcessHandler{
		runner: runner.New(&factory),
	}
}

func (h *HttpProcessHandler) Start(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("startHandler start")
	//for name, headers := range req.Header {
	//	for _, h := range headers {
	//		_, err := fmt.Fprintf(rw, "%v: %v\n", name, h)
	//		if err != nil {
	//			fmt.Println("Cannot write response")
	//			return
	//		}
	//	}
	//}
	//

	var data protocol.StartData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	//err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Invalid JSON data: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(data)

	//h.runner.Start("./test.sh", "5 .5")
	h.runner.Start(data.Command, data.Args)

	fmt.Printf("startHandler end (%d)\n", len(h.runner.Processes))
}

func (h *HttpProcessHandler) Stop(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("stopHandler start")

	var data protocol.StopData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Invalid JSON data: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(data)

	//h.runner.Start("./test.sh", "5 .5")
	h.runner.Stop(data.Pid)

	fmt.Printf("stopHandler end (%d)\n", len(h.runner.Processes))
}
