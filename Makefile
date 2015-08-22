PROGRAM=dogpack
VERSION=0.0.1
LDFLAGS="-X main.programVersion=$(VERSION)"

all: build

deps:
	go get ./...

tools:
	go get github.com/mjibson/esc

build:
	go build

generate:
	go generate

server:
	USE_LOCAL_FS=true ./dogpack

test: deps
	go test -v ./...

qa:
	go vet
	golint
	go test -coverprofile=.cover~
	go tool cover -html=.cover~

dist:
	@for os in linux; do \
		for arch in 386 amd64; do \
			target=$(PROGRAM)-$$os-$$arch-$(VERSION); \
			echo Building $$target; \
			GOOS=$$os GOARCH=$$arch go build -ldflags $(LDFLAGS) -o $$target/$(PROGRAM) ; \
			cp ./README.md ./LICENSE $$target; \
			tar -zcf $$target.tar.gz $$target; \
			rm -rf $$target;                   \
		done                                 \
	done

clean:
	rm -rf *.tar.gz
