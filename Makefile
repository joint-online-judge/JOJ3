.PHONY: all clean test

BUILD_DIR = ./build
APPS := $(notdir $(wildcard ./cmd/*))

all:
	$(foreach APP,$(APPS), go build -o $(BUILD_DIR)/$(APP) ./cmd/$(APP);)

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out

test:
	go test -v ./...
