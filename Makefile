PROTO := $(shell find . -type f -name '*.proto' -not -path "./vendor/*")
PBGO := $(patsubst %.proto,%.pb.go,$(PROTO))
PROTO_DEPS := $(shell go list -m -f "{{.Dir}}" gitlab.wikimedia.org/dduvall/protoyaml)
PROTOC_INC := $(addprefix -I,$(PROTO_DEPS))

.PHONY: proto

proto: $(PBGO)
%.pb.go: %.proto
	protoc $(PROTOC_INC) -I. --proto_path=. --go_out=. --go_opt=paths=source_relative $<

clean:
	rm $(PBGO)

test: proto
	go test ./...
