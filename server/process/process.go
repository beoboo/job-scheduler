package process

type Process interface {
	Id() string
	Start() string
	Stop() error
	Wait()
	Output() []OutputStream
	Error() []OutputStream
	Status() string
}
