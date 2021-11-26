module github.com/beoboo/job-scheduler/server

go 1.16

replace github.com/beoboo/job-scheduler/library => ../library

require (
	github.com/beoboo/job-scheduler/library v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.42.0
)
