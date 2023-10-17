package common

import(
	"strings"
)

// See https://github.com/opencontainers/image-spec/blob/main/config.md#properties
// and https://go.dev/doc/install/source#environment
#OS:
	"aix" |
	"android" |
	"darwin" |
	"dragonfly" |
	"freebsd" |
	"illumos" |
	"ios" |
	"js" |
	"linux" |
	"netbsd" |
	"openbsd" |
	"plan9" |
	"solaris" |
	"wasip1" |
	"windows"

#Architecture:
	"386" |
	"amd64" |
	"arm" |
	"arm64" |
	"loong64" |
	"mips" |
	"mips64" |
	"mips64le" |
	"mipsle" |
	"ppc64" |
	"ppc64le" |
	"riscv64" |
	"s390x" |
	"wasm"

#Variants: {
	arm:   "v6" | "v7" | "v8"
	arm64: "v8"
}

#Platform: {
	name?: string

	os:           #OS
	architecture: #Architecture
	variant?:     *#Variants[architecture] | ""

	if name != _|_ {
		_parts:    strings.SplitN(name, "/", 3)
		os: _parts[0]
		architecture: _parts[1]
		if len(_parts) > 2 {
			variant: _parts[2]
		}
	}
}
