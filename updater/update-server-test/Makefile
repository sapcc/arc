export GO15VENDOREXPERIMENT=1
REPO_PATH   := github.com/sapcc/arc
LDFLAGS     := -s -w -X github.com/sapcc/arc/version.GITCOMMIT=`git rev-parse --short HEAD`

ifneq ($(BUILD_VERSION),)
LDFLAGS:=$(LDFLAGS) -X github.com/sapcc/arc/version.Version=$(BUILD_VERSION)
endif

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build        - build static binary, output to bin/update-site"

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
	go build -o bin/update-site -ldflags="$(LDFLAGS)"