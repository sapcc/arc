REPO_PATH   := gitHub.***REMOVED***/monsoon/arc
LDFLAGS     := -s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT=`git rev-parse --short HEAD`
REPOSITORY  := hub.***REMOVED***/monsoon/arc-api
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := hub.***REMOVED***/monsoon/gobuild:1.7-alpine

ifneq ($(BUILD_VERSION),)
LDFLAGS += -X gitHub.***REMOVED***/monsoon/arc/version.Version=$(BUILD_VERSION)
endif

ifneq ($(http_proxy),)
BUILD_ARGS+= --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) --build-arg no_proxy=$(no_proxy)
endif

PKG_ROOT := $(realpath $(CURDIR)/.. )
GOFILES:= $(shell find $(PKG_ROOT) -name '*.go' -not -path '$(PKG_ROOT)/vendor/*' )


.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/api-server"
	@echo "  * build-docker - build static binary via golang container, output to bin/api-server"
	@echo "  * image        - build docker image containing the update site"

build: bin/api-server bin/gencert

bin/api-server: $(GOFILES)
	@mkdir -p bin/
	go build -o bin/api-server -ldflags="$(LDFLAGS)" $(REPO_PATH)/api-server

bin/gencert:
	@mkdir -p bin/
	go get -d gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert
	go build -o bin/gencert -ldflags="-s -w" gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert

.PHONY: build-docker
build-docker:
	docker run \
		--rm \
		-v $(CURDIR)/..:/go/src/$(REPO_PATH) \
		-w /go/src/$(REPO_PATH)/api-server \
		$(BUILD_IMAGE) \
		make build BUILD_VERSION=$(BUILD_VERSION)

.PHONY: image
image: build-docker
image:
	docker build $(BUILD_ARGS) -t $(IMAGE) --rm .

.PHONY: clean
clean:
	rm -f bin/api-server
	rm -f bin/gencert
