SHELL       := /bin/sh
GOPATH      := $(CURDIR)/../.gopath:$(CURDIR)/../Godeps/_workspace
ORG_PATH    := gitHub.***REMOVED***/monsoon
REPO_PATH   := $(ORG_PATH)/arc
GITVERSION  := -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT `git rev-parse --short HEAD`
REPOSITORY  := docker.***REMOVED***/monsoon/arc-api
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/arc-build

DOCKER = docker
WAIT   = $(DOCKER) run --rm --link $(WAIT_ID):wait  \
				 docker.***REMOVED***/monsoon/docker-build:1.4.0 \
				 wait

### Variables that are expanded dynamically
postgres = $(shell cat postgres 2> /dev/null)

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/api-server"
	@echo "  * build-docker - build static binary via golang container, output to bin/api-server"
	@echo "  * image        - build docker image containing the update site"

.PHONY: build
build: setup
	mkdir -p bin/ 
	go build -o bin/api-server -ldflags="-w -s $(GITVERSION)" $(REPO_PATH)/api-server

.PHONY: gencert
gencert:
	go get -d gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert
	go build -o bin/gencert -ldflags="-s -w" gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert

.PHONY: build-docker
build-docker: setup
	$(DOCKER) run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/api-server build gencert

.PHONY: gopath 
gopath: setup
	@echo $(GOPATH)

.PHONY: setup
setup: ../.gopath/src/$(REPO_PATH)

../.gopath/src/$(REPO_PATH):
	mkdir -p ../.gopath/src/$(ORG_PATH)
	ln -s ../../../.. ../.gopath/src/$(REPO_PATH)

.PHONY: image
image: build-docker 
image: 
	$(DOCKER) pull $(REPOSITORY):latest || true
	$(DOCKER) build -t $(IMAGE) --rm . 
	echo $(IMAGE) > image 

.PHONY: test
test: postgres migrate-test
	$(DOCKER) run \
		--rm \
		--link $(postgres):postgres \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/api-server run-tests

.PHONY: run-tests
run-tests: setup
	go test ./... -v
	go test ./... -v -p=1 -timeout=60s -tags=integration

.PHONY:
migrate-%:
	$(DOCKER) exec $(postgres) createdb -U postgres arc_$* || true
	$(DOCKER) run --rm \
		--link $(postgres):postgres \
		-v $(CURDIR):/api-server \
		-w /api-server \
		$(BUILD_IMAGE) \
	  goose --env=$* up 

# ----------------------------------------------------------------------------------
#   postgres 
# ----------------------------------------------------------------------------------
#
# Start postgres database and wait for it to become available. 
#
postgres: WAIT_ID = $$(cat postgres)
postgres: 
	$(DOCKER) run -d postgres > postgres 
	$(WAIT)

.PHONY: clean
clean:
	rm -f bin/api-server
	rm -f bin/gencert
