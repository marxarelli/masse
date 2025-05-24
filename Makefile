.PHONY: all masse massed test

TAG ?= $(shell git describe --tags --abbrev=0)
REPO ?= marxarelli/masse
IMAGE_NAME := $(REPO):$(TAG)
IMAGE_BUILD_OUTPUT := --output type=oci,name=$(IMAGE_NAME)

MODULE := $(shell go list -m)
GOBUILD_FLAGS := -ldflags "-X $(MODULE)/schema.Tag=$(TAG)"
PACKAGE ?= ./...
PACKAGE_BASENAME := $(shell basename $(PACKAGE))
GOTEST_FLAGS ?=

# Compile without optimizations (when debugging)
GOBUILD_DEBUG_FLAGS := -gcflags "all=-N -l"

BUILDX_BUILD_FLAGS ?=

define buildx_build
	docker buildx build \
		--file .pipeline/masse.cue \
		--target gateway \
		--build-arg masse:parameter:tag='"$(TAG)"' \
		--tag $(IMAGE_NAME) \
		--attest type=provenance,mode=max \
		$(BUILDX_BUILD_FLAGS) \
		$(1) \
		.
endef

all: bin/masse bin/massed
bin/masse bin/massed:
	CGO_ENABLED=0 go build $(GOBUILD_FLAGS) -o $@ ./cmd/$(notdir $@)

.PHONY: benchmark
benchmark:
	go test -bench=. $(GOBUILD_FLAGS) $(GOTEST_FLAGS) $(PACKAGE)

.PHONY: test
test:
	go test $(GOBUILD_FLAGS) $(GOTEST_FLAGS) $(PACKAGE)

.PHONY: debug-test
debug-test:
	go test -c $(GOBUILD_FLAGS) $(GOBUILD_DEBUG_FLAGS) $(GOTEST_FLAGS) $(PACKAGE)
	dlv exec $(PACKAGE_BASENAME).test
	rm $(PACKAGE_BASENAME).test

.PHONY: clean
clean:
	rm -f masse massed bin/*
	rm -rf build/*

.PHONY: image
image:
	mkdir -p build/gateway
	$(call buildx_build,$(IMAGE_BUILD_OUTPUT)) | tar -C build/gateway -xf -
	jq . < build/gateway/index.json

.PHONY: release
release:
	$(call buildx_build,--attest type=sbom --push)

.PHONY: publish-schema
publish-schema:
	cd schema && cue mod publish $(TAG)
