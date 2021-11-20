package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type ProcessImpl struct {
	executable *exec.Cmd
	Pid        int
	streams    map[string]string
	done       chan bool
}

func New(cmd string, args ...string) *ProcessImpl {
	fmt.Printf("NewProcessImpl for %s %s\n", cmd, strings.Join(args, " "))
	command := exec.Command(cmd, args...)

	p := &ProcessImpl{
		executable: command,
		streams:    make(map[string]string),
		done:       make(chan bool),
	}

	return p
}

func (p *ProcessImpl) Start() int {
	pid := make(chan int)

	fmt.Println("Running process")

	go func() {
		err := p.run(pid)
		if err != nil {
			fmt.Println(err)
		}
	}()

	p.Pid = <-pid

	fmt.Printf("Process PID: %d\n", p.Pid)

	return p.Pid
}

func (p *ProcessImpl) Stop() error {
	fmt.Println("Killing the process")
	err := p.executable.Process.Kill()

	if err != nil {
		return fmt.Errorf("Cannot kill process (%+v)\n", err)
	}

	return nil
}

func (p *ProcessImpl) Wait() {
	<-p.done
}

func (p *ProcessImpl) Output() string {
	return p.streams["output"]
}

func (p *ProcessImpl) Error() string {
	return p.streams["error"]
}

func (p *ProcessImpl) Status() string {
	return "unknown status"
}

func (p *ProcessImpl) pipe(stream string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanWords)

	p.streams[stream] = ""

	for scanner.Scan() {
		m := scanner.Text()
		p.streams[stream] += m
		fmt.Println(m)
	}
}

func (p *ProcessImpl) run(pid chan int) error {
	stdout, _ := p.executable.StdoutPipe()
	stderr, _ := p.executable.StderrPipe()

	err := p.executable.Start()
	if err != nil {
		return fmt.Errorf("Cannot start process (%+v)\n", err)
	}
	pid <- p.executable.Process.Pid

	p.pipe("output", stdout)
	p.pipe("error", stderr)

	if err != nil {
		return fmt.Errorf("Cannot start process (%+v)\n", err)
	}

	err = p.executable.Wait()

	if err != nil {
		return fmt.Errorf("Cannot wait for process %d (%+v)\n", p.Pid, err)
	}

	p.done <- true

	return nil
}
