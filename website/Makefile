REPOSITORY  := hub.***REMOVED***/monsoon/arc-website-builder
TAG         ?= latest
IMAGE       := $(REPOSITORY):$(TAG)


build: 
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)

