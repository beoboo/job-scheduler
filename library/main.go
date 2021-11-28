package main

import (
	"flag"
	"github.com/beoboo/job-scheduler/library/helpers"
	"github.com/beoboo/job-scheduler/library/job"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"github.com/beoboo/job-scheduler/library/stream"
	"io"
	"os"
	"os/exec"
	"sync"
)

var logger *log.Logger

func main() {
	logger = log.New()

	usage := "Usage: examples|schedule|run"

	if len(os.Args) < 2 {
		fatalf(usage)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "examples":
		runExamples()
	case "schedule":
		schedule()
	case "run":
		run(args...)
	case "child":
		child(args...)
	default:
		fatalf(usage)
	}
}

// Runs a job through the scheduler
func schedule() {
	s := scheduler.New(logger)

	id := do(s.Start("../test.sh", "5 1"))
	infof("Job \"%s\" started", id)

	status := do(s.Status(id))
	infof("Job status: %s", status)

	o, err := s.Output(id)
	if err != nil {
		fatalf("Cannot retrieve job: %s", id)
	}

	printOutput(o)

	status = do(s.Status(id))
	infof("Job status: %s", status)

	var wg sync.WaitGroup
	wg.Wait()
}

// Runs a job without the scheduler, to verify job configuration and implementation
// Wraps the job execution through a call to /proc/self/exe in order to provide proper isolation
func run(args ...string) {
	if len(args) < 2 {
		fatalf("Usage: run [--mem N] [--pids N] [--cpu N] EXECUTABLE [ARGS]")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		fatalf("Cannot parse arguments: %s", err)
	}
	remaining := flag.Args()

	executable := remaining[0]
	params := remaining[1:]

	j := job.New(logger, executable, params...)
	id, err := j.Start()
	if err != nil {
		fatalf("Cannot run \"%s\": %s", helpers.FormatCmdLine(executable, params...), err)
	}

	infof("Job \"%s\" started", id)

	o := j.Output()

	printOutput(o)

	j.Wait()
}

// Runs a job
func child(args ...string) {
	//fs := flag.FlagSet{}
	//mem := fs.Int("mem", 0, "Max memory usage in MB")
	//pids := fs.Int("pids", 0, "Max number of processes allowed")
	//cpus := fs.Float64("cpus", 0, "Max CPU core usage")

	if len(args) < 2 {
		fatalf("Usage: child [--mem N] [--pids N] [--cpu N] EXECUTABLE [ARGS]")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		fatalf("Cannot parse arguments: %s", err)
	}
	remaining := flag.Args()

	executable := remaining[0]
	params := remaining[1:]

	infof("Executing \"%s\"", helpers.FormatCmdLine(executable, params...))
	cmd := exec.Command(executable, params...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fatalf("Unexpected: %s", err)
	}
}

func printOutput(o *stream.Stream) {
	for {
		l, err := o.Read()
		if err == io.EOF {
			break
		}

		check(err)

		infof(l.String())
	}
}
