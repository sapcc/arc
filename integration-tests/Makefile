PKG_NAME    := gitHub.***REMOVED***/monsoon/arc
GITVERSION  := -X command-line-arguments.GITCOMMIT=`git rev-parse --short HEAD`
LDFLAGS     := -w -s $(GITVERSION)
REPOSITORY  := hub.***REMOVED***/monsoon/arc-tests
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := hub.***REMOVED***/monsoon/gobuild:1.7-alpine
BIN_DIR     := bin

ifneq ($(http_proxy),)
BUILD_ARGS+= --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) --build-arg no_proxy=$(no_proxy)
endif

.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binaries, $(BIN_DIR)"
	@echo "  * build-docker - build static test binaries via golang container, output to $(BIN_DIR)"
	@echo "  * image        - build docker image containing the tests"

.PHONY: build
build:  bin/smoke bin/updated-test bin/job-test bin/fact-test

bin/smoke: smoke_test.go helper_test.go
	@mkdir -p $(BIN_DIR)
	go test -c smoke_test.go helper_test.go -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/smoke

bin/updated-test: updated_online_service_test.go  helper_test.go
	@mkdir -p $(BIN_DIR)
	go test -c updated_online_service_test.go helper_test.go  -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/updated-test
bin/job-test: job_service_test.go helper_test.go
	@mkdir -p $(BIN_DIR)
	go test -c job_service_test.go helper_test.go -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/job-test
bin/fact-test: fact_service_test.go helper_test.go
	@mkdir -p $(BIN_DIR)
	go test -c fact_service_test.go helper_test.go  -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/fact-test

.PHONY: build-docker
build-docker:
	docker run \
		--rm \
		-v $(CURDIR)/..:/go/src/$(PKG_NAME) \
		-w /go/src/$(PKG_NAME)/integration-tests \
		$(BUILD_IMAGE) \
		make build

.PHONY: image
image: build-docker
image:
	docker build $(BUILD_ARGS) -t $(IMAGE) --rm .

.PHONY: clean
clean:
	rm -f bin/*

