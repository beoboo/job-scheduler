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
	fmt.Printf("New job for \"%s %s\"\n", executable, args)
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
	fmt.Println("Starting job...")

	errCh := make(chan error, 1)

	go func() {
		err := j.run(errCh)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Waits for the job to have been started successfully
	err := <-errCh

	return j.id, err
}

func (j *Job) Stop() error {
	fmt.Println("Killing the job")
	err := j.cmd.Process.Kill()

	if err != nil {
		return fmt.Errorf("cannot kill job %d: (%s)", j.pid(), err)
	}

	j.status = status.KILLED
	return nil
}

func (j *Job) Wait() {
	<-j.done
}

func (j *Job) Output() *stream.Stream {
	j.lock()
	defer j.unlock()

	j.output.ResetPos()

	return j.output
}

func (j *Job) Status() string {
	j.lock()
	defer j.unlock()

	return j.status
}

func (j *Job) lock() {
	j.m.Lock()
}

func (j *Job) unlock() {
	j.m.Unlock()
}

func (j *Job) run(started chan error) error {
	stdout, _ := j.cmd.StdoutPipe()
	stderr, _ := j.cmd.StderrPipe()

	err := j.cmd.Start()

	started <- err
	if err != nil {
		return err
	}

	j.updateStatus(status.RUNNING)
	pid := j.cmd.Process.Pid

	j.pipe("output", stdout)
	j.pipe("error", stderr)

	if err != nil {
		return fmt.Errorf("cannot start job: (%s)", err)
	}

	fmt.Printf("Job PID: %d\n", pid)

	err = j.cmd.Wait()

	if err != nil {
		return fmt.Errorf("cannot wait for job %d: (%s)", pid, err)
	}

	j.updateStatus(status.EXITED)

	j.done <- true

	return nil
}

func (j *Job) pipe(channel string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		go func() {
			j.lock()
			defer j.unlock()

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
	fmt.Printf("Updating status to %s\n", st)
	j.lock()
	defer j.unlock()

	j.status = st
	if st != status.RUNNING {
		j.output.Close()
	}
}
