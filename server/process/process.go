package process

type Process interface {
	Start() int
	Stop() error
	Wait()
	Output() string
	Error() string
	Status() string
}
