package job

type Job interface {
	Id() string
	Start(listener OnJobListener) string
	Stop() error
	Wait()
	Output() []OutputStream
	Error() []OutputStream
	Status() string
}

type OnJobListener interface {
	OnFinishedJob(job Job)
}
