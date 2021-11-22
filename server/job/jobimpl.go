package job

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type JobImpl struct {
	id      string
	cmd     *exec.Cmd
	streams map[string][]OutputStream
	done    chan bool
	status  string
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

func (p *JobImpl) Id() string {
	return p.id
}

func (p *JobImpl) Start(listener OnJobListener) string {
	fmt.Println("Running job")

	started := make(chan bool, 1)

	go func() {
		err := p.run(started, listener)
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-started

	return p.id
}

func (p *JobImpl) Stop() error {
	fmt.Println("Killing the job")
	err := p.cmd.Process.Kill()

	if err != nil {
		return fmt.Errorf("cannot kill job %d: (%s)", p.pid(), err)
	}

	p.status = "killed"
	return nil
}

func (p *JobImpl) Wait() {
	<-p.done
}

func (p *JobImpl) Output() []OutputStream {
	return p.streams["output"]
}

func (p *JobImpl) Error() []OutputStream {
	return p.streams["error"]
}

func (p *JobImpl) Status() string {
	return p.status
}

func (p *JobImpl) run(started chan bool, listener OnJobListener) error {
	stdout, _ := p.cmd.StdoutPipe()
	stderr, _ := p.cmd.StderrPipe()

	err := p.cmd.Start()

	started <- true
	p.status = "started"
	pid := p.cmd.Process.Pid

	p.pipe("output", stdout)
	p.pipe("error", stderr)

	if err != nil {
		return fmt.Errorf("cannot start job: (%s)", err)
	}

	fmt.Printf("Job PID: %d\n", pid)

	err = p.cmd.Wait()

	if err != nil {
		return fmt.Errorf("cannot wait for job %d: (%s)", pid, err)
	}

	p.status = "exited"

	if listener != nil {
		listener.OnFinishedJob(p)
	}

	p.done <- true

	return nil
}

func (p *JobImpl) pipe(stream string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)

	p.streams[stream] = []OutputStream{}

	for scanner.Scan() {
		m := scanner.Text()
		p.streams[stream] = append(p.streams[stream], OutputStream{
			Channel: stream,
			Time:    time.Now().UnixNano(),
			Text:    m,
		})
		fmt.Println(m)
	}
}

func (p *JobImpl) pid() int {
	return p.cmd.Process.Pid
}
