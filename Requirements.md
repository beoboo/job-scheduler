# Requirements

## Level 1

### Library

* Worker library with methods to start/stop/query status and get the output of a job.

### API

* HTTPS API to start/stop/get status of a running process.
* Use HTTP Basic Authentication.
* Use a simple authorization scheme.

### CLI

* CLI should be able to connect to worker service and start, stop, get status, and output of a job.

## Level 2

### Library

* Worker library with methods to start/stop/query status and get the output of a job.

### API

* HTTPS API to start/stop/get status of a running process.
* Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
  setup for certificates. Do not use any other authentication protocols on top of mTLS.
* Use a simple authorization scheme.

### CLI

* CLI should be able to connect to worker service and start, stop, get status, and output of a job.

## Level 3

### Library

* Worker library with methods to start/stop/query status and get the output of a job.
* Library should be able to stream the output of a running job.
    * Output should be from start of process execution.
    * Multiple concurrent clients should be supported.

### API

* GRPC API to start/stop/get status/stream output of a running process.
* Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
  setup for certificates. Do not use any other authentication protocols on top of mTLS.
* Use a simple authorization scheme.

### Client

* CLI should be able to connect to worker service and start, stop, get status, and stream output of a job.

## Level 4

### Library

* Worker library with methods to start/stop/query status and get the output of a job.
* Library should be able to stream the output of a running job.
    * Output should be from start of process execution.
    * Multiple concurrent clients should be supported.
* Add resource control for CPU, Memory and Disk IO per job using cgroups.

### API

* GRPC API to start/stop/get status/stream output of a running process.
* Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
  setup for certificates. Do not use any other authentication protocols on top of mTLS.
* Use a simple authorization scheme.

### Client

* CLI should be able to connect to worker service and start, stop, get status, and stream output of a job.

## Level 5

### Library

* Worker library with methods to start/stop/query status and get the output of a job.
* Library should be able to stream the output of a running job.
    * Output should be from start of process execution.
    * Multiple concurrent clients should be supported.
* Add resource control for CPU, Memory and Disk IO per job using cgroups.
* Add resource isolation for using PID, mount, and networking namespaces.

### API

* GRPC API to start/stop/get status/stream output of a running process.
* Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
  setup for certificates. Do not use any other authentication protocols on top of mTLS.
* Use a simple authorization scheme.

### Client

* CLI should be able to connect to worker service and start, stop, get status, and stream output of a job.
