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
	jobs map[string]*job.Job
	mtx  sync.Mutex
}

// New creates a scheduler.
func New() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*job.Job),
	}
}

// Start runs a new job.Job.
func (s *Scheduler) Start(executable string, args ...string) (string, error) {
	log.Debugf("Starting executable: \"%s\"\n", helpers.FormatCmdLine(executable, args...))
	j := job.New()

	err := j.StartChild(executable, args...)
	if err != nil {
		return "", err
	}

	log.Debugf("Job ID: %s\n", j.Id())
	log.Debugf("Status: %s\n", j.Status())
	//log.Debugf(utput: %s\n", j.Output())
	//log.Debugf(rror: %s\n", j.Error())

	s.lock("Start")
	defer s.unlock("Start")

	s.jobs[j.Id()] = j
	return j.Id(), err
}

// Stop stops a running job.Job, or an error if the job.Job doesn't exist.
func (s *Scheduler) Stop(id string) (string, error) {
	log.Debugf("Stopping job %s\n", id)

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

// Status returns the status of a job.Job, or an error if the job.Job doesn't exist.
func (s *Scheduler) Status(id string) (string, error) {
	log.Debugf("Checking status for job \"%s\"\n", id)

	s.lock("Status")
	defer s.unlock("Status")

	j, ok := s.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	return j.Status(), nil
}

// Output returns the stream of the stdout/stderr of a job.Job, or an error if the job.Job doesn't exist.
func (s *Scheduler) Output(id string) (*stream.Stream, error) {
	log.Debugf("Streaming output for job \"%s\"\n", id)

	s.lock("Output")
	defer s.unlock("Output")
	j, ok := s.jobs[id]

	if !ok {
		return nil, &errors.NotFoundError{Id: id}
	}

	return j.Output(), nil
}

// Size returns the number of stored jobs.
func (s *Scheduler) Size() int {
	s.lock("Size")
	defer s.unlock("Size")

	return len(s.jobs)
}

func (s *Scheduler) lock(id string) {
	log.Tracef("Scheduler locking %s\n", id)
	s.mtx.Lock()
}

func (s *Scheduler) unlock(id string) {
	log.Tracef("Scheduler unlocking %s\n", id)
	s.mtx.Unlock()
}
