module github.com/beoboo/job-worker-service/client

go 1.16

replace github.com/beoboo/job-worker-service/protocol => ../protocol

require (
	github.com/beoboo/job-worker-service/protocol v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/grpc v1.42.0
)
