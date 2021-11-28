package main

import (
	"flag"
	"github.com/beoboo/job-scheduler/library/helpers"
	"github.com/beoboo/job-scheduler/library/job"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"os"
	"strings"
	"sync"
)

func main() {
	usage := "Usage: examples|schedule|run"

	if len(os.Args) < 2 {
		log.Fatalf(usage)
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
		log.Fatalf(usage)
	}
}

// Runs a job through the scheduler
func schedule() {
	s := scheduler.New()

	id := do(s.Start("../test.sh", "5 1"))
	log.Infof("Job \"%s\" started\n", id)

	status := do(s.Status(id))
	log.Infof("Job status: %s\n", status)

	o, err := s.Output(id)
	if err != nil {
		log.Fatalf("Cannot retrieve job: %s\n", id)
	}

	printOutput(o)

	status = do(s.Status(id))
	log.Infof("Job status: %s\n", status)

	var wg sync.WaitGroup
	wg.Wait()
}

// Runs a job without the scheduler, to verify job configuration and implementation
// Wraps the job execution through a call to /proc/self/exe in order to provide proper isolation
func run(args ...string) {
	if len(args) < 1 {
		log.Fatalf("Usage: run [--mem N] [--pids N] [--cpu N] EXECUTABLE [ARGS]\n")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		log.Fatalf("Cannot parse arguments: %s\n", err)
	}
	remaining := flag.Args()

	executable := remaining[0]
	params := remaining[1:]

	j := job.New()

	err = j.StartChild(executable, params...)
	if err != nil {
		log.Fatalf("Cannot run \"%s\": %s\n", strings.Join(params, " "), err)
	}

	log.Infof("Job \"%s\" started\n", j.Id())

	o := j.Output()

	printOutput(o)

	j.Wait()

	status := j.Status()
	log.Infof("Job status: %s\n", status)
}

// Runs a job
func child(args ...string) {
	//fs := flag.FlagSet{}
	//mem := fs.Int("mem", 0, "Max memory usage in MB")
	//pids := fs.Int("pids", 0, "Max number of processes allowed")
	//cpus := fs.Float64("cpus", 0, "Max CPU core usage")

	if len(args) < 2 {
		log.Fatalf("Usage: child [--mem N] [--pids N] [--cpu N] JOB_ID EXECUTABLE [ARGS]\n")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		log.Fatalf("Cannot parse arguments: %s\n", err)
	}
	remaining := flag.Args()

	jobId := remaining[0]
	executable := remaining[1]
	params := remaining[2:]

	log.Infof("[%s] Executing \"%s\"\n", jobId, helpers.FormatCmdLine(executable, params...))
	c := job.Child{}
	c.Run(jobId, executable, params...)
}
