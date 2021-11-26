.ONESHELL:

.PHONY: test build run

test:
	cd ${COMPONENT} && go test ${VERBOSE} ./...

build:
	cd ${COMPONENT} && go build .

run:
	cd ${COMPONENT} && go run ${VERBOSE} main.go ${ARGS}

docker-build:
	docker build . -t golang

docker-run:
	docker run -it golang
