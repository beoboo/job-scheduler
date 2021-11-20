package process

type ProcessFactoryImpl struct{}

func (f *ProcessFactoryImpl) Create(executable string, args ...string) Process {
	return New(executable, args...)
}
