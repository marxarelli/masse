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

REGISTRY_AUTH = $(shell jq -r '.registries."registry.cue.works" | (.token_type + " " + .access_token)' ~/.config/cue/logins.json)
export REGISTRY_AUTH

# Compile without optimizations (when debugging)
GOBUILD_DEBUG_FLAGS := -gcflags "all=-N -l"

BUILDX_BUILD_FLAGS ?=

define buildx_build
	docker buildx build \
		--build-arg CUE_REGISTRY_AUTH_SECRET.registry.cue.works=auth \
		--secret id=auth,env=REGISTRY_AUTH \
		--file .pipeline/masse.cue \
		--target gateway \
		--build-arg PARAMETER_version='"$(TAG)"' \
		--tag $(IMAGE_NAME) \
		--attest type=provenance,mode=max \
		$(BUILDX_BUILD_FLAGS) \
		$(1) \
		.
endef

all: masse massed
masse massed:
	CGO_ENABLED=0 go build $(GOBUILD_FLAGS) ./cmd/$@

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
	rm -f masse massed

.PHONY: image
image:
	mkdir -p build/gateway
	$(call buildx_build,$(IMAGE_BUILD_OUTPUT)) | tar -C build/gateway -xf -
	jq . < build/gateway/index.json

.PHONY: release
release:
	$(call buildx_build,--push)
