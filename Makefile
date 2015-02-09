all: build

build:
	go build

generate:
	go generate

server:
	USE_LOCAL_FS=true ./dogpack

deps:
	go get github.com/mjibson/esc
