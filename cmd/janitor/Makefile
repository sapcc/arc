PKG_NAME:=gitHub.***REMOVED***/monsoon/arc/cmd/janitor
BUILD_DIR:=../../bin
ARC_JANITOR_BINARY:=$(BUILD_DIR)/janitor
IMAGE := sapcc/arc_janitor
VERSION=0.1.0
LDFLAGS=-X main.GITCOMMIT=`git rev-parse --short HEAD`
CGO_ENABLED=0
GOARCH=amd64

ifneq ($(BUILD_VERSION),)
LDFLAGS += -X gitHub.***REMOVED***/monsoon/arc/version.Version=$(BUILD_VERSION)
endif

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(ARC_JANITOR_BINARY) -ldflags="$(LDFLAGS)" $(PKG_NAME)

linux: export GOOS=linux
linux: build

docker: linux
	docker build -t $(IMAGE):$(VERSION) .

push:
	docker push $(IMAGE):$(VERSION)
