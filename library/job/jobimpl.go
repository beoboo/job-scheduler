package job

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type JobImpl struct {
	id      string
	cmd     *exec.Cmd
	streams map[string][]OutputStream
	done    chan bool
	status  string
	mtx     sync.Mutex
}

func New(executable string, args ...string) *JobImpl {
	fmt.Printf("New job for \"%s %s\"\n", executable, strings.Join(args, " "))
	cmd := exec.Command(executable, args...)
	fmt.Printf("%+v", cmd)

	// TODO: This could be a UUID or some other generated value
	now := time.Now()

	p := &JobImpl{
		cmd:     cmd,
		id:      fmt.Sprintf("%d", now.UnixNano()),
		streams: make(map[string][]OutputStream),
		done:    make(chan bool),
		status:  "idle",
	}

	return p
}

func (j *JobImpl) Id() string {
	return j.id
}

func (j *JobImpl) Start(listener OnJobListener) string {
	fmt.Println("Running job")

	started := make(chan bool, 1)

	go func() {
		err := j.run(started, listener)
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-started

	return j.id
}

func (j *JobImpl) Stop() error {
	fmt.Println("Killing the job")
	err := j.cmd.Process.Kill()

	if err != nil {
		return fmt.Errorf("cannot kill job %d: (%s)", j.pid(), err)
	}

	j.status = "killed"
	return nil
}

func (j *JobImpl) Wait() {
	<-j.done
}

func (j *JobImpl) Output() []OutputStream {
	return j.streams["output"]
}

func (j *JobImpl) Error() []OutputStream {
	return j.streams["error"]
}

func (j *JobImpl) Status() string {
	j.lock()
	defer j.unlock()
	status := j.status

	return status
}

func (j *JobImpl) lock() {
	j.mtx.Lock()
}

func (j *JobImpl) unlock() {
	j.mtx.Unlock()
}

func (j *JobImpl) run(started chan bool, listener OnJobListener) error {
	stdout, _ := j.cmd.StdoutPipe()
	stderr, _ := j.cmd.StderrPipe()

	err := j.cmd.Start()

	started <- true

	j.updateStatus("started")
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

	j.updateStatus("exited")

	if listener != nil {
		listener.OnFinishedJob(j)
	}

	j.done <- true

	return nil
}

func (j *JobImpl) pipe(stream string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)

	j.streams[stream] = []OutputStream{}

	for scanner.Scan() {
		m := scanner.Text()
		j.streams[stream] = append(j.streams[stream], OutputStream{
			Channel: stream,
			Time:    time.Now().UnixNano(),
			Text:    m,
		})
		fmt.Println(m)
	}
}

func (j *JobImpl) pid() int {
	return j.cmd.Process.Pid
}

func (j *JobImpl) updateStatus(status string) {
	j.lock()
	defer j.unlock()

	j.status = status
}
