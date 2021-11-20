package httphandler

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

	var data protocol.StartRequestData
	err := json.NewDecoder(req.Body).Decode(&data)

	if err != nil {
		log.Printf("Invalid JSON data: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pid, status := h.runner.Start(data.Executable, data.Args)

	sendResponse(rw, protocol.StartResponseData{
		Pid:    pid,
		Status: status,
	})

	fmt.Printf("startHandler end (%d)\n", len(h.runner.Processes))
}

func sendResponse(rw http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		log.Printf("Cannot build JSON response: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HttpProcessHandler) Stop(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("stopHandler start")

	var data protocol.StopRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Invalid JSON data: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := h.runner.Stop(data.Pid)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.StopResponseData{
		Pid:    data.Pid,
		Status: status,
	})

	fmt.Printf("stopHandler end (%d)\n", len(h.runner.Processes))
}
