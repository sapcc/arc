export GO15VENDOREXPERIMENT=1
REPO_PATH   := gitHub.***REMOVED***/monsoon/arc
LDFLAGS     := -s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT=`git rev-parse --short HEAD`
REPOSITORY  := docker.***REMOVED***/monsoon/arc-update-site
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/gobuild@1.5.3

ifneq ($(BUILD_VERSION),)
LDFLAGS:=$(LDFLAGS) -X gitHub.***REMOVED***/monsoon/arc/version.Version=$(BUILD_VERSION)
endif

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/update-site"
	@echo "  * build-docker - build static binary via golang container, output to bin/update-site"
	@echo "  * image        - build docker image containing the update site"

.PHONY: ensure_gopath
ensure_gopath:
	@goDir=$${GOPATH%%:*}/src/$(REPO_PATH) && \
				mkdir -p $$(dirname $$goDir) && \
				if [ ! -e "$$goDir" ]; then \
					ln -sfv "$(realpath $(CURDIR)/..)" "$$goDir"; \
				fi

.PHONY: build
build: ensure_gopath 
	mkdir -p bin/ 
	go build -o bin/update-site -ldflags="$(LDFLAGS)" $(REPO_PATH)/update-server

.PHONY: build-docker
build-docker:
	docker run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/update-server build BUILD_VERSION=$(BUILD_VERSION) 

.PHONY: image
image: build-docker 
image: 
	docker pull $(REPOSITORY):latest || true
	docker build -t $(IMAGE) --rm . 
	echo $(IMAGE) > image 

.PHONY: metalint
metalint:
	gometalinter -E goimports -D unconvert -D goconst -D gotype -D aligncheck -D structcheck -D deadcode -D gocyclo -D interfacer -D errcheck -D dupl -D vetshadow -D golint -D vet -D varcheck ./...

.PHONY: test
test: metalint unit

.PHONY: unit
unit: ensure_gopath 
	go test -v -timeout=4s $(REPO_PATH)/update-server

.PHONY: clean
clean:
	rm -f bin/update-site
