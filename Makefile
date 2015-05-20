GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc
BUILD_DIR:=bin
BINARY:=$(BUILD_DIR)/arc

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build       - build the binary, output to $(BINARY)"
	@echo "  * test        - run all tests"
	@echo "  * gopath      - print custom GOPATH external use" 
	@echo "  * build-deps  - build and cache dependencies (speeds up make build)" 

.PHONY: build
build: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(REPO_PATH)

.PHONY: test
test: setup
	go test ./... -v

.PHONY: gopath 
gopath: setup
	@echo $(GOPATH)

.PHONY: setup
setup: .gopath/src/$(REPO_PATH)

.PHONY: install-deps
install-deps:
	jq -r .Deps[].ImportPath < Godeps/Godeps.json |xargs -L1 go install

#file targets below

.gopath/src/$(REPO_PATH):
	mkdir -p .gopath/src/$(ORG_PATH)
	ln -s ../../../.. .gopath/src/$(REPO_PATH)
