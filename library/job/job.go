package job

import (
	"bufio"
	"fmt"
	"github.com/beoboo/job-scheduler/library/log"
	"github.com/beoboo/job-scheduler/library/status"
	"github.com/beoboo/job-scheduler/library/stream"
	"io"
	"os/exec"
	"sync"
	"time"
)

type Job struct {
	logger *log.Logger
	id     string
	cmd    *exec.Cmd
	output *stream.Stream
	done   chan bool
	status string
	m      sync.Mutex
}

func New(logger *log.Logger, executable string, args ...string) *Job {
	cmd := exec.Command(executable, args...)

	// TODO: This could be a UUID or some other random generated value
	now := time.Now()

	p := &Job{
		logger: logger,
		cmd:    cmd,
		id:     fmt.Sprintf("%d", now.UnixNano()),
		done:   make(chan bool),
		output: stream.New(logger),
		status: status.IDLE,
	}

	return p
}

func (j *Job) Id() string {
	return j.id
}

func (j *Job) Start() (string, error) {
	errCh := make(chan error, 1)

	go func() {
		j.run(errCh)
	}()

	// Waits for the job to have been started successfully
	err := <-errCh

	return j.id, err
}

func (j *Job) Stop() error {
	j.lock("Stop")
	err := j.cmd.Process.Kill()
	j.unlock("Stop")

	if err != nil {
		return fmt.Errorf("cannot kill job %d: (%s)", j.pid(), err)
	}

	j.updateStatus(status.KILLED)
	return nil
}

func (j *Job) Wait() {
	<-j.done
}

func (j *Job) Output() *stream.Stream {
	j.lock("Output")
	defer j.unlock("Output")

	j.output.ResetPos()

	return j.output
}

func (j *Job) Status() string {
	j.lock("Status")
	defer j.unlock("Status")

	return j.status
}

func (j *Job) run(started chan error) {
	stdout, _ := j.cmd.StdoutPipe()
	stderr, _ := j.cmd.StderrPipe()

	err := j.cmd.Start()

	started <- err
	if err != nil {
		return
	}

	j.updateStatus(status.RUNNING)

	j.pipe("output", stdout)
	j.pipe("error", stderr)

	_ = j.cmd.Wait()

	j.updateStatus(status.EXITED)

	j.done <- true
}

func (j *Job) pipe(channel string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		go func() {
			j.lock("pipe")
			defer j.unlock("pipe")

			_ = j.output.Write(stream.Line{
				Channel: channel,
				Time:    time.Now(),
				Text:    text,
			})
		}()
	}
}

func (j *Job) pid() int {
	return j.cmd.Process.Pid
}

func (j *Job) updateStatus(st string) {
	j.lock("updateStatus")
	defer j.unlock("updateStatus")

	switch j.status {
	case status.IDLE:
		if st == status.RUNNING {
			j.status = st
		}
	case status.RUNNING:
		if st == status.EXITED || st == status.KILLED {
			j.status = st
			j.output.Close()
		}
	default:
		return
	}
}

func (j *Job) lock(id string) {
	j.debug("Job locking %s", id)
	j.m.Lock()
}

func (j *Job) unlock(id string) {
	j.debug("Job unlocking %s", id)
	j.m.Unlock()
}

func (j *Job) debug(format string, args ...interface{}) {
	if j.logger != nil {
		j.logger.Debugf(format+"\n", args...)
	}
}
