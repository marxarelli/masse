PROTO := $(shell find . -type f -name '*.proto' -not -path "./vendor/*")
PBGO := $(patsubst %.proto,%.pb.go,$(PROTO))
PROTOC_INC := $(shell go list -m -f "{{.Dir}}" gitlab.wikimedia.org/dduvall/protoyaml)

.PHONY: proto

proto: $(PBGO)
%.pb.go: %.proto
	protoc -I$(PROTOC_INC) -I. --proto_path=. --go_out=. --go_opt=paths=source_relative $<

clean:
	rm $(PBGO)

test: proto
	go test ./...
