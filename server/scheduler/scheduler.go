package scheduler

import (
	"fmt"
	"github.com/beoboo/job-worker-service/server/errors"
	"github.com/beoboo/job-worker-service/server/job"
	"strings"
)

type Scheduler struct {
	factory job.JobFactory
	jobs    map[string]job.Job
}

func New(factory job.JobFactory) *Scheduler {
	return &Scheduler{
		factory: factory,
		jobs:    make(map[string]job.Job),
	}
}

func (r *Scheduler) Start(executable string, args string) (string, string) {
	fmt.Printf("Starting executable: \"%s %s\"\n", executable, args)
	j := r.factory.Create(executable, strings.Split(args, " ")...)

	id := j.Start(r)
	fmt.Printf("Job ID: %s\n", id)
	//fmt.Printf("Output: %s\n", j.Output())
	//fmt.Printf("Error: %s\n", j.Error())
	fmt.Printf("Status: %s\n", j.Status())

	r.jobs[j.Id()] = j
	return j.Id(), j.Status()
}

func (r *Scheduler) Stop(id string) (string, error) {
	fmt.Printf("Stopping job %s\n", id)

	j, ok := r.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	err := j.Stop()
	if err != nil {
		return "", fmt.Errorf("cannot stop job: %s", id)
	}

	delete(r.jobs, id)

	return j.Status(), nil
}

func (r *Scheduler) Status(id string) (string, error) {
	fmt.Printf("Checking status for job \"%s\"\n", id)

	j, ok := r.jobs[id]

	if !ok {
		return "", &errors.NotFoundError{Id: id}
	}

	return j.Status(), nil
}

func (r *Scheduler) Output(id string) ([]job.OutputStream, error) {
	fmt.Printf("Streaming output for job \"%s\"\n", id)

	j, ok := r.jobs[id]

	if !ok {
		return nil, &errors.NotFoundError{Id: id}
	}

	return j.Output(), nil
}

func (r *Scheduler) OnFinishedJob(j job.Job) {
	fmt.Printf("Job \"%s\" exited\n", j.Id())
	delete(r.jobs, j.Id())
}
