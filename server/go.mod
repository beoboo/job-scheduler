module github.com/beoboo/job-worker-service/server

go 1.16

replace github.com/beoboo/job-worker-service/protocol => ../protocol

require (
	github.com/beoboo/job-worker-service/protocol v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.42.0
)
