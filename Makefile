#Workaround for concourse not using ENV statements in registry v2 images
ifeq ($(GOPATH),)
  PATH ?= /go/bin:/usr/local/bin:$(PATH)
  export http_proxy := http://proxy.***REMOVED***:8080
  export https_proxy := http://proxy.***REMOVED***:8080
  export no_proxy := sap.corp,127.0.0.1
  export GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
else
  GOPATH:=$(CURDIR)/.gopath:$(CURDIR)/Godeps/_workspace
endif
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc
BUILD_DIR:=bin
ARC_BINARY:=$(BUILD_DIR)/arc
US_BINARY:=$(BUILD_DIR)/update-site
API_BINARY:=$(BUILD_DIR)/api-server
LDFLAGS:=-s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT `git rev-parse --short HEAD`
TARGETS:=linux/amd64 windows/amd64
BUILD_IMAGE:=docker.***REMOVED***/monsoon/arc-build

ARC_BIN_TPL:=arc_{{.OS}}_{{.Arch}}
ifneq ($(BUILD_VERSION),)
LDFLAGS += -X gitHub.***REMOVED***/monsoon/arc/version.Version $(BUILD_VERSION)
ARC_BIN_TPL:=arc_$(BUILD_VERSION)_{{.OS}}_{{.Arch}}
endif


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
	go build -o $(ARC_BINARY) -ldflags="$(LDFLAGS)" $(REPO_PATH)

.PHONY: build-update-site
build-update-site: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(US_BINARY) $(REPO_PATH)/update-server

.PHONY: build-api
build-api: setup
	@mkdir -p $(BUILD_DIR)
	go build -o $(API_BINARY) -ldflags="$(LDFLAGS)" $(REPO_PATH)/api-server

.PHONY: build-all
build-all: build build-update-site build-api

.PHONY: test
test: fmt unit

.PHONY: unit
unit: setup
	go test -v -timeout=4s ./...

.PHONY: fmt
fmt:
	which goimports > /dev/null
	dirs=`go list -f "{{.Dir}}" ./...|grep -v update-server`; \
		test -z "`for d in $$dirs; do goimports -l $$d/*.go | tee /dev/stderr; done`"


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


.PHONY: build-image
build-image: gonative_linux
	docker build -t $(BUILD_IMAGE) .

gonative_linux:
	docker run --rm -i -e https_proxy=http://proxy.***REMOVED***:8080 golang@1.4.2 bash -c "go get -u github.com/inconshreveable/gonative && cat /go/bin/gonative" > gonative_linux
	chmod +x gonative_linux

.PHONY: cross
cross:
	@# -w omit DWARF symbol table -> smaller
	@# -s stip binary
	docker run \
		--rm \
		-v $(CURDIR):/arc \
		$(BUILD_IMAGE) \
		make -C /arc cross-compile TARGETS="$(TARGETS)" BUILD_VERSION=$(BUILD_VERSION)

.PHONY: cross-compile
cross-compile: setup
	gox -osarch="$(TARGETS)" -output="bin/$(ARC_BIN_TPL)" -ldflags="$(LDFLAGS)"

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

.PHONY: clean
clean:
	make -C api-server clean
	make -C update-server clean
	rm -f gonative_linux
