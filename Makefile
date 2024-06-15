.PHONY: all build clean prepare-test test

BUILD_DIR = ./build
TMP_DIR = ./tmp
APPS := $(notdir $(wildcard ./cmd/*))
FLAGS := "-s -w"

all: build

build:
	$(foreach APP,$(APPS), go build -ldflags=$(FLAGS) -o $(BUILD_DIR)/$(APP) ./cmd/$(APP);)

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf $(TMP_DIR)/*
	rm -rf *.out

prepare-test:
	git submodule update --init --remote

test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	go test -coverprofile cover.out -v ./...
