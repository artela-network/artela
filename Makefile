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

all: build

clean:
	rm -rf ./build