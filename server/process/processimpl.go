package process

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type ProcessImpl struct {
	command *exec.Cmd
	Pid     int
	streams map[string]string
	done    chan bool
}

func New(cmd string, args ...string) *ProcessImpl {
	fmt.Printf("NewProcessImpl for %s %s\n", cmd, strings.Join(args, " "))
	command := exec.Command(cmd, args...)

	p := &ProcessImpl{
		command: command,
		streams: make(map[string]string),
		done:    make(chan bool),
	}

	return p
}

func (p *ProcessImpl) Start() int {
	pid := make(chan int)

	fmt.Println("Running process")

	go p.run(pid)

	p.Pid = <-pid

	fmt.Printf("Process PID: %d\n", p.Pid)

	return p.Pid
}

func (p *ProcessImpl) Stop() {
	fmt.Println("Killing the process")
	err := p.command.Process.Kill()

	if err != nil {
		log.Fatalf("Cannot kill process (%+v)\n", err)
	}
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

func (p *ProcessImpl) run(pid chan int) {
	stdout, _ := p.command.StdoutPipe()
	stderr, _ := p.command.StderrPipe()

	err := p.command.Start()
	if err != nil {
		log.Fatalf("Cannot start process (%+v)\n", err)
		return
	}
	pid <- p.command.Process.Pid

	p.pipe("output", stdout)
	p.pipe("error", stderr)

	if err != nil {
		log.Fatalf("Cannot start process (%+v)\n", err)
		return
	}

	err = p.command.Wait()
	if err != nil {
		log.Fatalf("Cannot wait for process (%+v)\n", err)
		return
	}

	p.done <- true
}
