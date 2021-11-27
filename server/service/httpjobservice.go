package service

import (
	"encoding/json"
	"fmt"
	"github.com/beoboo/job-scheduler/library/errors"
	http2 "github.com/beoboo/job-scheduler/library/protocol/http"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"log"
	"net/http"
	"strings"
)

type HttpJobService struct {
	scheduler *scheduler.Scheduler
}

func NewHttpJobService() *HttpJobService {
	factory := scheduler.Factory{}

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
	var data http2.StartRequestData
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

	sendResponse(rw, http2.StartResponseData{
		Id:     id,
		Status: status,
	})
}

func (h *HttpJobService) Stop(rw http.ResponseWriter, req *http.Request) {
	var data http2.StopRequestData
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

	sendResponse(rw, http2.StopResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpJobService) Status(rw http.ResponseWriter, req *http.Request) {
	var data http2.StatusRequestData
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

	sendResponse(rw, http2.StatusResponseData{
		Id:     data.Id,
		Status: status,
	})
}

func (h *HttpJobService) Output(rw http.ResponseWriter, req *http.Request) {
	var data http2.OutputRequestData
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

	sendResponse(rw, http2.OutputResponseData{
		Id:     data.Id,
		Output: converted,
	})
}

func convertOutput(from []output.OutputLine) []http2.OutputStream {
	result := make([]http2.OutputStream, len(from))

	for i, o := range from {
		result[i] = http2.OutputStream{
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

	_ = json.NewEncoder(rw).Encode(http2.ErrorResponseData{
		Message: error.Error(),
	})

	fmt.Fprintln(rw, error)
}
