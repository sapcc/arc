export GO15VENDOREXPERIMENT=1
PKG_NAME:=gitHub.***REMOVED***/monsoon/arc
BUILD_DIR:=bin
ARC_BINARY:=$(BUILD_DIR)/arc
US_BINARY:=$(BUILD_DIR)/update-site
API_BINARY:=$(BUILD_DIR)/api-server
LDFLAGS:=-s -w -X gitHub.***REMOVED***/monsoon/arc/version.GITCOMMIT=`git rev-parse --short HEAD`
TARGETS:=linux/amd64 windows/amd64
BUILD_IMAGE:=hub.***REMOVED***/monsoon/gobuild:1.8

ARC_BIN_TPL:=arc_{{.OS}}_{{.Arch}}
ifneq ($(BUILD_VERSION),)
LDFLAGS += -X gitHub.***REMOVED***/monsoon/arc/version.Version=$(BUILD_VERSION)
ARC_BIN_TPL:=arc_$(BUILD_VERSION)_{{.OS}}_{{.Arch}}
endif

packages = $(PKG_NAME) $(shell go list -f '{{ join .Deps "\n" }}' | grep -v vendor | grep $(PKG_NAME))

.PHONY: help 
help:
	@echo
	@echo "Available targets:"
	@echo "  * build             - build the binary, output to $(ARC_BINARY)"
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
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(ARC_BINARY) -ldflags="$(LDFLAGS)" $(PKG_NAME)

.PHONY: test
test: metalint unit

.PHONY: unit
unit:
	go test -v -timeout=4s $(packages)

.PHONY: metalint
metalint:
	gometalinter -E goimports -D goconst -D unconvert -D gotype -D aligncheck -D structcheck -D deadcode -D gocyclo -D interfacer -D errcheck -D dupl -D vetshadow -D golint -D vet -D varcheck --deadline=10s $(shell echo $(packages) | sed -e 's!$(PKG_NAME)!.!g')

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


.PHONY: build-image
build-image:
	docker build -t $(BUILD_IMAGE) .

.PHONY: cross
cross:
	@# -w omit DWARF symbol table -> smaller
	@# -s stip binary
	docker run \
		--rm \
		-v $(CURDIR):/go/src/$(PKG_NAME) \
		-w /go/src/$(PKG_NAME) \
		$(BUILD_IMAGE) \
		make cross-compile TARGETS="$(TARGETS)" BUILD_VERSION=$(BUILD_VERSION)

.PHONY: cross-compile
cross-compile:
	gox -osarch="$(TARGETS)" -output="bin/$(ARC_BIN_TPL)" -ldflags="$(LDFLAGS)" $(PKG_NAME)

.PHONY: up
up:
	osascript $(CURDIR)/scripts/arcup.applescript

.PHONY: assets
assets: service/assets_linux/runsv service/assets_linux/svlogd service/assets_linux/sv service/assets_windows/nssm.exe
	go generate $(PKG_NAME)/service

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

.PHONY: clean
clean:
	make -C api-server clean
	make -C update-server clean
