.PHONY: all masse massed test

TAG ?= $(shell git describe --tags --abbrev=0)
REPO ?= marxarelli/masse

MODULE := $(shell go list -m)
GOBUILD_FLAGS := -ldflags "-X $(MODULE)/schema.Tag=$(TAG)"
PACKAGE ?= ./...
GOTEST_FLAGS ?=

REGISTRY_AUTH = $(shell jq -r '.registries."registry.cue.works" | (.token_type + " " + .access_token)' ~/.config/cue/logins.json)
export REGISTRY_AUTH

# Compile without optimizations (when debugging)
GOBUILD_DEBUG_FLAGS := -gcflags "all=-N -l"
GOBUILD_FLAGS ?=

define buildx_build
	docker buildx build \
		--build-arg CUE_REGISTRY_AUTH_SECRET.registry.cue.works=auth \
		--secret id=auth,env=REGISTRY_AUTH \
		--file .pipeline/masse.cue \
		--target gateway \
		--build-arg PARAMETER_version='"$(TAG)"' \
		--tag $(REPO):$(TAG) \
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
	dlv exec $(PACKAGE).test
	rm $(PACKAGE).test

.PHONY: clean
clean:
	rm -f masse massed

.PHONY: image
image:
	$(call buildx_build,"--load")

.PHONY: release
release:
	$(call buildx_build,"--push")
