.DEFAULT_GOAL: help

.PHONY: help build modules test build_docker test_docker clean lint push kube_yaml

DOCKER          ?= docker
DOCKER_TAG      := test
GO              ?= go
MAKEFILE_PATH   := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
APP				:= app

help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dep-dev:
	echo "Would run dep-dev"

setup:
	$(GO) mod init github.com/AClarkie/timescale
	make modules
	make format

modules: ## Tidy up and update vendor dependencies
	$(GO) mod tidy
	$(GO) mod vendor

format:
	$(GO) fmt $$($(GO) list ./...)

build: ## Builds the app
	echo "To build $(APP) binary"
	$(GO) build -o $(APP) ./cmd

test: modules ## Run the tests
	echo "To run tests"
	## Run tests in current dir and subdirectory and disable test caching
	$(GO) test -coverprofile=coverage.out -mod=vendor ./... -count=1 -v -coverprofile=coverage.out

clean: ## Clean the object files, build, test caches along with the built binary
	$(GO) clean -cache -testcache
	rm -rf $(APP)

docker: ## Build the docker image for timescale
	docker build \
		-t aclarkie/timescale:$(DOCKER_TAG) \
		.
