GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc
BUILD_DIR:=bin
ARC_BINARY:=$(BUILD_DIR)/arc
US_BINARY:=$(BUILD_DIR)/update-site
GITVERSION:=-X main.GITCOMMIT `git rev-parse --short HEAD`

TARGETS:=linux/amd64 windows/amd64 darwin/amd64

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build             - build the binary, output to $(ARC_BINARY)"
	@echo "  * build-update-site - build the update site, output to $(US_BINARY)"
	@echo "  * test              - run all tests"
	@echo "  * test-win          - run tests on windows (requires running vagrant vm)"
	@echo "  * gopath            - print custom GOPATH external use" 
	@echo "  * install-deps      - build and cache dependencies (speeds up make build)" 
	@echo "  * cross             - cross compile for darwin, windows, linux (requires docker)" 

.PHONY: build
build: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(ARC_BINARY) -ldflags="$(GITVERSION)" $(REPO_PATH)

.PHONY: build-update-site
build-update-site: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(US_BINARY) $(REPO_PATH)/update-server

.PHONY: test
test: setup
	go test ./... -v

.PHONY: test-win
test-win: 
	vagrant provision --provision-with shell

.PHONY: gopath 
gopath: setup
	@echo $(GOPATH)

.PHONY: setup
setup: .gopath/src/$(REPO_PATH)

.PHONY: install-deps
install-deps:
	jq -r .Deps[].ImportPath < Godeps/Godeps.json |xargs -L1 go install

.PHONY: cross
cross:
	@# -w omit DWARF symbol table -> smaller
	docker run \
		--rm \
		-v $(CURDIR):/gonative/src/gitHub.***REMOVED***/monsoon/arc \
		gonative \
		gox -osarch="$(TARGETS)" -output="bin/arc_{{.OS}}" -ldflags="-w $(GITVERSION)"

.gopath/src/$(REPO_PATH):
	mkdir -p .gopath/src/$(ORG_PATH)
	ln -s ../../../.. .gopath/src/$(REPO_PATH)
