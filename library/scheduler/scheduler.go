package scheduler

import (
	"fmt"
	"github.com/beoboo/job-scheduler/library/errors"
	"github.com/beoboo/job-scheduler/library/job"
	"strings"
	"sync"
)

type Scheduler struct {
	factory job.JobFactory
	jobs    map[string]job.Job
	mtx     sync.Mutex
}

func New(factory job.JobFactory) *Scheduler {
	return &Scheduler{
		factory: factory,
		jobs:    make(map[string]job.Job),
	}
}

func (s *Scheduler) Start(executable string, args string) (string, string) {
	fmt.Printf("Starting executable: \"%s %s\"\n", executable, args)
	j := s.factory.Create(executable, strings.Split(args, " ")...)

	id := j.Start(s)
	fmt.Printf("Job ID: %s\n", id)
	//fmt.Printf("Output: %s\n", j.Output())
	//fmt.Printf("Error: %s\n", j.Error())
	fmt.Printf("Status: %s\n", j.Status())

	s.lock()
	defer s.unlock()

	s.jobs[j.Id()] = j
	return j.Id(), j.Status()
}

func (s *Scheduler) Stop(id string) (string, error) {
	fmt.Printf("Stopping job %s\n", id)

	j, ok := s.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	err := j.Stop()
	if err != nil {
		return "", fmt.Errorf("cannot stop job: %s", id)
	}

	delete(s.jobs, id)

	return j.Status(), nil
}

func (s *Scheduler) Status(id string) (string, error) {
	fmt.Printf("Checking status for job \"%s\"\n", id)

	s.lock()
	defer s.unlock()

	j, ok := s.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	return j.Status(), nil
}

func (s *Scheduler) Output(id string) ([]job.OutputStream, error) {
	fmt.Printf("Streaming output for job \"%s\"\n", id)

	j, ok := s.jobs[id]

	if !ok {
		return nil, &errors.NotFoundError{Id: id}
	}

	return j.Output(), nil
}

func (s *Scheduler) OnFinishedJob(j job.Job) {
	fmt.Printf("Job \"%s\" exited\n", j.Id())

	s.lock()
	defer s.unlock()

	delete(s.jobs, j.Id())
}

func (s *Scheduler) Size() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return len(s.jobs)
}

func (s *Scheduler) lock() {
	s.mtx.Lock()
}

func (s *Scheduler) unlock() {
	s.mtx.Unlock()
}
