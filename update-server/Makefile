SHELL       := /bin/sh
GOPATH      := $(CURDIR)/../.gopath:$(CURDIR)/../Godeps/_workspace
ORG_PATH    := gitHub.***REMOVED***/monsoon
REPO_PATH   := $(ORG_PATH)/arc
LDFLAGS     := -s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT `git rev-parse --short HEAD`
REPOSITORY  := docker.***REMOVED***/monsoon/arc-update-site
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/arc-build

ifneq ($(BUILD_VERSION),)
LDFLAGS:=$(LDFLAGS) -X gitHub.***REMOVED***/monsoon/arc/version.Version $(BUILD_VERSION)
endif

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
	go build -o bin/update-site -ldflags="$(LDFLAGS)" $(REPO_PATH)/update-server

.PHONY: build-docker
build-docker:
	docker run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/update-server build BUILD_VERSION=$(BUILD_VERSION) 

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

.PHONY: fmt
fmt:
	which goimports > /dev/null
	dirs=`go list -f "{{.Dir}}" ./...`; \
		test -z "`for d in $$dirs; do goimports -l $$d/*.go | tee /dev/stderr; done`"

.PHONY: test
test: fmt unit

.PHONY: unit
unit: setup
	go test -v -timeout=4s ./...

.PHONY: clean
clean:
	rm -f bin/update-site
