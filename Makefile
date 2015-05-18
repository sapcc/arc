GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc
BUILD_DIR:=bin
BINARY:=$(BUILD_DIR)/arc

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build   - build the binary, output to $(BINARY)"
	@echo "  * test    - run all tests"
	@echo "  * gopath  - print custom GOPATH external use" 

.PHONY: build
build: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(REPO_PATH)

.PHONY: test
test: setup
	go test $(REPO_PATH) -v

.PHONY: gopath 
gopath: setup
	@echo $(GOPATH)

.PHONY: setup
setup: .gopath/src/$(REPO_PATH)

#file targets below

.gopath/src/$(REPO_PATH):
	mkdir -p .gopath/src/$(ORG_PATH)
	ln -s ../../../.. .gopath/src/$(REPO_PATH)
