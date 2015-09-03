SHELL       := /bin/sh
REPOSITORY  := docker.***REMOVED***/monsoon/arc-website-builder
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)

### Executables
DOCKER      = docker
BUILD_IMAGE = docker.***REMOVED***/monsoon/docker-build:1.5.0 
WAIT        = $(DOCKER) run --rm --link $(WAIT_ID):wait $(BUILD_IMAGE) wait
MTIMES      = $(DOCKER) run --rm $(MTIMES_OPTS)         $(BUILD_IMAGE) reset_mtimes

# ----------------------------------------------------------------------------------
#   image 
# ----------------------------------------------------------------------------------
image: build
	echo $(IMAGE) > image

# ----------------------------------------------------------------------------------
#   build 
# ----------------------------------------------------------------------------------
#
# Build and tags an image from a Dockerfile.
#
# We need to reset the modification times of all files in a git repository to the 
# date of the latest commit that changed it. This is required because Docker takes 
# the mtime of a file into account when checking for modifications. Git on the other
# hand does not. Without this the cache will be busted by Git and is basically
# useless.
#
build: MTIMES_OPTS = -v $(shell pwd):/src
build: 
	$(MTIMES)
	$(DOCKER) pull $(REPOSITORY):latest || true
	$(DOCKER) build -t $(IMAGE) --rm . 
	echo $(IMAGE) > build

# ----------------------------------------------------------------------------------
#   clean 
# ----------------------------------------------------------------------------------
#
# Kill and remove all containers. Remove intermediate files. 
#
.PHONY: 
clean: 	
	$(RM) image build
