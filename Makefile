.PHONY: all build clean lint prepare-test test ci-test

BUILD_DIR = ./build
TMP_DIR = ./tmp
APPS := $(notdir $(wildcard ./cmd/*))
COMMIT_HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +"%Y%m%d-%H%M%S")
VERSION := $(COMMIT_HASH)-$(DATE)
FLAGS := "-s -w -X main.Version=$(VERSION)"

all: build

build:
	$(foreach APP,$(APPS), go build -ldflags=$(FLAGS) -o $(BUILD_DIR)/$(APP) ./cmd/$(APP);)
	cp ./build/repo-health-checker ./build/healthcheck

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf $(TMP_DIR)/*
	rm -rf *.out

lint:
	golangci-lint run -v

prepare-test:
	git submodule update --init --remote

test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	go test -coverprofile cover.out -v ./...

ci-test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	./scripts/run_foreach_test_repos.sh $(TMP_DIR) "sed -i '2i \ \ \"sandboxExecServer\": \"172.17.0.1:5051\",' conf.json"
	GITHUB_ACTIONS="test" go test -coverprofile cover.out -v ./...
