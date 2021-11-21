package scheduler

import (
	"fmt"
	"github.com/beoboo/job-worker-service/protocol"
	"github.com/beoboo/job-worker-service/server/process"
	"strings"
)

type Scheduler struct {
	factory   process.ProcessFactory
	Processes map[string]process.Process
}

func New(factory process.ProcessFactory) *Scheduler {
	return &Scheduler{
		factory:   factory,
		Processes: make(map[string]process.Process),
	}
}

func (r *Scheduler) Start(executable string, args string) (string, string) {
	fmt.Printf("Starting executable: \"%s %s\"\n", executable, args)
	proc := r.factory.Create(executable, strings.Split(args, " ")...)

	id := proc.Start()
	fmt.Printf("Process PID: %s\n", id)
	fmt.Printf("Output: %s\n", proc.Output())
	fmt.Printf("Error: %s\n", proc.Error())
	fmt.Printf("Status: %s\n", proc.Status())

	r.Processes[proc.Id()] = proc
	return proc.Id(), proc.Status()
}

func (r *Scheduler) Stop(id string) (string, error) {
	fmt.Printf("Stopping process %s: \n", id)

	proc, ok := r.Processes[id]

	if !ok {
		return "", fmt.Errorf("Cannot find process: %s\n", id)
	}

	err := proc.Stop()
	if err != nil {
		return "", fmt.Errorf("Cannot stop process: %s\n", id)
	}

	delete(r.Processes, id)

	return proc.Status(), nil
}

func (r *Scheduler) Status(id string) (string, error) {
	proc, ok := r.Processes[id]

	if !ok {
		return "", fmt.Errorf("Cannot find process: %s\n", id)
	}

	return proc.Status(), nil
}

func (r *Scheduler) Output(id string) ([]protocol.OutputStream, error) {
	fmt.Printf("Retrieving output for process %s: \n", id)

	proc, ok := r.Processes[id]

	if !ok {
		return nil, fmt.Errorf("Cannot find process: %ss\n", id)
	}

	output := convertStream(proc.Output())

	return output, nil
}

func convertStream(from []process.OutputStream) []protocol.OutputStream {
	result := make([]protocol.OutputStream, len(from))

	for i, o := range from {
		result[i] = o.ToProtocol()
	}
	return result
}
