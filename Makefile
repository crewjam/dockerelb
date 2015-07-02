
.PHONY: _dockerelb container

IMAGE_NAME=crewjam/dockerelb

all: container

container: dockerelb
	docker build -t $(IMAGE_NAME) .

dockerelb: dockerelb.go
	docker run -v $(PWD):/go/src/github.com/crewjam/dockerelb golang \
		make -C /go/src/github.com/crewjam/dockerelb _dockerelb

_dockerelb:
	go get ./...
	CGO_ENABLED=0 go install -a -installsuffix cgo -ldflags '-s' .
	ldd /go/bin/dockerelb | grep "not a dynamic executable"
	install /go/bin/dockerelb dockerelb
	
lint:
	go fmt ./...
	goimports -w *.go
