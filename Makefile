GOPATH:=$(CURDIR)/.gopath/:$(CURDIR)/Godeps/_workspace
ORG_PATH:=gitHub.***REMOVED***/monsoon
REPO_PATH:=$(ORG_PATH)/arc

bin:
	mkdir bin
.PHONY: build
build: gopath bin
	go build -o bin/arc $(REPO_PATH)

.PHONY: test
test: gopath
	go test $(REPO_PATH) -v

.PHONY: gopath
gopath: .gopath/$(REPO_PATH)

.PHONY: env
env: gopath
	echo GOPATH=$(GOPATH)

.gopath/$(REPO_PATH):
	mkdir -p .gopath/$(ORG_PATH)
	ln -s ../../../.. .gopath/$(REPO_PATH)
