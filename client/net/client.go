package net

type Client interface {
	Start(executable string, args string) (string, error)
	Stop(id string) (string, error)
	Status(id string) (string, error)
	Output(id string) (string, error)
}
