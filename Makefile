.PHONY: all clean test

BUILD_DIR = ./build
APPS := $(notdir $(wildcard ./cmd/*))
FLAGS := "-s -w"

all:
	$(foreach APP,$(APPS), go build -ldflags=$(FLAGS) -o $(BUILD_DIR)/$(APP) ./cmd/$(APP);)

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out

test:
	go test -v ./...
