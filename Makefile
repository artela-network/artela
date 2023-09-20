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
	@if ! [ -f testnet/node0/artelad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/testnet:/artela:Z artela-network/artela "./artelad testnet init-files --chain-id artela_11820-1 --v 4 -o /artela --keyring-backend=test --starting-ip-address 172.16.10.2"; fi
	docker-compose up -d

remove-testnet:
	docker-compose down
	# if ! [[ "$(docker images -q artela-network/artela:latest 2> /dev/null)" == "" ]]; then docker rmi artela-network/artela:latest; fi
	rm -rf ./testnet

clean:
	rm -rf ./build
