package process

type ProcessFactoryImpl struct {}

func (f* ProcessFactoryImpl) Create(command string, args ...string) Process {
	return New(command, args...)
}
