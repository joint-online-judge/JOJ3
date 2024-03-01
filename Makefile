.PHONY: all clean

BUILD_DIR = ./build

all:
	go build -o $(BUILD_DIR)/tiger ./cmd/tiger

clean:
	rm -rf $(BUILD_DIR)/*
	rm -rf *.out
