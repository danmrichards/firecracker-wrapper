GOARCH=amd64
BINARY=fcw

.PHONY: build
build:
	GOOS=linux go build -ldflags="-s -w" -o bin/${BINARY}-linux-${GOARCH} ./cmd/main.go

.PHONY: lint
lint:
	golangci-lint run ./cmd/...

.PHONY: test
test:
	go test -v -race -count=1 ./...

.PHONY: deps
deps:
	go mod verify && \
	go mod tidy

.PHONY: ubuntu-firecracker
ubuntu-firecracker:
	git submodule init
	git submodule update --remote 'ubuntu-firecracker'

.PHONY: images
images:
	cd ubuntu-firecracker && \
	docker build -t ubuntu-firecracker -f ./ubuntu-firecracker/Dockerfile ./ubuntu-firecracker && \
	docker run --privileged -it --rm -v $(pwd)/output:/output ubuntu-firecracker
