module github.com/beoboo/job-scheduler/client

go 1.16

replace github.com/beoboo/job-scheduler/library => ../library

require (
	github.com/beoboo/job-scheduler/library v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	google.golang.org/grpc v1.42.0
)
