package httphandler

import (
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"github.com/beoboo/job-worker-service/server/process"
	"github.com/beoboo/job-worker-service/server/scheduler"
	"log"
	"net/http"
)

type HttpProcessHandler struct {
	scheduler *scheduler.Scheduler
}

func NewHttpProcessHandler() *HttpProcessHandler {
	factory := process.ProcessFactoryImpl{}

	return &HttpProcessHandler{
		scheduler: scheduler.New(&factory),
	}
}

func (h *HttpProcessHandler) Start(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StartRequestData
	err := json.NewDecoder(req.Body).Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	id, status := h.scheduler.Start(data.Executable, data.Args)

	sendResponse(rw, protocol.StartResponseData{
		Id:     id,
		Status: status,
	})
}

func (h *HttpProcessHandler) Stop(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StopRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	status, err := h.scheduler.Stop(data.Id)
	if err != nil {
		sendErrorResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.StopResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpProcessHandler) Status(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StatusRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	status, err := h.scheduler.Status(data.Id)
	if err != nil {
		sendErrorResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.StatusResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpProcessHandler) Output(rw http.ResponseWriter, req *http.Request) {
	var data protocol.OutputRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Sprintf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	output, err := h.scheduler.Output(data.Id)
	if err != nil {
		sendErrorResponse(rw, err.Error(), http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.OutputResponseData{
		Id:     data.Id,
		Output: output,
	})
}

func (h *HttpProcessHandler) NotFound(rw http.ResponseWriter, req *http.Request) {
	sendErrorResponse(rw, fmt.Sprintf("No \"%s\" path found", req.URL.Path), http.StatusNotFound)
}

func sendResponse(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		sendErrorResponse(rw, fmt.Sprintf("Cannot build JSON response: %s", err), http.StatusInternalServerError)
	}
}

func sendErrorResponse(rw http.ResponseWriter, error string, code int) {
	log.Println(error)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	_ = json.NewEncoder(rw).Encode(protocol.ErrorResponseData{
		Message: error,
	})

	fmt.Fprintln(rw, error)
}
