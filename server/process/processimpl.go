package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type ProcessImpl struct {
	id      string
	cmd     *exec.Cmd
	streams map[string][]OutputStream
	done    chan bool
	status  string
}

func New(executable string, args ...string) *ProcessImpl {
	fmt.Printf("New process for \"%s %s\"\n", executable, strings.Join(args, " "))
	cmd := exec.Command(executable, args...)
	fmt.Printf("%+v", cmd)

	// TODO: This could be a UUID or some other generated value
	now := time.Now()

	p := &ProcessImpl{
		cmd:     cmd,
		id:      fmt.Sprintf("%d", now.UnixNano()),
		streams: make(map[string][]OutputStream),
		done:    make(chan bool),
		status:  "idle",
	}

	return p
}

func (p *ProcessImpl) Id() string {
	return p.id
}

func (p *ProcessImpl) Start() string {
	fmt.Println("Running process")

	started := make(chan bool, 1)

	go func() {
		err := p.run(started)
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-started

	return p.id
}

func (p *ProcessImpl) Stop() error {
	fmt.Println("Killing the process")
	err := p.cmd.Process.Kill()

	if err != nil {
		return fmt.Errorf("cannot kill process %d: (%s)", p.pid(), err)
	}

	p.status = "killed"
	return nil
}

func (p *ProcessImpl) Wait() {
	<-p.done
}

func (p *ProcessImpl) Output() []OutputStream {
	return p.streams["output"]
}

func (p *ProcessImpl) Error() []OutputStream {
	return p.streams["error"]
}

func (p *ProcessImpl) Status() string {
	return p.status
}

func (p *ProcessImpl) run(started chan bool) error {
	stdout, _ := p.cmd.StdoutPipe()
	stderr, _ := p.cmd.StderrPipe()

	err := p.cmd.Start()

	started <- true
	p.status = "started"
	pid := p.cmd.Process.Pid

	p.pipe("output", stdout)
	p.pipe("error", stderr)

	if err != nil {
		return fmt.Errorf("cannot start process: (%s)", err)
	}

	fmt.Printf("Process PID: %d\n", pid)

	err = p.cmd.Wait()

	if err != nil {
		return fmt.Errorf("cannot wait for process %d: (%s)", pid, err)
	}

	p.status = "exited"
	p.done <- true

	return nil
}

func (p *ProcessImpl) pipe(stream string, pipe io.ReadCloser) {
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

func (p *ProcessImpl) pid() int {
	return p.cmd.Process.Pid
}
