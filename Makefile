PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags1 = -X github.com/cosmos/cosmos-sdk/version.Name=dbChain \
       	-X github.com/cosmos/cosmos-sdk/version.ServerName=dbchaind \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=dbchaincli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS1 := -ldflags '$(ldflags1)'

ldflags2 = -X github.com/cosmos/cosmos-sdk/version.Name=dbChainCommunity \
       	-X github.com/cosmos/cosmos-sdk/version.ServerName=dbchaind \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=dbchaincli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS2 := -ldflags '$(ldflags2)'

#include Makefile.ledger
all: install

install: go.sum
	make -j 2 daemon cli

installc: go.sum
	make -j 2 daemonc clic

daemon:
	go install  $(BUILD_FLAGS1) ./cmd/dbchaind
cli:
	go install  $(BUILD_FLAGS1) ./cmd/dbchaincli

daemonc:
	go install  $(BUILD_FLAGS2) ./cmd/dbchaind
clic:
	go install  $(BUILD_FLAGS2) ./cmd/dbchaincli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)
