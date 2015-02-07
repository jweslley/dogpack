all: build

build:
	go build

generate:
	go generate

deps:
	go get github.com/mjibson/esc
