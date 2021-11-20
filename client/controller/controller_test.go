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

func (c *DummyClient) Start(command string, args string) (string, error) {
	c.executed["start"] = []string{command, args}

	return "ok", nil
}

func (c *DummyClient) Stop(pid int) (string, error) {
	c.executed["stop"] = []string{strconv.Itoa(pid)}

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

	controller.Start("echo", "hello")

	args, ok := dummyClient.executed["start"]
	if !ok {
		t.Fatalf("Client has not executed the %s command", "start")
	}
	expected := []string{"echo", "hello"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}

func TestStop(t *testing.T) {
	controller := New(dummyClient)

	controller.Stop(1)

	args, ok := dummyClient.executed["stop"]
	if !ok {
		t.Fatalf("Client has not executed the %s command", "stop")
	}
	expected := []string{"1"}

	if !equal(args, expected) {
		t.Fatalf("Expected \"%s\", got \"%s\"", strings.Join(args, " "), strings.Join(expected, " "))
	}
}
