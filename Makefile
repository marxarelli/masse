.PHONY: all masse massed test

GOBUILD_FLAGS :=

# Compile without optimizations (when debugging)
#GOBUILD_FLAGS := -gcflags "all=-N -l"

all: masse massed
masse massed:
	CGO_ENABLED=0 go build $(GOBUILD_FLAGS) ./cmd/$@

test:
	go test ./...

clean:
	rm -f masse
