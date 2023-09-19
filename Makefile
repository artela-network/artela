default_target: all

.PHONY: default_target

VERSION=$(shell git describe --tags --always)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_COMMIT_DATE=$(shell git log -n1 --pretty='format:%cd' --date=format:'%Y%m%d')
REPO=github.com/artela-network/artela

GOBIN ?= $$(go env GOPATH)/bin

ldflags = -X $(REPO)/version.AppVersion=$(VERSION) \
          -X $(REPO)/version.GitCommit=$(GIT_COMMIT) \
          -X $(REPO)/version.GitCommitDate=$(GIT_COMMIT_DATE)

build:
	go build -o build/artelad -ldflags="$(ldflags)" ./cmd/artelad/main.go

build-linux:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

install:
	go install -ldflags="$(ldflags)" ./cmd/artelad/main.go

all: build

build-testnet:
	docker build --no-cache --tag artela-network/artela ../. -f ./Dockerfile

start-testnet: remove-testnet build-testnet
	@if ! [ -f testnet/node0/artelad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build/artela:/artela:Z artela-network/artela "./artela testnet init-files --chain-id artela_11820-1 --v 4 -o /artela --keyring-backend=test --starting-ip-address 172.17.0.2"; fi
	docker-compose up -d

remove-testnet:
	docker-compose down
	# docker rm $(docker ps -q --filter ancestor=artela-network/artela:latest)
	# docker rmi artela-network/artela:latest
	rm -rf ./testnet

clean:
	rm -rf ./build
