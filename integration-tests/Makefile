SHELL       := /bin/sh
GOPATH      := $(CURDIR)/../.gopath:$(CURDIR)/../Godeps/_workspace
ORG_PATH    := gitHub.***REMOVED***/monsoon
REPO_PATH   := $(ORG_PATH)/arc
GITVERSION  := -X command-line-arguments.GITCOMMIT `git rev-parse --short HEAD`
LDFLAGS     := -w -s $(GITVERSION)
REPOSITORY  := docker.***REMOVED***/monsoon/arc-test
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
build: 
	mkdir -p bin/ 
	go test -c $(CURDIR)/smoke_test.go $(CURDIR)/helper_test.go -ldflags="$(LDFLAGS)" -o bin/smoke 
	go test -c $(CURDIR)/updated_online_service_test.go $(CURDIR)/helper_test.go  -ldflags="$(LDFLAGS)" -o bin/updated-test
	go test -c $(CURDIR)/job_service_test.go $(CURDIR)/helper_test.go -ldflags="$(LDFLAGS)" -o bin/job-test
	go test -c $(CURDIR)/fact_service_test.go $(CURDIR)/helper_test.go  -ldflags="$(LDFLAGS)" -o bin/fact-test 

.PHONY: build-docker
build-docker: setup
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

# Unused (i think)
.PHONY: gopath 
gopath: setup
	@echo $(GOPATH)

.PHONY: setup
setup: ../.gopath/src/$(REPO_PATH)

../.gopath/src/$(REPO_PATH):
	mkdir -p ../.gopath/src/$(ORG_PATH)
	ln -s ../../../.. ../.gopath/src/$(REPO_PATH)
