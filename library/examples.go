package main

import (
	"github.com/beoboo/job-scheduler/library/scheduler"
	"io"
	"sync"
)

func runExamples() {
	sched := scheduler.New(logger)

	var wg sync.WaitGroup

	// These examples will run concurrently, simulating several processes running at the same time.
	runExample(1, example1, sched, &wg)
	runExample(2, example2NoExecutable, sched, &wg)

	for i := 0; i < 10; i++ {
		runExample(i, example1, sched, &wg)
	}

	wg.Wait()
}

func runExample(id int, example func(s *scheduler.Scheduler), s *scheduler.Scheduler, wg *sync.WaitGroup) {
	infof("Example #%d", id)

	wg.Add(1)

	go func() {
		defer wg.Done()
		example(s)
	}()
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
