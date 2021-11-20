package net

type Client interface {
	Start(command string, args string) (string, error)
	Stop(pid int) (string, error)
}
