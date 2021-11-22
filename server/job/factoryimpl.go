package job

type JobFactoryImpl struct{}

func (f *JobFactoryImpl) Create(executable string, args ...string) Job {
	return New(executable, args...)
}
