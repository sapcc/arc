SHELL       := /bin/sh
GOPATH      := $(CURDIR)/../.gopath:$(CURDIR)/../Godeps/_workspace
ORG_PATH    := gitHub.***REMOVED***/monsoon
REPO_PATH   := $(ORG_PATH)/arc
GITVERSION  := -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT `git rev-parse --short HEAD`
REPOSITORY  := docker.***REMOVED***/monsoon/arc-update-site
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/arc-build

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/update-site"
	@echo "  * build-docker - build static binary via golang container, output to bin/update-site"
	@echo "  * image        - build docker image containing the update site"

.PHONY: build
build: setup
	mkdir -p bin/ 
	CGO_ENABLED=0 go build -o bin/update-site -a --installsuffix cgo -ldflags="-w -s $(GITVERSION)" $(REPO_PATH)/update-server

.PHONY: build-docker
build-docker:
	docker run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/update-server build

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
	docker pull $(REPOSITORY):latest || true
	docker build -t $(IMAGE) --rm . 
	echo $(IMAGE) > image 

.PHONY: clean
clean:
	rm -f bin/update-site
