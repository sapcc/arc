GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc
BUILD_DIR:=bin
ARC_BINARY:=$(BUILD_DIR)/arc
US_BINARY:=$(BUILD_DIR)/update-site
API_BINARY:=$(BUILD_DIR)/api-server
GITVERSION:=-X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT `git rev-parse --short HEAD`

TARGETS:=linux/amd64 windows/amd64 darwin/amd64

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build             - build the binary, output to $(ARC_BINARY)"
	@echo "  * build-update-site - build the update site, output to $(US_BINARY)"
	@echo "  * build-api         - build the api server, output to $(API_BINARY)"
	@echo "  * build-all         - build everything" 
	@echo "  * test              - run all tests"
	@echo "  * unit              - run unit tests"
	@echo "  * integration test  - run integration tests"
	@echo "  * test-win          - run tests on windows (requires running vagrant vm)"
	@echo "  * gopath            - print custom GOPATH external use" 
	@echo "  * install-deps      - build and cache dependencies (speeds up make build)" 
	@echo "  * cross             - cross compile for darwin, windows, linux (requires docker)" 
	@echo "  * run-ubuntu        - run bin/arc_linux in a docker container" 
	@echo "  * run-rhel          - run bin/arc_linux in a docker container" 
	@echo "  * run-sles          - run bin/arc_linux in a docker container" 
	@echo "  * up                - run dev stack in iTerm tabs" 

.PHONY: build
build: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(ARC_BINARY) -ldflags="$(GITVERSION)" $(REPO_PATH)

.PHONY: build-update-site
build-update-site: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(US_BINARY) $(REPO_PATH)/update-server

.PHONY: build-api
build-api: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(API_BINARY) -ldflags="$(GITVERSION)" $(REPO_PATH)/api-server

.PHONY: build-all
build-all: build build-update-site build-api

.PHONY: test
test: test-gofmt unit integration 

.PHONY: unit
unit: setup
	go test -v -timeout=2s ./...

.PHONY: integration
integration: setup
	go test ./... -v -p=1 -timeout=30s -tags=integration

.PHONY: test-gofmt
test-gofmt:
	@fmt_fails=`gofmt -l **/*.go | grep -v '^Godep'`; \
		if [ -n "$$fmt_fails" ]; then \
		echo The following files are not gofmt compatiable:; \
		echo $$fmt_fails; \
		exit 1; \
		fi;


.PHONY: test-win
test-win: 
	vagrant provision --provision-with shell

.PHONY: run-ubuntu 
run-ubuntu: 
	docker run \
		--rm \
		-v $(CURDIR)/bin/arc_linux:/arc \
		ubuntu-arc \
		/arc $(ARGS)

.PHONY: run-rhel 
run-rhel: 
	docker run \
		--rm \
		-v $(CURDIR)/bin/arc_linux:/arc \
		rhel7-arc \
		/arc $(ARGS)
.PHONY: run-sles
run-sles: 
	docker run \
		--rm \
		-v $(CURDIR)/bin/arc_linux:/arc \
		sles11-arc \
		/arc $(ARGS)

.PHONY: build-sles
build-sles:
	docker build -f scripts/Dockerfile.sles11 -t sles11-arc scripts/
.PHONY: build-rhel
build-rhel:
	docker build -f scripts/Dockerfile.rhel7 -t rhel7-arc scripts/
.PHONY: build-ubuntu
build-ubuntu:
	docker build -f scripts/Dockerfile.ubuntu -t ubuntu-arc scripts/

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

.PHONY: up
up:
	osascript $(CURDIR)/scripts/arcup.applescript

.PHONY: assets
assets: service/assets_linux/runsv service/assets_linux/svlogd service/assets_linux/sv service/assets_windows/nssm.exe
	go generate $(ORG_PATH)/arc/service

service/assets_linux/%:
	mkdir -p service/assets_linux
	docker build -f scripts/Dockerfile.runit -t static-runit scripts/
	docker run --rm static-runit cat /musl/src/$* > $@
	chmod +x $@

service/assets_windows/nssm.exe:
	mkdir -p service/assets_windows
	wget -O $(dir $@)nssm.zip http://www.nssm.cc/release/nssm-2.24.zip
	unzip -p $(dir $@)nssm.zip nssm-2.24/win64/nssm.exe > $@
	rm -f $(dir $@)nssm.zip

.gopath/src/$(REPO_PATH):
	mkdir -p .gopath/src/$(ORG_PATH)
	ln -s ../../../.. .gopath/src/$(REPO_PATH)
