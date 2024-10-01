.PHONY: all build clean lint prepare-test test ci-test

BUILD_DIR = ./build
TMP_DIR = ./tmp
APPS := $(notdir $(wildcard ./cmd/*))
VERSION := $(shell git rev-parse --short HEAD)
FLAGS := "-s -w -X main.Version=$(VERSION)"

all: build

build:
	$(foreach APP,$(APPS), go build -ldflags=$(FLAGS) -o $(BUILD_DIR)/$(APP) ./cmd/$(APP);)

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf $(TMP_DIR)/*
	rm -rf *.out

lint:
	golangci-lint run

prepare-test:
	git submodule update --init --remote

test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	go test -coverprofile cover.out -v ./...

ci-test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	./scripts/run_foreach_test_repos.sh $(TMP_DIR) "sed -i '2i \ \ \"sandboxExecServer\": \"172.17.0.1:5051\",' conf.json"
	go test -coverprofile cover.out -v ./...
