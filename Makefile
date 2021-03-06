PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sandblockchain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=sbd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=sbcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

build: go.sum
		go build $(BUILD_FLAGS) -o build/sbd ./cmd/sbd
		go build $(BUILD_FLAGS) -o build/sbcli ./cmd/sbcli

install: go.sum
		go install  $(BUILD_FLAGS) ./cmd/sbd
		go install  $(BUILD_FLAGS) ./cmd/sbcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify