package job

import (
	"bufio"
	"fmt"
	"github.com/beoboo/job-scheduler/library/status"
	"github.com/beoboo/job-scheduler/library/stream"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Job struct {
	id     string
	cmd    *exec.Cmd
	output *stream.Stream
	done   chan bool
	status string
	m      sync.Mutex
}

func NewJob(executable string, args string) *Job {
	cmd := exec.Command(executable, strings.Split(args, " ")...)

	// TODO: This could be a UUID or some other generated value
	now := time.Now()

	p := &Job{
		cmd:    cmd,
		id:     fmt.Sprintf("%d", now.UnixNano()),
		done:   make(chan bool),
		output: stream.NewStream(),
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
	j.lock(1)
	err := j.cmd.Process.Kill()
	j.unlock(1)

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
	j.lock(2)
	defer j.unlock(2)

	j.output.ResetPos()

	return j.output
}

func (j *Job) Status() string {
	j.lock(3)
	defer j.unlock(3)

	return j.status
}

func (j *Job) lock(id int) {
	//println("locking %d", id)
	j.m.Lock()
}

func (j *Job) unlock(id int) {
	//println("unlocking %d", id)
	j.m.Unlock()
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
			j.lock(4)
			defer j.unlock(4)

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
	j.lock(5)
	defer j.unlock(5)

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
