PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags1 = -X github.com/dbchaincloud/cosmos-sdk/version.Name=dbChain \
       	-X github.com/dbchaincloud/cosmos-sdk/version.ServerName=dbchaind \
	-X github.com/dbchaincloud/cosmos-sdk/version.ClientName=dbchaincli \
	-X github.com/dbchaincloud/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/dbchaincloud/cosmos-sdk/version.Commit=$(COMMIT) 

ldflagsoracle1 = -X github.com/dbchaincloud/cosmos-sdk/version.Name=dbChain \
       	-X github.com/dbchaincloud/cosmos-sdk/version.ServerName=dbchaind \
	-X github.com/dbchaincloud/cosmos-sdk/version.ClientName=dbchainoracle \
	-X github.com/dbchaincloud/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/dbchaincloud/cosmos-sdk/version.Commit=$(COMMIT)

ldflagsnft1 = -X github.com/dbchaincloud/cosmos-sdk/version.Name=dbChain \
        -X github.com/dbchaincloud/cosmos-sdk/version.ServerName=dbchaind \
    -X github.com/dbchaincloud/cosmos-sdk/version.ClientName=dbchainnft \
    -X github.com/dbchaincloud/cosmos-sdk/version.Version=$(VERSION) \
    -X github.com/dbchaincloud/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS1 := -ldflags '$(ldflags1)'
BUILD_FLAGS_ORACLE1 := -ldflags '$(ldflagsoracle1)'
BUILD_FLAGS_NFT1 := -ldflags '$(ldflagsnft1)'

ldflags2 = -X github.com/dbchaincloud/cosmos-sdk/version.Name=dbChainCommunity \
       	-X github.com/dbchaincloud/cosmos-sdk/version.ServerName=dbchaind \
	-X github.com/dbchaincloud/cosmos-sdk/version.ClientName=dbchaincli \
	-X github.com/dbchaincloud/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/dbchaincloud/cosmos-sdk/version.Commit=$(COMMIT) 

ldflagsoracle2 = -X github.com/dbchaincloud/cosmos-sdk/version.Name=dbChainCommunity \
        -X github.com/dbchaincloud/cosmos-sdk/version.ServerName=dbchaind \
    -X github.com/dbchaincloud/cosmos-sdk/version.ClientName=dbchainoracle \
    -X github.com/dbchaincloud/cosmos-sdk/version.Version=$(VERSION) \
    -X github.com/dbchaincloud/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS2 := -ldflags '$(ldflags2)'
BUILD_FLAGS_ORACLE2 := -ldflags '$(ldflagsoracle2)'

#include Makefile.ledger
all: install

install: go.sum
	make -j 3 nft

installc: go.sum
	make -j 2 daemonc clic oraclec

daemon:
	go install  $(BUILD_FLAGS1) ./cmd/dbchaind
cli:
	go install  $(BUILD_FLAGS1) ./cmd/dbchaincli
nft:
	go install  $(BUILD_FLAGS_NFT1) ./cmd/dbchainnft

daemonc:
	go install  $(BUILD_FLAGS2) ./cmd/dbchaind
clic:
	go install  $(BUILD_FLAGS2) ./cmd/dbchaincli
oraclec:
	go install  $(BUILD_FLAGS_ORACLE2) ./cmd/dbchainoracle

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test  $(PACKAGES)
