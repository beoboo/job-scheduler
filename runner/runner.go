package runner

import (
	"fmt"
	"job-worker-service/process"
	"strings"
)

type Runner struct {
	factory process.ProcessFactory
	Processes map[int]process.Process
}

func New(factory process.ProcessFactory) *Runner {
	return &Runner{
		factory: factory,
		Processes: make(map[int]process.Process),
	}
}

func (r *Runner) Start(command string, args string) int {
	fmt.Printf("Starting command: \"%s %s\"\n", command, args)
	proc := r.factory.Create(command, strings.Split(args, " ")...)

	pid := proc.Start()
	fmt.Printf("Process PID: %d\n", pid)
	fmt.Printf("Output: %s", proc.Output())
	fmt.Printf("Error: %s", proc.Error())
	fmt.Printf("Status: %s", proc.Status())

	r.Processes[pid] = proc
	return pid
}

func (r *Runner) Status(pid int) {
	fmt.Printf("Stopping process %d: \n", pid)

	proc, ok := r.Processes[pid]

	if !ok {
		_ = fmt.Errorf("Cannot find process: %d\n", pid)
		return
	}

	proc.Stop()

	delete(r.Processes, pid)
}

func (r *Runner) Stop(pid int) {
	fmt.Printf("Stopping process %d: \n", pid)

	proc, ok := r.Processes[pid]

	if !ok {
		_ = fmt.Errorf("Cannot find process: %d\n", pid)
		return
	}

	proc.Stop()

	delete(r.Processes, pid)
}

func (r *Runner) Output(pid int) string {
	fmt.Printf("Retrieving output for process %d: \n", pid)

	proc, ok := r.Processes[pid]

	if !ok {
		_ = fmt.Errorf("Cannot find process: %d\n", pid)
		return ""
	}

	return proc.Output()
}
