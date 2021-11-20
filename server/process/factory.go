package process

type ProcessFactory interface {
	Create(executable string, args ...string) Process
}
