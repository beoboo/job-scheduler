package process

type Process interface {
	Start() int
	Stop()
	Wait()
	Output() string
	Error() string
	Status() string
}
