.PHONY: all clean

BUILD_DIR = ./build
APP_NAME = joj3

all:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out
