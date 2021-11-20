package net

type Client interface {
	Start(executable string, args string) (string, error)
	Stop(pid int) (string, error)
}
