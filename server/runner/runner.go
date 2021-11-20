package runner

import (
	"fmt"
	"github.com/beoboo/job-worker-service/server/process"
	"strings"
)

type Runner struct {
	factory   process.ProcessFactory
	Processes map[int]process.Process
}

func New(factory process.ProcessFactory) *Runner {
	return &Runner{
		factory:   factory,
		Processes: make(map[int]process.Process),
	}
}

func (r *Runner) Start(executable string, args string) (int, string) {
	fmt.Printf("Starting executable: \"%s %s\"\n", executable, args)
	proc := r.factory.Create(executable, strings.Split(args, " ")...)

	pid := proc.Start()
	fmt.Printf("Process PID: %d\n", pid)
	fmt.Printf("Output: %s", proc.Output())
	fmt.Printf("Error: %s\n", proc.Error())
	fmt.Printf("Status: %s\n", proc.Status())

	r.Processes[pid] = proc
	return pid, proc.Status()
}

func (r *Runner) Stop(pid int) (string, error) {
	fmt.Printf("Stopping process %d: \n", pid)

	proc, ok := r.Processes[pid]

	if !ok {
		return "", fmt.Errorf("Cannot find process: %d\n", pid)
	}

	err := proc.Stop()
	if err != nil {
		return "", fmt.Errorf("Cannot stop process: %d\n", pid)
	}

	delete(r.Processes, pid)

	return proc.Status(), nil
}

func (r *Runner) Status(pid int) {
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
