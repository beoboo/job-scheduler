package main

import (
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/scheduler"
	"io"
	"sync"
)

var logger *log.Logger

func main() {
	logger = log.New()
	sched := scheduler.New(logger)

	var wg sync.WaitGroup

	run(1, example1, sched, &wg)
	run(2, example2NoExecutable, sched, &wg)

	wg.Wait()
}

func fatalf(format string, args ...interface{}) {
	logger.Fatalf(format+"\n", args...)
}

func warnf(format string, args ...interface{}) {
	logger.Warnf(format+"\n", args...)
}

func infof(format string, args ...interface{}) {
	logger.Infof(format+"\n", args...)
}

func example1(s *scheduler.Scheduler) {
	id := do(s.Start("../test.sh", "5 1"))
	infof("Job \"%s\" started", id)

	status := do(s.Status(id))
	infof("Job status: %s", status)

	o, err := s.Output(id)
	if err != nil {
		fatalf("Cannot retrieve job: %s", id)
	}

	for {
		l, err := o.Read()
		if err == io.EOF {
			break
		}

		check(err)

		infof(l.String())
	}

	status = do(s.Status(id))
	infof("Job status: %s", status)
}

func example2NoExecutable(s *scheduler.Scheduler) {
	_, err := s.Start("../unknown", "")
	warnf("Expected error: %s", err)
}

func run(id int, example func(s *scheduler.Scheduler), s *scheduler.Scheduler, wg *sync.WaitGroup) {
	infof("Example #%d", id)

	wg.Add(1)

	go func() {
		defer wg.Done()
		example(s)
	}()
}

func do(val string, err error) string {
	check(err)

	return val
}

func check(err error) {
	if err != nil {
		fatalf("Unexpected: %s", err)
	}
}
