VERSION=0.0.1

all: build

build:
	go build

generate:
	go generate

server:
	USE_LOCAL_FS=true ./dogpack

deps:
	go get github.com/mjibson/esc
	go get github.com/gobuild/gobuild3/packer

dist:
	packer --os linux --arch amd64 --output dogpack-linux-amd64-$(VERSION).zip
	packer --os linux --arch 386 --output dogpack-linux-386-$(VERSION).zip

clean:
	rm -f *.zip
