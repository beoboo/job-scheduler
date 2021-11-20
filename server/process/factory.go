package process

type ProcessFactory interface {
	Create(command string, args ...string) Process
}
