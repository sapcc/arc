export GO15VENDOREXPERIMENT=1
REPO_PATH   := gitHub.***REMOVED***/monsoon/arc
LDFLAGS     := -s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT=`git rev-parse --short HEAD`
REPOSITORY  := docker.***REMOVED***/monsoon/arc-api
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)
BUILD_IMAGE := docker.***REMOVED***/monsoon/gobuild@1.6

ifneq ($(BUILD_VERSION),)
LDFLAGS += -X gitHub.***REMOVED***/monsoon/arc/version.Version=$(BUILD_VERSION)
endif

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

.PHONY: ensure_gopath
ensure_gopath:
	goDir=$${GOPATH%%:*}/src/$(REPO_PATH) && \
				mkdir -p $$(dirname $$goDir) && \
				if [ ! -e "$$goDir" ]; then \
					ln -sfv "$(realpath $(CURDIR)/..)" "$$goDir"; \
				fi

.PHONY: build
build: ensure_gopath 
	mkdir -p bin/ 
	go build -o bin/api-server -ldflags="$(LDFLAGS)" $(REPO_PATH)/api-server

.PHONY: gencert
gencert:
	go get -d gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert
	go build -o bin/gencert -ldflags="-s -w" gitHub.***REMOVED***/monsoon/arc-pki.git/src/gencert

.PHONY: build-docker
build-docker: ensure_gopath 
	$(DOCKER) run \
		--rm \
		-v $(CURDIR)/..:/arc \
		$(BUILD_IMAGE) \
		make -C /arc/api-server build gencert BUILD_VERSION=$(BUILD_VERSION)

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
run-tests: ensure_gopath
	go test -v $(REPO_PATH)/api-server
	go test -v -p=1 -timeout=60s -tags=integration $(REPO_PATH)/api-server

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
	rm -f postgres
