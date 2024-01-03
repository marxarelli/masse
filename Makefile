.PHONY: all masse massed test

all: masse massed
masse massed:
	CGO_ENABLED=0 go build ./cmd/$@

test:
	go test ./...

clean:
	rm -f masse
