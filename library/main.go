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

	wg.Wait()
}

func fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

func info(format string, args ...interface{}) {
	logger.Infof(format+"\n", args...)
}

func example1(s *scheduler.Scheduler) {
	id := do(s.Start("../test.sh", "5 1"))
	info("Job \"%s\" started\n", id)

	status := do(s.Status(id))
	info("Job status: %s", status)

	o, err := s.Output(id)
	if err != nil {
		fatal("Cannot retrieve job: %s", id)
	}

	for {
		l, err := o.Read()
		if err == io.EOF {
			break
		}

		check(err)

		info(l.String())
	}

	status = do(s.Status(id))
	info("Job status: %s", status)
}

func run(id int, example func(s *scheduler.Scheduler), s *scheduler.Scheduler, wg *sync.WaitGroup) {
	info("Example #%d", id)

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
		fatal(err)
	}
}
