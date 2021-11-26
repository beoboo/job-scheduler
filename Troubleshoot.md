# Troubleshoot

## Generate proto

Inside the `protocol` folder, run

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    job_worker_service.proto
```

Then run

```
go mod tidy
```

to download all the required packages.

## protoc-gen-go: program not found or is not executable

Install protoc-gen-go with

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

The executable will be installed in the `GOBIN` folder (usually ``{HOME}/go/bin). Make sure that folder is in your path.

## protoc-gen-go: unable to determine Go import path for

Add the package option to the proto file

```
option go_package = "example.com/package/name";
```

## Check installed cgroups versions

```
grep cgroup /proc/filesystems
```

## Check certificate

openssl x509 -text -noout -in cert.name