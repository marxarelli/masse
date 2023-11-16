.PHONY: all phyton test

all: phyton
phyton:
	go build ./cmd/phyton

test:
	go test ./...

clean:
	rm -f phyton
