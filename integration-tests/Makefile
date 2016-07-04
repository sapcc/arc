SHELL       := /bin/sh
ORG_PATH    := gitHub.***REMOVED***/monsoon
REPO_PATH   := $(ORG_PATH)/arc
GITVERSION  := -X command-line-arguments.GITCOMMIT `git rev-parse --short HEAD`
LDFLAGS     := -w -s $(GITVERSION)
REPOSITORY  := docker.***REMOVED***/monsoon/arc-tests
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/arc-build

DOCKER = docker

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/api-server"
	@echo "  * build-docker - build static binary via golang container, output to bin/api-server"
	@echo "  * image        - build docker image containing the update site"

.PHONY: build 
build:  bin/smoke bin/updated-test bin/job-test bin/fact-test

bin/smoke: smoke_test.go helper_test.go ensure_gopath
	@mkdir -p bin/ 
	go test -c smoke_test.go helper_test.go -ldflags="$(LDFLAGS)" -o bin/smoke 
bin/updated-test: updated_online_service_test.go  helper_test.go ensure_gopath
	@mkdir -p bin/
	go test -c updated_online_service_test.go helper_test.go  -ldflags="$(LDFLAGS)" -o bin/updated-test
bin/job-test: job_service_test.go helper_test.go ensure_gopath
	@mkdir -p bin/
	go test -c job_service_test.go helper_test.go -ldflags="$(LDFLAGS)" -o bin/job-test
bin/fact-test: fact_service_test.go helper_test.go ensure_gopath
	@mkdir -p bin/
	go test -c fact_service_test.go helper_test.go  -ldflags="$(LDFLAGS)" -o bin/fact-test

.PHONY: build-docker
build-docker:
	$(DOCKER) run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/integration-tests build


.PHONY: image
image: build-docker 
image: 
	$(DOCKER) pull $(REPOSITORY):latest || true
	$(DOCKER) build -t $(IMAGE) --rm . 
	echo $(IMAGE) > image 

.PHONY: clean
clean:
	rm -f bin/*
	rm image


.PHONY: ensure_gopath
ensure_gopath:
	@goDir=$${GOPATH%%:*}/src/$(PKG_NAME) && \
				mkdir -p $$(dirname $$goDir) && \
				if [ ! -e "$$goDir" ]; then \
					ln -sfv "$(CURDIR)" "$$goDir"; \
				fi
