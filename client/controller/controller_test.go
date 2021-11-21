package controller

import (
	"strconv"
	"strings"
	"testing"
)

type DummyClient struct {
	executed map[string][]string
}

func NewDummyClient() *DummyClient {
	return &DummyClient{
		executed: make(map[string][]string),
	}
}

func (c *DummyClient) Start(executable string, args string) (string, error) {
	c.executed["start"] = []string{executable, args}

	return "ok", nil
}

func (c *DummyClient) Stop(pid int) (string, error) {
	c.executed["stop"] = []string{strconv.Itoa(pid)}

	return "ok", nil
}

func (c *DummyClient) Status(pid int) (string, error) {
	c.executed["status"] = []string{strconv.Itoa(pid)}

	return "ok", nil
}

func (c *DummyClient) Output(pid int) (string, error) {
	c.executed["output"] = []string{strconv.Itoa(pid)}

	return "ok", nil
}

var dummyClient = NewDummyClient()

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestStart(t *testing.T) {
	controller := New(dummyClient)

	_, _ = controller.Start("echo", "hello")

	args, ok := dummyClient.executed["start"]
	if !ok {
		t.Fatalf("Controller has not started the \"%s\" executable", "echo")
	}
	expected := []string{"echo", "hello"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}

func TestStop(t *testing.T) {
	controller := New(dummyClient)

	_, _ = controller.Stop(1)

	args, ok := dummyClient.executed["stop"]
	if !ok {
		t.Fatalf("Controller has not stopped the %d pid", 1)
	}
	expected := []string{"1"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}

func TestStatus(t *testing.T) {
	controller := New(dummyClient)

	_, _ = controller.Status(1)

	args, ok := dummyClient.executed["status"]
	if !ok {
		t.Fatalf("Controller has not checked status for the %d pid", 1)
	}
	expected := []string{"1"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}

func TestOutput(t *testing.T) {
	controller := New(dummyClient)

	_, _ = controller.Output(1)

	args, ok := dummyClient.executed["output"]
	if !ok {
		t.Fatalf("Controller has not returned any output for the %d pid", 1)
	}

	expected := []string{"1"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}
