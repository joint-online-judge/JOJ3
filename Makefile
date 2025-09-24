.PHONY: all build clean lint prepare-test test ci-test

BUILD_DIR := ./build
TMP_DIR := ./tmp
APPS := $(notdir $(wildcard ./cmd/*))
COMMIT_HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +"%Y%m%d-%H%M%S")
VERSION := $(COMMIT_HASH)-$(DATE)
LDFLAGS := -s -w -X main.Version=$(VERSION)
GOFLAGS := -trimpath -mod=readonly -buildvcs=false

all: build

build:
	$(foreach APP,$(APPS), \
		CGO_ENABLED=0 \
		go build $(GOFLAGS) -ldflags='$(LDFLAGS)' -o $(BUILD_DIR)/$(APP) ./cmd/$(APP) \
		|| exit 1; \
	)

clean:
	$(RM) -rv $(BUILD_DIR) $(TMP_DIR) *.out

lint:
	golangci-lint run -v

prepare-test:
	git submodule update --init --remote

test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	# no students.csv & clang-tidy-18 locally
	rm -rf $(TMP_DIR)/submodules/JOJ3-examples/examples/healthcheck
	rm -rf $(TMP_DIR)/submodules/JOJ3-examples/examples/keyword/clangtidy
	go test -count=1 -v ./...

local-test:
	rm -rf $(TMP_DIR)/submodules/JOJ3-examples/examples/
	mkdir -p $(TMP_DIR)/submodules/JOJ3-examples/examples/
	go test -count=1 -v ./...

ci-test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	./scripts/run_foreach_test_repos.sh $(TMP_DIR) "sed -i '3i \ \ \"sandboxExecServer\": \"172.17.0.1:5051\",' conf.json"
	GITHUB_ACTOR="" go test -count=1 -v ./...
