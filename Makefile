GOARCH=amd64
NIC?=eth0

.PHONY: build
build: build-wrapper build-manager

.PHONY: build-wrapper
build-wrapper:
	GOOS=linux go build -ldflags="-s -w" -o bin/wrapper-linux-${GOARCH} ./cmd/wrapper/main.go

.PHONY: build-manager
build-manager:
	GOOS=linux go build -ldflags="-s -w" -o bin/manager-linux-${GOARCH} ./cmd/manager/main.go

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

.PHONY: tap
tap:
	sudo ip tuntap add tap0 mode tap && \
	sudo ip addr add 172.16.0.1/24 dev tap0 && \
    sudo ip link set tap0 up && \
    sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward" && \
    sudo iptables -t nat -A POSTROUTING -o ${NIC} -j MASQUERADE && \
    sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT && \
    sudo iptables -A FORWARD -i tap0 -o ${NIC} -j ACCEPT

.PHONY: cleantap
cleantap:
	sudo ip link del tap0
