.PHONY: all build clean lint prepare-test test ci-test

BUILD_DIR := ./build
TMP_DIR := ./tmp
APPS := $(notdir $(wildcard ./cmd/*))
COMMIT_HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +"%Y%m%d-%H%M%S")
VERSION := $(COMMIT_HASH)-$(DATE)
LDFLAGS := -s -w -X main.Version=$(VERSION)
GOFLAGS := -trimpath -mod=readonly -buildvcs=false
COVERAGE_FILE ?= coverage.out

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

test: build
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	# no clang-tidy-18 locally
	rm -rf $(TMP_DIR)/submodules/JOJ3-examples/examples/keyword/clangtidy
	go test -count=1 -v -coverpkg=./... -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE) | tail -n 1

ci-test:
	./scripts/prepare_test_repos.sh $(TMP_DIR)
	./scripts/run_foreach_test_repos.sh $(TMP_DIR) "sed -i '2i \ \ \"sandboxExecServer\": \"172.17.0.1:5051\",' conf.json"
	GITHUB_ACTOR="" go test -count=1 -v -coverpkg=./... -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE) | tail -n 1
