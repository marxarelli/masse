.PHONY: all phyton phytond test

all: phyton phytond
phyton phytond:
	CGO_ENABLED=0 go build ./cmd/$@

test:
	go test ./...

clean:
	rm -f phyton
