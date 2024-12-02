.PHONY: all masse massed test

GCFLAGS := ""

# Compile without optimizations (when debugging)
#GCFLAGS := "all=-N -l"

all: masse massed
masse massed:
	CGO_ENABLED=0 go build -gcflags $(GCFLAGS) ./cmd/$@

test:
	go test ./...

clean:
	rm -f masse
