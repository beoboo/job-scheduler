package main

import (
	"flag"
	"github.com/beoboo/job-scheduler/library/helpers"
	"github.com/beoboo/job-scheduler/library/job"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"os"
	"strings"
)

func main() {
	usage := "Usage: examples|schedule|run"

	if !isRoot() {
		log.Fatalln("Please run this with root privileges.")
	}

	if len(os.Args) < 2 {
		log.Fatalf(usage)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "examples":
		runExamples()
	case "schedule":
		schedule(args...)
	case "run":
		run(args...)
	case "child":
		child(args...)
	default:
		log.Fatalf(usage)
	}
}

func isRoot() bool {
	return os.Geteuid() == 0
}

// Runs a job through the scheduler
func schedule(args ...string) {
	if len(args) < 1 {
		log.Fatalf("Usage: schedule [--cpu N] [--io N] [--mem N] EXECUTABLE [ARGS]\n")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		log.Fatalf("Cannot parse arguments: %s\n", err)
	}
	remaining := flag.Args()

	executable := remaining[0]
	params := remaining[1:]

	runScheduler(executable, params...)
}

func runScheduler(executable string, params ...string) {
	s := scheduler.New(true)

	id := do(s.Start(executable, params...))
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
}

// Runs a job without the scheduler, to verify job configuration and implementation
// Wraps the job execution through a call to /proc/self/exe in order to provide proper isolation
func run(args ...string) {
	if len(args) < 1 {
		log.Fatalf("Usage: run [--cpu N] [--io N] [--mem N] EXECUTABLE [ARGS]\n")
	}

	err := flag.CommandLine.Parse(args)
	if err != nil {
		log.Fatalf("Cannot parse arguments: %s\n", err)
	}
	remaining := flag.Args()

	executable := remaining[0]
	params := remaining[1:]

	startJob(err, executable, params...)
}

func startJob(err error, executable string, params ...string) {
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
	//cpus := fs.Float64("cpus", 0, "Max CPU core usage")
	//mem := fs.Int("mem", 0, "Max memory usage in MB")
	//io := fs.Int("io", 0, "Max number of processes allowed")

	if len(args) < 2 {
		log.Fatalf("Usage: child [--cpu N] [--io N] [--mem N] JOB_ID EXECUTABLE [ARGS]\n")
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
