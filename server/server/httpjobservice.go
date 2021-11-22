package server

import (
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"github.com/beoboo/job-worker-service/server/errors"
	"github.com/beoboo/job-worker-service/server/job"
	"github.com/beoboo/job-worker-service/server/scheduler"
	"log"
	"net/http"
	"strings"
)

type HttpJobService struct {
	scheduler *scheduler.Scheduler
}

func NewHttpJobService() *HttpJobService {
	factory := job.JobFactoryImpl{}

	service := &HttpJobService{
		scheduler: scheduler.New(&factory),
	}

	http.HandleFunc("/start", service.Start)
	http.HandleFunc("/stop", service.Stop)
	http.HandleFunc("/status", service.Status)
	http.HandleFunc("/output", service.Output)
	http.HandleFunc("/", service.NotFound)

	return service
}

func (h *HttpJobService) Start(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StartRequestData
	err := json.NewDecoder(req.Body).Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Errorf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	executable := strings.TrimSpace(data.Executable)
	if executable == "" {
		sendErrorResponse(rw, fmt.Errorf("Invalid executable: \"%s\"", executable), http.StatusBadRequest)
		return
	}

	id, status := h.scheduler.Start(data.Executable, data.Args)

	sendResponse(rw, protocol.StartResponseData{
		Id:     id,
		Status: status,
	})
}

func (h *HttpJobService) Stop(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StopRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Errorf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	status, err := h.scheduler.Stop(data.Id)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.StopResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpJobService) Status(rw http.ResponseWriter, req *http.Request) {
	var data protocol.StatusRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Errorf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	status, err := h.scheduler.Status(data.Id)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	sendResponse(rw, protocol.StatusResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpJobService) Output(rw http.ResponseWriter, req *http.Request) {
	var data protocol.OutputRequestData
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)

	if err != nil {
		sendErrorResponse(rw, fmt.Errorf("Invalid JSON data: %s", err), http.StatusBadRequest)
		return
	}

	output, err := h.scheduler.Output(data.Id)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	converted := convertOutput(output)

	sendResponse(rw, protocol.OutputResponseData{
		Id:     data.Id,
		Output: converted,
	})
}

func convertOutput(from []job.OutputStream) []protocol.OutputStream {
	result := make([]protocol.OutputStream, len(from))

	for i, o := range from {
		result[i] = protocol.OutputStream{
			Channel: o.Channel,
			Time:    int(o.Time),
			Text:    o.Text,
		}
	}
	return result
}

func (h *HttpJobService) NotFound(rw http.ResponseWriter, req *http.Request) {
	sendErrorResponse(rw, fmt.Errorf("No \"%s\" path found", req.URL.Path), http.StatusNotFound)
}

func sendResponse(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		sendErrorResponse(rw, fmt.Errorf("Cannot build JSON response: %s", err), http.StatusInternalServerError)
	}
}

func sendErrorResponse(rw http.ResponseWriter, error error, code int) {
	log.Printf("%+v", error)
	rw.Header().Set("Content-Type", "application/json")
	switch error.(type) {
	case *errors.NotFoundError:
		rw.WriteHeader(http.StatusNotFound)
	default:
		rw.WriteHeader(code)
	}

	_ = json.NewEncoder(rw).Encode(protocol.ErrorResponseData{
		Message: error.Error(),
	})

	fmt.Fprintln(rw, error)
}
