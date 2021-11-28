package scheduler

import (
	"fmt"
	"github.com/beoboo/job-scheduler/library/errors"
	"github.com/beoboo/job-scheduler/library/helpers"
	"github.com/beoboo/job-scheduler/library/job"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/stream"
	"sync"
)

type Scheduler struct {
	logger *log.Logger
	jobs   map[string]*job.Job
	mtx    sync.Mutex
}

func New(logger *log.Logger) *Scheduler {
	return &Scheduler{
		logger: logger,
		jobs:   make(map[string]*job.Job),
	}
}

func (s *Scheduler) Start(executable string, args ...string) (string, error) {
	s.debug("Starting executable: \"%s\"\n", helpers.FormatCmdLine(executable, args...))
	j := job.New(s.logger, executable, args...)

	id, err := j.Start()
	if err != nil {
		return "", err
	}

	s.debug("Job ID: %s\n", id)
	s.debug("Status: %s\n", j.Status())
	//s.debug("Output: %s\n", j.Output())
	//s.debug("Error: %s\n", j.Error())

	s.lock("Start")
	defer s.unlock("Start")

	s.jobs[j.Id()] = j
	return j.Id(), err
}

func (s *Scheduler) Stop(id string) (string, error) {
	s.debug("Stopping job %s\n", id)

	s.lock("Stop")
	defer s.unlock("Stop")
	j, ok := s.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	err := j.Stop()
	if err != nil {
		return "", fmt.Errorf("cannot stop job: %s", id)
	}

	return j.Status(), nil
}

func (s *Scheduler) Status(id string) (string, error) {
	s.debug("Checking status for job \"%s\"\n", id)

	s.lock("Status")
	defer s.unlock("Status")

	j, ok := s.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	return j.Status(), nil
}

func (s *Scheduler) Output(id string) (*stream.Stream, error) {
	s.debug("Streaming output for job \"%s\"\n", id)

	j, ok := s.jobs[id]

	if !ok {
		return nil, &errors.NotFoundError{Id: id}
	}

	return j.Output(), nil
}

func (s *Scheduler) Size() int {
	s.lock("Size")
	defer s.unlock("Size")

	return len(s.jobs)
}

func (s *Scheduler) lock(id string) {
	s.debug("Scheduler locking %s", id)
	s.mtx.Lock()
}

func (s *Scheduler) unlock(id string) {
	s.debug("Scheduler unlocking %s", id)
	s.mtx.Unlock()
}

func (s *Scheduler) debug(format string, args ...interface{}) {
	if s.logger != nil {
		s.logger.Debugf(format+"\n", args...)
	}
}
