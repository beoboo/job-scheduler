package job

type JobFactory interface {
	Create(executable string, args ...string) Job
}
