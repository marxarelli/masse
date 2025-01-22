// syntax=registry:5000/masse:experimental
package main

import (
	"github.com/marxarelli/foocue/foo"
)

chains: {
	example: [
		{ image: "debian:stable" },
		{ sh: "echo '\(foo.Foo.foo)' > /foo" },
	]
}

targets: {
	example: {
		build: "example"
	}
}
